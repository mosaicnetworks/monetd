package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/mosaicnetworks/monetd/src/common"
)

// [{"NetAddr":"192.168.1.3:1337","PubKeyHex":"0X04541909581932BD007515661FCA68FF9607447C5ECCFD262E33A74ECA2E5DE0428BCA9261D7E1D52513F38DB5A6EE100689A1C5FD6CB605EDE45AEEBF90D4AE82","Moniker":"banquet"}]

type peer struct {
	NetAddr   string
	PubKeyHex string
	Moniker   string
}

type peers []*peer

var (
	peerlist    peers
	isPublished = false
	genesisJSON string
	networkTOML = "unused"
)

func servermain() {

	fmt.Println("Starting monetcfgsrv")
	fmt.Println(common.GetMyIP() + ":8088")

	http.HandleFunc("/", cfgHandler)

	if err := http.ListenAndServe(":8088", nil); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Started monetcfgsrv")

}

func cfgHandler(w http.ResponseWriter, r *http.Request) {

	switch r.URL.Path {
	case "/peersjson":
		outputPeers(w)
	case "/ispublished":
		outputIsPublished(w)
	case "/genesisjson":
		outputGenesisJSON(w)
	case "/networktoml":
		outputNetworkTOML(w)
	case "/addpeer":
		addPeer(w, r)
	case "/setgenesisjson":
		addGenesisJSON(w, r)
	case "/setnetworktoml":
		addNetworkTOML(w, r)
	case "/publish":
		publish(w)
	default:
		fmt.Fprintln(w, r.URL.Path)

	}

}

func publish(w http.ResponseWriter) {
	if genesisJSON == "" || networkTOML == "" || len(peerlist) < 1 {
		fmt.Fprintln(w, "false")
	} else {
		isPublished = true
		fmt.Fprintln(w, "true")
	}
}

func addGenesisJSON(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error in addGenesisJSON: ", err.Error())
		fmt.Fprint(w, "false")
		return
	}
	genesisJSON = string(b)
	fmt.Fprint(w, "true")
}

func addNetworkTOML(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error in addNetworkTOML: ", err.Error())
		fmt.Fprint(w, "false")
		return
	}
	networkTOML = string(b)
	fmt.Fprint(w, "true")
}

func addPeer(w http.ResponseWriter, r *http.Request) {

	if isPublished {
		fmt.Fprint(w, "false")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var p peer

	err := decoder.Decode(&p)
	if err != nil {
		fmt.Println("Error in addpeer: ", err.Error())
		fmt.Fprint(w, "false")
		return
	}

	peerlist = append(peerlist, &p)
	fmt.Fprint(w, "true")
}

func outputPeers(w http.ResponseWriter) {
	b, err := json.Marshal(peerlist)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Fprint(w, "false")
		return
	}
	fmt.Fprintln(w, string(b))
}

func outputGenesisJSON(w http.ResponseWriter) {
	if isPublished {
		fmt.Fprintln(w, genesisJSON)
	} else {
		fmt.Fprint(w, "false")
	}

}

func outputNetworkTOML(w http.ResponseWriter) {
	if isPublished {
		fmt.Fprintln(w, networkTOML)
	} else {
		fmt.Fprint(w, "false")
	}

}

func outputIsPublished(w http.ResponseWriter) {
	fmt.Fprintln(w, strconv.FormatBool(isPublished))
}
