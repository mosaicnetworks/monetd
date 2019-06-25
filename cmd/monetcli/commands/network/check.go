package network

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func newCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "check configuration",
		Long: `
Check configuration.`,
		RunE: checkconfig,
	}
	return cmd
}

func checkconfig(cmd *cobra.Command, args []string) error {

	err := loadConfig()

	if err != nil {
		fmt.Println("Cannot load configuration: ", err)
		return err
	}

	err = networkViper.Unmarshal(&config)
	if err != nil {
		fmt.Println("Error loading configuration: ", err)
	}

	//	message("Loaded configuration: ", config)

	if verboseLogging {
		fmt.Printf("%+v\n", config)
	}

	if addr := networkViper.GetString("poa.contractaddress"); isValidAddress(addr) {
		message("poa.contractaddress is a valid address: ", addr)
	} else {
		fmt.Println("Invalid address: ", "\""+addr+"\"")
		return errors.New("poa.contractaddress is not a valid address")
	}

	//	message("Address: ", networkViper.GetString("poa.contractaddress"))
	// networkViper.Debug()

	fmt.Println("All checks passed")
	return nil
}