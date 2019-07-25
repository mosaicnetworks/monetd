package server

import (
	"github.com/spf13/cobra"
)

//ServerCmd is the CLI command for the giverny server
var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "giverny server",
	Long: `Server
	
The giverny server is a simple REST server to facilitate the sharing of
Monet configurations prior to instantiation of the node. `,

	TraverseChildren: true,
}

func init() {
	//Subcommands
	ServerCmd.AddCommand(
		newStartCmd(),
		newStopCmd(),
		newStatusCmd(),
	)
	/*
		//Commonly used command line flags
		KeysCmd.PersistentFlags().StringVar(&passwordFile, "passfile", "", "the file that contains the passphrase for the keyfile")
		KeysCmd.PersistentFlags().BoolVar(&outputJSON, "json", false, "output JSON instead of human-readable format")
		viper.BindPFlags(KeysCmd.Flags())
	*/
}
