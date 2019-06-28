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

	"github.com/mosaicnetworks/monetd/cmd/monetcli/commands/network"

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
	myAddress        string
)

type peer struct {
	NetAddr   string
	PubKeyHex string
	Moniker   string
}

type peers []*peer

type copyFile struct {
	SourceFile string
	TargetFile string
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

	myAddress = key.Address.String()

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

	err := common.CreateNewConfig(testConfigDir)
	if err != nil {
		common.MessageWithType(common.MsgError, "Error creating cli network config", err)
		return err
	}
	//
	// Load Peers
	common.MessageWithType(common.MsgInformation, "Getting peers.json")
	url := CfgServer + "/peersjson"
	b, err := getRequest(url)
	if err != nil {
		common.MessageWithType(common.MsgError, "Error getting peers", err)
		return err
	}
	//save peers
	err = common.WriteToFile(filepath.Join(testConfigDir, "peers.json"), string(b))
	if err != nil {
		common.MessageWithType(common.MsgError, "Error writing peers", err)
		return err
	}

	//parse peers
	//	decoder := json.NewDecoder(b)
	var peerlist peers

	common.MessageWithType(common.MsgInformation, "Unmarshalling peers.json")

	err = json.Unmarshal(b, &peerlist)
	if err != nil {
		common.MessageWithType(common.MsgError, "Error decoding peers", err)
		return err
	}

	common.MessageWithType(common.MsgInformation, "Peers list unmarshalled: ", len(peerlist), peerlist)

	for i, p := range peerlist {

		common.MessageWithType(common.MsgInformation, "Adding... ", p.Moniker)
		safeLabel := common.GetNodeSafeLabel(p.Moniker)
		common.MessageWithType(common.MsgDebug, "safe label... ", safeLabel)

		err := network.AddValidatorParamaterised(testConfigDir, p.Moniker, safeLabel, "", p.PubKeyHex, p.NetAddr, true)
		if err != nil {
			common.MessageWithType(common.MsgError, "Error adding peers to network toml", i, err)
			return err
		}
	}

	// Compile
	err = network.CompileConfigWithParam(testConfigDir)
	if err != nil {
		common.MessageWithType(common.MsgError, "Error compiling: ", err)
		return err
	}

	//Read Genesis.Json
	b, err = ioutil.ReadFile(filepath.Join(testConfigDir, "genesis.json"))
	if err != nil {
		common.MessageWithType(common.MsgError, "Cannot read genesis.json from local disk: ", err)
		return err
	}

	// Set Json
	url = CfgServer + "/setgenesisjson"
	err = sendJSON(url, b, "text/text")

	if err != nil {
		common.MessageWithType(common.MsgError, "Genesis.Json publishing error", err)
		return err
	}

	/*
		// Set Json
		url = CfgServer + "/setnetworktoml"
		err = sendJSON(url, []byte("Network Toml"), "text/text")

		if err != nil {
			common.MessageWithType(common.MsgError, "Network.toml publishing error", err)
			return err
		}
	*/

	// Set CfgServer as published
	url = CfgServer + "/publish"
	b, err = getRequest(url)
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

	err = generateMonetdToml()
	if err != nil {
		return err
	}

	/*
		common.MessageWithType(common.MsgInformation, "Getting network.toml")
		url = CfgServer + "/networktoml"
		b, err = getRequest(url)
		if err != nil {
			common.MessageWithType(common.MsgError, "Error getting peers", err)
			return err
		}

		err = common.WriteToFile(filepath.Join(testConfigDir, "network.toml"), string(b))
		if err != nil {
			common.MessageWithType(common.MsgError, "Error writing network.toml", err)
			return err
		}
	*/
	// Get Genesis JSON
	// Get Peers JSON

	common.MessageWithType(common.MsgInformation, "All files downloaded")

	// Copy stuff into place

	/*
		testnet
		├── contract0.abi *
		├── contract0.sol *
		├── genesis.json
		├── keyfile.json
		├── network.toml
		├── peers.genesis.json *
		├── peers.json
		└── pwd.txt

		* files only on publishing machine.
	*/

confirmloop:
	for {
		confirm := common.RequestSelect("Confirm Overwriting Existing Configuration", []string{"No", "Yes"}, "No")
		if confirm == "Yes" {
			break confirmloop
		}
	}

	defaultMonetConfigDir, _ := common.DefaultHomeDir(common.MonetdTomlDir)

	common.SafeRenameDir(defaultMonetConfigDir)

	newdirs := []string{
		defaultMonetConfigDir,
		filepath.Join(defaultMonetConfigDir, "babble"),
		filepath.Join(defaultMonetConfigDir, "eth"),
		filepath.Join(defaultMonetConfigDir, "eth", "keystore"),
	}

	for _, dir := range newdirs {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			common.Message("Error creating empty config folder: ", err)
			return err
		}
	}

	copyfiles := []copyFile{
		copyFile{SourceFile: filepath.Join(testConfigDir, "monetd.toml"),
			TargetFile: filepath.Join(defaultMonetConfigDir, "monetd.toml")},
		copyFile{SourceFile: filepath.Join(testConfigDir, "genesis.json"),
			TargetFile: filepath.Join(defaultMonetConfigDir, "eth", "genesis.json")},
		copyFile{SourceFile: filepath.Join(testConfigDir, "peers.json"),
			TargetFile: filepath.Join(defaultMonetConfigDir, "babble", "peers.json")},
		copyFile{SourceFile: filepath.Join(testConfigDir, "peers.json"),
			TargetFile: filepath.Join(defaultMonetConfigDir, "babble", "peers.genesis.json")},

		copyFile{SourceFile: filepath.Join(testConfigDir, "pwd.txt"),
			TargetFile: filepath.Join(defaultMonetConfigDir, "eth", "pwd.txt")},

		copyFile{SourceFile: filepath.Join(testConfigDir, keys.DefaultKeyfile),
			TargetFile: filepath.Join(defaultMonetConfigDir, "eth", "keystore", keys.DefaultKeyfile)},

		copyFile{SourceFile: filepath.Join(testConfigDir, keys.DefaultKeyfile),
			TargetFile: filepath.Join(defaultMonetConfigDir, keys.DefaultKeyfile)},
	}

	for i, cf := range copyfiles {
		common.MessageWithType(common.MsgInformation, "Copying to ", i, cf.TargetFile)
		err = common.CopyFileContents(cf.SourceFile, cf.TargetFile)
		if err != nil {
			common.MessageWithType(common.MsgError, "Error copying ", i, cf.TargetFile)
			return err
		}
	}

	common.MessageWithType(common.MsgInformation, "Updating evmlc config")
	err = updateEvmlcConfig()
	if err != nil {
		common.MessageWithType(common.MsgError, "Error Updating evmlc config ", err)
		return err
	}
	common.MessageWithType(common.MsgWarning, "Try running:  monetd run")

	return nil

}

