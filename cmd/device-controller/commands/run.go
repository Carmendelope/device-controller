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
		//server := server.NewService(config)
		//server.Run()
	},
}

func init() {
	runCmd.Flags().IntVar(&config.Port, "port", 6020, "Port to launch the Device gRPC API")

	rootCmd.AddCommand(runCmd)
}
