/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package server

import (
	"github.com/nalej/derrors"
	"github.com/rs/zerolog/log"
	"github.com/nalej/device-controller/version"
	"strings"
)

type Config struct {
	// Port where the gRPC API service will listen  to incoming requests
	Port int
	// HTTPPort where the HTTP gRPC gateway will be listening.
	HTTPPort int
	// ClusterAPIHostname with the hostname of the cluster API on the management cluster
	ClusterAPIHostname string
	// ClusterAPIPort with the port where the cluster API is listening.
	ClusterAPIPort uint32
	// LoginHostname with the hostname of the login API on the management cluster.
	LoginHostname string
	// LoginPort with the port where the login API is listening
	LoginPort uint32
	UseTLSForLogin bool
	// Email to log into the management cluster.
	Email string
	// Password to log into the managment cluster.
	Password string
	// Threshold in milliseconds by which it will be considered if a latency is acceptable or not
	Threshold int
}

func (conf *Config) Validate() derrors.Error {

	if conf.Port <= 0 {
		return derrors.NewInvalidArgumentError("ports must be valid")
	}
	if conf.Threshold <= 0 {
		return derrors.NewInvalidArgumentError("Threshold must be valid")
	}
	return nil
}

func (conf *Config) Print() {
	log.Info().Str("app", version.AppVersion).Str("commit", version.Commit).Msg("Version")
	log.Info().Int("port", conf.Port).Msg("gRPC port")
	log.Info().Int("port", conf.HTTPPort).Msg("HTTP port")
	log.Info().Int("Threshold", conf.Threshold).Msg("Threshold in milliseconds")
	log.Info().Str("URL", conf.ClusterAPIHostname).Uint32("port", conf.ClusterAPIPort).Msg("Cluster API on management cluster")
	log.Info().Str("URL", conf.LoginHostname).Uint32("port", conf.LoginPort).Bool("UseTLSForLogin", conf.UseTLSForLogin).Msg("Login API on management cluster")
	log.Info().Str("Email", conf.Email).Str("password", strings.Repeat("*", len(conf.Password))).Msg("Application cluster credentials")


}
