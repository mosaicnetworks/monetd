package network

import (
	"fmt"

	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newLocationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "location",
		Short: "show the location of the configuration files",
		Long: `
giverny network location
		`,
		Args: cobra.ArbitraryArgs,
		RunE: networkLocation,
	}

	addLocationFlags(cmd)

	return cmd
}

func addLocationFlags(cmd *cobra.Command) {
	//	cmd.Flags().StringVar(&addressParam, "address", addressParam, "IP/hostname of this node")
	//	cmd.Flags().StringVar(&passwordFile, "passfile", "", "file containing the passphrase")
	viper.BindPFlags(cmd.Flags())
}

func networkLocation(cmd *cobra.Command, args []string) error {

	fmt.Println(configuration.GivernyConfigDir)

	return nil
}
