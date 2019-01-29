/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package server

import (
	"github.com/nalej/derrors"
	"github.com/rs/zerolog/log"
	"github.com/nalej/device-controller/version"
)

type Config struct {
	// Port where the gRPC API service will listen  to incoming requests
	Port int
}

func (conf *Config) Validate() derrors.Error {

	if conf.Port <= 0 {
		return derrors.NewInvalidArgumentError("ports must be valid")
	}
	return nil
}

func (conf *Config) Print() {
	log.Info().Str("app", version.AppVersion).Str("commit", version.Commit).Msg("Version")
	log.Info().Int("port", conf.Port).Msg("gRPC port")
}
