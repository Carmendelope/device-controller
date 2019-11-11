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

package commands

import (
	"github.com/nalej/device-controller/pkg/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var config = server.Config{}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run Device Controller",
	Long:  `Run Device Controller`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		log.Info().Msg("Launching API!")
		server := server.NewService(config)
		server.Run()
	},
}

func init() {
	runCmd.Flags().IntVar(&config.Port, "port", 6020, "Port to launch the Device gRPC API")
	runCmd.Flags().IntVar(&config.HTTPPort, "httpPort", 6021, "Port to launch the Device HTTP API")
	runCmd.Flags().IntVar(&config.Threshold, "threshold", 100, "Threshold for latency")
	runCmd.Flags().StringVar(&config.ClusterAPIHostname, "clusterAPIHostname", "", "Hostname of the cluster API on the management cluster")
	runCmd.Flags().Uint32Var(&config.ClusterAPIPort, "clusterAPIPort", 8000, "Port where the cluster API is listening")
	runCmd.Flags().StringVar(&config.LoginHostname, "loginHostname", "", "Hostname of the login service")
	runCmd.Flags().Uint32Var(&config.LoginPort, "loginPort", 31683, "port where the login service is listening")
	runCmd.Flags().BoolVar(&config.UseTLSForLogin, "useTLSForLogin", true, "Use TLS to connect to the Login API")
	runCmd.Flags().StringVar(&config.Email, "email", "", "email address")
	runCmd.Flags().StringVar(&config.Password, "password", "", "password")
	runCmd.PersistentFlags().StringVar(&config.AuthHeader, "authHeader", "", "Authorization Header")
	runCmd.PersistentFlags().StringVar(&config.AuthConfigPath, "authConfigPath", "", "Authorization config path")
	runCmd.Flags().StringVar(&config.CACertPath, "caCertPath", "", "Path for the CA certificate")
	runCmd.Flags().StringVar(&config.ClientCertPath, "clientCertPath", "", "Path for the client certificate")
	runCmd.Flags().BoolVar(&config.SkipServerCertValidation, "skipServerCertValidation", true, "Skip CA authentication validation")
	rootCmd.AddCommand(runCmd)
}
