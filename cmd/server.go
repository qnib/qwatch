package cmd

import (
	"github.com/qnib/qwatch/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// watchSrv loops over nodes, services and tasks
var qServer = &cobra.Command{
	Use:   "server",
	Short: "Starts server daemon",
	Long:  ``,
	Run:   qserver.ServeQlog,
}

func init() {
	RootCmd.AddCommand(qServer)

	RootCmd.PersistentFlags().String("collectors", "Gelf,DockerEvents", "Comma separated list of collectors to start")
	viper.BindPFlag("collectors", RootCmd.PersistentFlags().Lookup("collectors"))
	RootCmd.PersistentFlags().Int("ticker-interval", 15000, "Interval of global ticker in milliseconds")
	viper.BindPFlag("ticker-interval", RootCmd.PersistentFlags().Lookup("ticker-interval"))
	RootCmd.PersistentFlags().Int("gelf-port", 12201, "UDP port of GELF collector")
	viper.BindPFlag("gelf-port", RootCmd.PersistentFlags().Lookup("gelf-port"))
}
