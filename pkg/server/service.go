/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/nalej/derrors"
	"github.com/nalej/device-controller/pkg/login_helper"
	"github.com/nalej/device-controller/pkg/server/ping"
	"github.com/nalej/grpc-cluster-api-go"
	"github.com/nalej/grpc-device-controller-go"
	"github.com/nalej/grpc-utils/pkg/tools"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"strings"
)

// Service structure with the configuration and the gRPC server.
type Service struct {
	Configuration Config
	Server        *tools.GenericGRPCServer
}

// Clients structure with the gRPC clients for remote services.
func NewService(conf Config) *Service {
	return &Service {
		conf,
		tools.NewGenericGRPCServer(uint32(conf.Port)),
	}
}

type Clients struct {
}

func (s *Service) GetClients() (*Clients, derrors.Error) {

	return &Clients{}, nil
}

func (s *Service) getClusterAPIConnection(hostname string, port int) (*grpc.ClientConn, derrors.Error) {
	// Build connection with cluster API
	tlsConfig := &tls.Config{
		ServerName:   hostname,
		InsecureSkipVerify: true,
	}
	targetAddress := fmt.Sprintf("%s:%d", hostname, port)
	log.Debug().Str("address", targetAddress).Msg("creating cluster API connection")

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

	go s.LaunchGRPC()
	return s.LaunchHTTP()

}

func (s * Service) LaunchGRPC() error {
	// create clients
	clusterAPILoginHelper := login_helper.NewLogin(s.Configuration.LoginHostname, int(s.Configuration.LoginPort), s.Configuration.UseTLSForLogin,
		s.Configuration.Email, s.Configuration.Password)

	cErr := clusterAPILoginHelper.Login()
	if cErr != nil {
		log.Fatal().Errs("there was an error requesting cluster-api login", []error{cErr})
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Configuration.Port))
	if err != nil {
		log.Fatal().Errs("failed to listen: %v", []error{err})
	}

	// Build connection with cluster-api
	log.Debug().Str("hostname", s.Configuration.ClusterAPIHostname).Msg("connecting with cluster api")
	clusterAPIConn, errCond := s.getClusterAPIConnection(s.Configuration.ClusterAPIHostname, int(s.Configuration.ClusterAPIPort))
	if errCond != nil {
		log.Fatal().Errs("impossible to connect with cluster api", []error{cErr})
	}
	cClient := grpc_cluster_api_go.NewDeviceManagerClient(clusterAPIConn)

	// Create handlers and managers
	pingManager := ping.NewManager(s.Configuration.Threshold, clusterAPILoginHelper, cClient)
	pingHandler := ping.NewHandler(pingManager)

	grpcServer := grpc.NewServer()
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

func (s * Service) LaunchHTTP() error {

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