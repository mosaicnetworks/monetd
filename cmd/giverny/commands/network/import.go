package network

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mosaicnetworks/monetd/src/config"
	mconfig "github.com/mosaicnetworks/monetd/src/configuration"

	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"
	"github.com/mosaicnetworks/monetd/src/files"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//system flags
var (
	useExportDir = false
	srcdir       = ""
	srvAddress   = ""
)

func newImportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import [network] [node]",
		Short: "import the configuration for a node on the named network",
		Long: `
giverny network import
		`,
		Args: cobra.ExactArgs(2),
		RunE: networkImport,
	}

	addImportFlags(cmd)

	return cmd
}

func addImportFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&useExportDir, "from-exports", useExportDir, "source config from exports directory")
	cmd.Flags().StringVar(&srvAddress, "server", srvAddress, "giverny server address")
	cmd.Flags().StringVar(&srcdir, "dir", srcdir, "dir containing export folder")

	viper.BindPFlags(cmd.Flags())
}

func networkImport(cmd *cobra.Command, args []string) error {

	networkName := args[0]
	nodeName := args[1]

	var zipName = ""

	if useExportDir {
		zipName = filepath.Join(configuration.GivernyConfigDir, configuration.GivernyExportDir, networkName+"_"+nodeName+".zip")
	} else {
		if srcdir != "" {
			zipName = filepath.Join(srcdir, networkName+"_"+nodeName+".zip")
		} else {
			if srvAddress != "" {

				if !strings.Contains(srvAddress, ":") {
					srvAddress += ":" + configuration.GivernyServerPort
				}

				tmpDir := filepath.Join(configuration.GivernyConfigDir, givernyTmpDir)
				files.CreateDirsIfNotExists([]string{tmpDir})
				zipName = filepath.Join(tmpDir, networkName+"_"+nodeName+".zip")
				url := "http://" + srvAddress + "/import/" + networkName + "/" + nodeName
				err := downloadFile(zipName, url)
				if err != nil {
					return err
				}
			} else {
				return errors.New("you must specify exactly one of --from-exports, --server or --dir")
			}
		}
	}

	if !files.CheckIfExists(zipName) {
		return errors.New("cannot read file " + zipName)
	}

	return config.ImportZip(zipName, mconfig.Global.DataDir)

}

func downloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