func generateMonetdToml() error {

	ip := common.RequestString("Enter your ip without the port: ", getMyIP())
	defaultMonetConfigDir, _ := common.DefaultHomeDir(common.MonetdTomlDir)

	toml := `datadir = "` + defaultMonetConfigDir + `"
log = "debug"

[babble]
listen = "` + ip + `:1337"
service-listen = ":8000"
heartbeat = "500ms"
timeout = "1s"
cache-size = 50000
sync-limit = 1000
max-pool = 2

[eth]
listen = ":8080"
cache = 128
`

	err := common.WriteToFile(filepath.Join(testConfigDir, "monetd.toml"), toml)
	if err != nil {
		common.MessageWithType(common.MsgError, "Error writing monetd.toml", err)
		return err
	}

	return nil

}

func updateEvmlcConfig() error {
	defaultEVMLCConfigDir, _ := common.DefaultHomeDir(".evmlc")
	defaultMonetConfigDir, _ := common.DefaultHomeDir(common.MonetdTomlDir)
	tomlFile := filepath.Join(defaultEVMLCConfigDir, "config.toml")
	keystoreFile := filepath.Join(defaultMonetConfigDir, "eth", "keystore")

	tree, err := common.LoadToml(tomlFile)
	if err != nil {
		common.MessageWithType(common.MsgError, "Error loading toml: ", tomlFile)
	}

	tree.SetPath([]string{"storage", "keystore"}, keystoreFile)
	tree.SetPath([]string{"connection", "host"}, getMyIP())

	if myAddress != "" {
		tree.SetPath([]string{"defaults", "from"}, myAddress)
	}
	tree.SetPath([]string{"defaults", "gas"}, 100000000.0)
	tree.SetPath([]string{"defaults", "gasPrice"}, 0.0)

	err = common.SaveToml(tree, tomlFile)
	if err != nil {
		common.Message("Cannot save toml file")
		return err
	}

	return nil
}
