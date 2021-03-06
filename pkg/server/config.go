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
	"github.com/nalej/authx/pkg/interceptor"
	"github.com/nalej/derrors"
	"github.com/nalej/device-controller/version"
	"github.com/rs/zerolog/log"
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
	LoginPort      uint32
	UseTLSForLogin bool
	// Email to log into the management cluster.
	Email string
	// Password to log into the managment cluster.
	Password string
	// Threshold in milliseconds by which it will be considered if a latency is acceptable or not
	Threshold int
	// AuthHeader contains the name of the target header.
	AuthHeader string
	// AuthConfigPath contains the path of the file with the authentication configuration.
	AuthConfigPath string
	// Path for the certificate of the CA
	CACertPath string
	// Client Cert Path
	ClientCertPath string
	// Skip Server validation
	SkipServerCertValidation bool
}

// LoadAuthConfig loads the security configuration.
func (conf *Config) LoadAuthConfig() (*interceptor.AuthorizationConfig, derrors.Error) {
	return interceptor.LoadAuthorizationConfig(conf.AuthConfigPath)
}

func (conf *Config) Validate() derrors.Error {

	if conf.Port <= 0 {
		return derrors.NewInvalidArgumentError("ports must be valid")
	}
	if conf.Threshold <= 0 {
		return derrors.NewInvalidArgumentError("Threshold must be valid")
	}
	if conf.AuthConfigPath == "" {
		return derrors.NewInvalidArgumentError("authConfigPath must be set")
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
	log.Info().Str("header", conf.AuthHeader).Msg("Authorization")
	log.Info().Str("path", conf.AuthConfigPath).Msg("Permissions file")
}
