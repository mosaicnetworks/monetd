package network

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/mosaicnetworks/monetd/src/files"

	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"
	mconfiguration "github.com/mosaicnetworks/monetd/src/configuration"
)

func buildZip(configDir string, networkName, nodeName string) error {
	sourceDir := filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir, networkName)
	monetdtoml := filepath.Join(sourceDir, mconfiguration.MonetTomlFile)
	peersjson := filepath.Join(sourceDir, mconfiguration.PeersJSON)
	peersgenesisjson := filepath.Join(sourceDir, mconfiguration.PeersGenesisJSON)
	genesisjson := filepath.Join(sourceDir, mconfiguration.GenesisJSON)
	acctjson := filepath.Join(sourceDir, givernyKeystoreDir, nodeName+".json")
	passphrase := filepath.Join(sourceDir, givernyKeystoreDir, nodeName+".txt")

	filesList := []string{
		monetdtoml,
		peersjson,
		peersgenesisjson,
		genesisjson,
	}

	if checkIfKeyfileContainsPrivateKey(acctjson) {
		filesList = append(filesList, acctjson)
		filesList = append(filesList, passphrase)
	}

	for _, f := range filesList {
		if !files.CheckIfExists(f) {
			return errors.New("missing file, " + f + ", configuration is incomplete and cannot be published")
		}
	}

	outputDir := filepath.Join(configuration.GivernyConfigDir, configuration.GivernyExportDir)
	err := files.CreateDirsIfNotExists([]string{outputDir})
	if err != nil {
		return err
	}

	outputFile := filepath.Join(outputDir, networkName+"_"+nodeName+".zip")

	common.InfoMessage("Writing to " + outputFile)
	err = zipFiles(outputFile, filesList, false)
	return err
}

func zipFiles(filename string, files []string, perservePath bool) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {

		zipfile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer zipfile.Close()

		// Get the file information
		info, err := zipfile.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if perservePath {
			header.Name = file
		}

		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if _, err = io.Copy(writer, zipfile); err != nil {
			return err
		}
	}
	return nil
}

func checkIfKeyfileContainsPrivateKey(keyfilename string) bool {

	common.DebugMessage("Checking " + keyfilename)
	jsonFile, err := os.Open(keyfilename)
	if err != nil {
		common.ErrorMessage("Error opening " + keyfilename)
		common.ErrorMessage(err.Error())

		return false
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	if _, ok := result["crypto"]; ok {
		return true
	}

	return false
}
