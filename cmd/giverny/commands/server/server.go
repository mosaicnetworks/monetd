package server

import (
	"github.com/spf13/cobra"
)

const (
	//ServerDir is the subfolder of the configuration folder that contains files for
	//the server process. It is currently set to giverny which is shared with all
	//giverney subcommands
	ServerDir = "server"
	//ServerPIDFile is the file name where the process ID for the server process is written.
	ServerPIDFile = "server.pid"
	//ServerLogFile is the log file within ServerDir for the server process
	ServerLogFile = "server.log"
)

//ServerCmd is the CLI command for the giverny server
var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "giverny server",
	Long: `Server
	
The giverny server is a simple REST server to facilitate the sharing of
Monet configurations prior to instantiation of the node. The server 
is started with the run subcommand for foreground operation, start for background
operation.`,

	TraverseChildren: true,
}

func init() {
	//Subcommands
	ServerCmd.AddCommand(
		newStartCmd(),
		newStopCmd(),
		newStatusCmd(),
	)

}
