/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package server

import (
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-utils/pkg/tools"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
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

// Run the service, launch the service handler
func (s *Service) Run() error {
	cErr := s.Configuration.Validate()
	if cErr != nil {
		log.Fatal().Str("err", cErr.DebugReport()).Msg("invalid configuration")
	}
	s.Configuration.Print()

	// create clients

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Configuration.Port))
	if err != nil {
		log.Fatal().Errs("failed to listen: %v", []error{err})
	}

	// Create handlers and managers

	grpcServer := grpc.NewServer()

	// register

	reflection.Register(grpcServer)
	log.Info().Int("port", s.Configuration.Port).Msg("Launching gRPC server")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Errs("failed to serve: %v", []error{err})
	}
	return nil
}