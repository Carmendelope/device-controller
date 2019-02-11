/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package server

import (
	"crypto/tls"
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/device-controller/pkg/login_helper"
	"github.com/nalej/device-controller/pkg/server/ping"
	"github.com/nalej/grpc-device-controller-go"
	"github.com/nalej/grpc-utils/pkg/tools"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"net"
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

	// Build connection with conductor
	log.Debug().Str("hostname", s.Configuration.ClusterAPIHostname).Msg("connecting with cluster api")
	clusterAPIConn, errCond := s.getClusterAPIConnection(s.Configuration.ClusterAPIHostname, int(s.Configuration.ClusterAPIPort))
	if errCond != nil {
		log.Fatal().Errs("impossible to connect with cluster api", []error{cErr})
	}

	// Create handlers and managers
	pingManager := ping.NewManager(s.Configuration.Threshold, clusterAPIConn)
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

