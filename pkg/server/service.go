/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/nalej/authx/pkg/interceptor"
	"github.com/nalej/authx/pkg/interceptor/devinterceptor"
	"github.com/nalej/derrors"
	"github.com/nalej/device-controller/pkg/login_helper"
	"github.com/nalej/device-controller/pkg/server/ping"
	"github.com/nalej/grpc-cluster-api-go"
	"github.com/nalej/grpc-device-controller-go"
	"github.com/nalej/grpc-login-api-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

// Service structure with the configuration and the gRPC server.
type Service struct {
	Configuration Config
}

// Clients structure with the gRPC clients for remote services.
func NewService(conf Config) *Service {
	return &Service{
		conf,
	}
}

type Clients struct {
	DeviceManagerClient grpc_cluster_api_go.DeviceManagerClient
	LoginClient         grpc_login_api_go.LoginClient
}

func (s *Service) GetClients() (*Clients, derrors.Error) {

	dmConn, err := s.getSecureAPIConnection(s.Configuration.ClusterAPIHostname, int(s.Configuration.ClusterAPIPort), s.Configuration.CACertPath, s.Configuration.ClientCertPath, s.Configuration.SkipServerCertValidation)
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with the Cluster API manager")
	}
	deviceClient := grpc_cluster_api_go.NewDeviceManagerClient(dmConn)

	loginConn, err := s.getSecureAPIConnection(s.Configuration.LoginHostname, int(s.Configuration.LoginPort), s.Configuration.CACertPath, s.Configuration.ClientCertPath, s.Configuration.SkipServerCertValidation)
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with the Login API manager")
	}
	loginClient := grpc_login_api_go.NewLoginClient(loginConn)

	return &Clients{DeviceManagerClient: deviceClient, LoginClient: loginClient}, nil
}

func (s *Service) getSecureAPIConnection(hostname string, port int, caCertPath string, clientCertPath string, skipCAValidation bool) (*grpc.ClientConn, derrors.Error) {
	// Build connection with cluster API
	rootCAs := x509.NewCertPool()
	tlsConfig := &tls.Config{
		ServerName: hostname,
	}

	if caCertPath != "" {
		log.Debug().Str("caCertPath", caCertPath).Msg("loading CA cert")
		caCert, err := ioutil.ReadFile(caCertPath)
		if err != nil {
			return nil, derrors.NewInternalError("Error loading CA certificate")
		}
		added := rootCAs.AppendCertsFromPEM(caCert)
		if !added {
			return nil, derrors.NewInternalError("cannot add CA certificate to the pool")
		}
		tlsConfig.RootCAs = rootCAs
	}

	targetAddress := fmt.Sprintf("%s:%d", hostname, port)
	log.Debug().Str("address", targetAddress).Msg("creating cluster API connection")

	if clientCertPath != "" {
		log.Debug().Str("clientCertPath", clientCertPath).Msg("loading client certificate")
		clientCert, err := tls.LoadX509KeyPair(fmt.Sprintf("%s/tls.crt", clientCertPath), fmt.Sprintf("%s/tls.key", clientCertPath))
		if err != nil {
			log.Error().Str("error", err.Error()).Msg("Error loading client certificate")
			return nil, derrors.NewInternalError("Error loading client certificate")
		}

		tlsConfig.Certificates = []tls.Certificate{clientCert}
		tlsConfig.BuildNameToCertificate()
	}
	log.Debug().Str("address", targetAddress).Str("caCertPath", caCertPath).Bool("skipCAValidation", skipCAValidation).Msg("creating secure connection")

	if skipCAValidation {
		log.Debug().Msg("skipping server cert validation")
		tlsConfig.InsecureSkipVerify = true
	}

	creds := credentials.NewTLS(tlsConfig)

	log.Debug().Interface("creds", creds.Info()).Msg("Secure credentials")
	sConn, dErr := grpc.Dial(targetAddress, grpc.WithTransportCredentials(creds))
	if dErr != nil {
		return nil, derrors.AsError(dErr, "cannot create connection with the cluster API service")
	}
	return sConn, nil
}

// Run the service, launch the service handler
func (s *Service) Run() error {
	vErr := s.Configuration.Validate()
	if vErr != nil {
		log.Fatal().Str("err", vErr.DebugReport()).Msg("invalid configuration")
	}
	s.Configuration.Print()

	authConfig, authErr := s.Configuration.LoadAuthConfig()
	if authErr != nil {
		log.Fatal().Str("err", authErr.DebugReport()).Msg("cannot load authx config")
	}

	log.Info().Bool("AllowsAll", authConfig.AllowsAll).Int("permissions", len(authConfig.Permissions)).Msg("Auth config")

	go s.LaunchGRPC(authConfig)
	return s.LaunchHTTP()

}

func (s *Service) LaunchGRPC(authConfig *interceptor.AuthorizationConfig) error {
	// create clients
	clients, cErr := s.GetClients()
	if cErr != nil {
		log.Fatal().Str("err", cErr.DebugReport()).Msg("cannot generate clients")
		return cErr
	}

	clusterAPILoginHelper := login_helper.NewLogin(s.Configuration.LoginHostname, int(s.Configuration.LoginPort), s.Configuration.UseTLSForLogin,
		s.Configuration.Email, s.Configuration.Password, "", "", true)

	cErr = clusterAPILoginHelper.Login()
	if cErr != nil {
		log.Fatal().Str("err", cErr.DebugReport()).Msg("there was an error requesting cluster-api login")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Configuration.Port))
	if err != nil {
		log.Fatal().Int("port", s.Configuration.Port).Str("err", err.Error()).Msg("failed to listen")
	}

	// Create handlers and managers
	pingManager := ping.NewManager(s.Configuration.Threshold, clusterAPILoginHelper, clients.DeviceManagerClient)
	pingHandler := ping.NewHandler(pingManager)

	// Interceptor
	accessManager, aErr := devinterceptor.NewClusterApiSecretAccessWithClients(clients.LoginClient, clients.DeviceManagerClient,
		s.Configuration.Email, s.Configuration.Password, devinterceptor.DefaultCacheEntries)
	if err != nil {
		log.Fatal().Str("trace", aErr.DebugReport()).Msg("cannot create management secret access")
	}

	authxConfig := interceptor.NewConfig(authConfig, "", s.Configuration.AuthHeader)
	grpcServer := grpc.NewServer(interceptor.WithDeviceAuthxInterceptor(accessManager, authxConfig))

	//grpcServer := grpc.NewServer()
	grpc_device_controller_go.RegisterConnectionServer(grpcServer, pingHandler)

	// register

	reflection.Register(grpcServer)
	log.Info().Int("port", s.Configuration.Port).Msg("Launching gRPC server")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Errs("failed to serve: %v", []error{err})
	}
	return nil
}

func (s *Service) allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept", "Authorization"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
}

func (s *Service) LaunchHTTP() error {

	addr := fmt.Sprintf(":%d", s.Configuration.HTTPPort)
	clientAddr := fmt.Sprintf(":%d", s.Configuration.Port)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	mux := runtime.NewServeMux()

	if err := grpc_device_controller_go.RegisterConnectionHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatal().Err(err).Msg("failed to start device controller handler")
	}

	server := &http.Server{
		Addr:    addr,
		Handler: s.allowCORS(mux),
	}

	log.Info().Str("address", addr).Msg("HTTP Listening")
	return server.ListenAndServe()

}
