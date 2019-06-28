package testnet

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mosaicnetworks/monetd/cmd/monetcli/commands/keys"

	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//NetworkCmd controls network configuration
var (
	publishTarget    string
	monetConfigDir   string
	networkConfigDir string
	jumpToPublish    bool
	testConfigDir    string
	CfgServer        string
)

type peer struct {
	NetAddr   string
	PubKeyHex string
	Moniker   string
}

func init() {
	//	defaultConfigDir, _ := common.DefaultHomeDir(common.MonetcliTomlDir)
	//	defaultMonetConfigDir, _ := common.DefaultHomeDir(common.MonetdTomlDir)
}

//NewCheckCmd defines the CLI command config check
func NewTestNetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testnet",
		Short: "manage monetd testnets",
		Long: `
TestNet Config.`,
		Args: cobra.ArbitraryArgs,
		RunE: testnetCmd,
	}

	cmd.PersistentFlags().BoolVarP(&jumpToPublish, "publish", "p", false, "jump straight to polling for a configuration")
	viper.BindPFlags(cmd.Flags())

	return cmd
}

func testnetCmd(cmd *cobra.Command, args []string) error {

	return testNetWizard()
}

func testNetWizard() error {

	// Always request the server

	defaultConfigDir, _ := common.DefaultHomeDir(common.MonetcliTomlDir)
	testConfigDir = filepath.Join(defaultConfigDir, "testnet")

	if common.CheckIfExists(testConfigDir) {
		if !jumpToPublish {
			common.MessageWithType(common.MsgWarning, "This is a destructive operation. Remove/rename the following folder to proceed.")
			common.MessageWithType(common.MsgInformation, testConfigDir)
			return nil
		}
	} else {
		err := os.MkdirAll(testConfigDir, os.ModePerm)
		if err != nil {
			common.MessageWithType(common.MsgError, "Error creating config folder: ", testConfigDir)
			return err
		}
	}

cfgserverloop:
	for {

		CfgServer = common.RequestString("Configuration Server:", "http://localhost:8088")

		// Hacky, but an empty server string means that we want to abort
		if CfgServer == "" {
			return nil
		}

		if checkIsLiveServer(CfgServer) {
			break cfgserverloop
		}
	}

	if !jumpToPublish {
		err := enterParams()
		if err != nil {
			return err
		}
	}

	err := publishWizard()
	if err != nil {
		return err
	}

	return nil
}

func enterParams() error {
	var password, moniker, ip, pubkey string

	// request password
passwordloop:
	for {
		password = common.RequestPassword("Enter Keystore Password: ", "")
		password2 := common.RequestPassword("Confirm Keystore Password: ", "")

		if password == password2 {
			break passwordloop
		}
	}

	passwordFile := filepath.Join(testConfigDir, "pwd.txt")

	err := common.WriteToFile(passwordFile, password)
	if err != nil {
		common.MessageWithType(common.MsgError, "Error saving password: ", err)
		return err

	}

	// request name
	moniker = common.RequestString("Enter your moniker: ", "")

	// confirm your ipS
	ip = common.RequestString("Enter your ip without the port: ", getMyIP())

	// generate key

	keyfilepath := filepath.Join(testConfigDir, keys.DefaultKeyfile)
	key, err := keys.GenerateKeyPair(keyfilepath, passwordFile)
	if err != nil {
		common.MessageWithType(common.MsgError, "Error generating key: ", err)
		return err
	}

	common.MessageWithType(common.MsgInformation, "Building Data to push to Configuration Server")

	pubkey = hex.EncodeToString(
		crypto.FromECDSAPub(&key.PrivateKey.PublicKey))

	common.MessageWithType(common.MsgInformation, "Moniker  : ", moniker)
	common.MessageWithType(common.MsgInformation, "IP       : ", ip)
	common.MessageWithType(common.MsgInformation, "Pub Key  : ", pubkey)
	common.MessageWithType(common.MsgInformation, "Address  : ", key.Address.String())

	peer := peer{
		NetAddr:   ip + ":1337",
		PubKeyHex: "0x" + pubkey,
		Moniker:   moniker,
	}

	b, err := json.Marshal(peer)
	if err != nil {
		common.MessageWithType(common.MsgError, "Error Marshalling Peer JSON: ", peer)
		return err
	}

	url := CfgServer + "/addpeer"

	common.MessageWithType(common.MsgInformation, "URL      : ", url)

	err = sendJSON(url, b, "application/json")
	if err != nil {
		return err
	}
	return nil

}

func sendJSON(url string, b []byte, contenttype string) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("X-Custom-Header", "monetcfgsrv")
	req.Header.Set("Content-Type", contenttype)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	return nil
}

func checkIsLiveServer(addr string) bool {
	//TODO load a page from server to confirm it is live
	return true
}

func getMyIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {

				return ipnet.IP.String()
			}
		}
	}

	return ""
}

func publishWizard() error {

	const (
		CheckIfPublished = "Check if published"
		Publish          = "Publish, no more initial peers will be allowed to be added"
		Exit             = "Exit"
	)

publishloop:
	for {

		selectedOption := common.RequestSelect("Choose your action", []string{CheckIfPublished, Publish, Exit}, CheckIfPublished)
		switch selectedOption {
		case Exit:
			return nil
		case CheckIfPublished:
			if checkIfPublished() {
				break publishloop
			}
		case Publish:
			err := publish()
			if err != nil {
				return nil
			}
			break publishloop
		}

	}

	common.MessageWithType(common.MsgInformation, "Configuration has been published.")

	err := buildConfig()
	if err != nil {
		return nil
	}
	return nil
}

func checkIfPublished() bool {
	url := CfgServer + "/ispublished"

	b, err := getRequest(url)
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(b)) == "true"
}

func publish() error {

	// Set Json
	url := CfgServer + "/setgenesisjson"
	err := sendJSON(url, []byte("Genesis!!"), "text/text")

	if err != nil {
		common.MessageWithType(common.MsgError, "Genesis.Json publishing error", err)
		return err
	}

	url = CfgServer + "/publish"
	b, err := getRequest(url)
	if err != nil {
		common.MessageWithType(common.MsgError, "Publishing error", err)
		return err
	}

	common.MessageWithType(common.MsgInformation, "Publish result: "+string(b))
	return nil
}

func getRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return bytes, nil
}

func buildConfig() error {

	common.MessageWithType(common.MsgInformation, "Getting peers.json")
	url := CfgServer + "/peersjson"
	b, err := getRequest(url)
	if err != nil {
		common.MessageWithType(common.MsgError, "Error getting peers", err)
		return err
	}

	err = common.WriteToFile(filepath.Join(testConfigDir, "peers.json"), string(b))
	if err != nil {
		common.MessageWithType(common.MsgError, "Error writing peers", err)
		return err
	}

	common.MessageWithType(common.MsgInformation, "Getting genesis.json")
	url = CfgServer + "/genesisjson"
	b, err = getRequest(url)
	if err != nil {
		common.MessageWithType(common.MsgError, "Error getting genesis json", err)
		return err
	}

	err = common.WriteToFile(filepath.Join(testConfigDir, "genesis.json"), string(b))
	if err != nil {
		common.MessageWithType(common.MsgError, "Error writing peers", err)
		return err
	}

	// Get Genesis JSON
	// Get Peers JSON

	return nil
}
