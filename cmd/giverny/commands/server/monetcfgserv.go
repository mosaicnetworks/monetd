package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"
	"github.com/mosaicnetworks/monetd/src/files"

	"github.com/gorilla/mux"

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

	log.Println("Starting monetcfgsrv")
	log.Println(common.GetMyIP() + ":8088")

	r := mux.NewRouter()
	//	http.HandleFunc("/", cfgHandler)
	r.HandleFunc("/peersjson", outputPeers)
	r.HandleFunc("/ispublished", outputIsPublished)
	r.HandleFunc("/genesisjson", outputGenesisJSON)
	r.HandleFunc("/networktoml", outputNetworkTOML)
	r.HandleFunc("/addpeer", addPeer)
	r.HandleFunc("/setgenesisjson", addGenesisJSON)
	r.HandleFunc("/import/{network}/{node}", importHandler)
	r.HandleFunc("/setnetworktoml", addNetworkTOML)
	r.HandleFunc("/publish", publish)

	if err := http.ListenAndServe(":8088", r); err != nil {
		log.Fatal(err)
	}
	log.Println("Started monetcfgsrv")

}

/*
func cfgHandler(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

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
	case "/import":
		importHandler(w, r)
	case "/setnetworktoml":
		addNetworkTOML(w, r)
	case "/publish":
		publish(w)
	default:
		fmt.Fprintln(w, r.URL.Path)

	}

	// log request by who(IP address)
	requesterIP := r.RemoteAddr

	log.Printf(
		"%s\t\t%s\t\t%s\t\t%v",
		r.Method,
		r.RequestURI,
		requesterIP,
		time.Since(start),
	)

}
*/

func importHandler(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	vars := mux.Vars(r)
	network := vars["network"]
	node := vars["node"]

	shortname := network + "_" + node + ".zip"
	zipname := filepath.Join(configuration.GivernyConfigDir, configuration.GivernyExportDir, shortname)

	if !files.CheckIfExists(zipname) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Not found %v %v\n", network, node)

		requesterIP := r.RemoteAddr

		log.Printf(
			"%s\t\t%s\t\t%s\t\t%v",
			r.Method,
			r.RequestURI,
			requesterIP,
			time.Since(start),
		)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename="+shortname)
	http.ServeFile(w, r, zipname)

	requesterIP := r.RemoteAddr

	log.Printf(
		"%s\t\t%s\t\t%s\t\t%v",
		r.Method,
		r.RequestURI,
		requesterIP,
		time.Since(start),
	)
}

func publish(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	if genesisJSON == "" || networkTOML == "" || len(peerlist) < 1 {
		fmt.Fprintln(w, "false")
	} else {
		isPublished = true
		fmt.Fprintln(w, "true")
	}
	// log request by who(IP address)
	requesterIP := r.RemoteAddr

	log.Printf(
		"%s\t\t%s\t\t%s\t\t%v",
		r.Method,
		r.RequestURI,
		requesterIP,
		time.Since(start),
	)
}

func addGenesisJSON(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error in addGenesisJSON: ", err.Error())
		fmt.Fprint(w, "false")
		return
	}
	genesisJSON = string(b)
	fmt.Fprint(w, "true")
	// log request by who(IP address)
	requesterIP := r.RemoteAddr

	log.Printf(
		"%s\t\t%s\t\t%s\t\t%v",
		r.Method,
		r.RequestURI,
		requesterIP,
		time.Since(start),
	)
}

func addNetworkTOML(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error in addNetworkTOML: ", err.Error())
		fmt.Fprint(w, "false")
		return
	}
	networkTOML = string(b)
	fmt.Fprint(w, "true")
	// log request by who(IP address)
	requesterIP := r.RemoteAddr

	log.Printf(
		"%s\t\t%s\t\t%s\t\t%v",
		r.Method,
		r.RequestURI,
		requesterIP,
		time.Since(start),
	)
}

func addPeer(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	if isPublished {
		fmt.Fprint(w, "false")
		// log request by who(IP address)
		requesterIP := r.RemoteAddr

		log.Printf(
			"%s\t\t%s\t\t%s\t\t%v",
			r.Method,
			r.RequestURI,
			requesterIP,
			time.Since(start),
		)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var p peer

	err := decoder.Decode(&p)
	if err != nil {
		log.Println("Error in addpeer: ", err.Error())
		fmt.Fprint(w, "false")
		// log request by who(IP address)
		requesterIP := r.RemoteAddr

		log.Printf(
			"%s\t\t%s\t\t%s\t\t%v",
			r.Method,
			r.RequestURI,
			requesterIP,
			time.Since(start),
		)
		return
	}

	peerlist = append(peerlist, &p)
	fmt.Fprint(w, "true")
	// log request by who(IP address)
	requesterIP := r.RemoteAddr

	log.Printf(
		"%s\t\t%s\t\t%s\t\t%v",
		r.Method,
		r.RequestURI,
		requesterIP,
		time.Since(start),
	)
}

func outputPeers(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	b, err := json.Marshal(peerlist)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprint(w, "false")
		return
	}
	fmt.Fprintln(w, string(b))

	// log request by who(IP address)
	requesterIP := r.RemoteAddr

	log.Printf(
		"%s\t\t%s\t\t%s\t\t%v",
		r.Method,
		r.RequestURI,
		requesterIP,
		time.Since(start),
	)
}

func outputGenesisJSON(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	if isPublished {
		fmt.Fprintln(w, genesisJSON)
	} else {
		fmt.Fprint(w, "false")
	}
	// log request by who(IP address)
	requesterIP := r.RemoteAddr

	log.Printf(
		"%s\t\t%s\t\t%s\t\t%v",
		r.Method,
		r.RequestURI,
		requesterIP,
		time.Since(start),
	)

}

func outputNetworkTOML(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	if isPublished {
		fmt.Fprintln(w, networkTOML)
	} else {
		fmt.Fprint(w, "false")
	}

	// log request by who(IP address)
	requesterIP := r.RemoteAddr

	log.Printf(
		"%s\t\t%s\t\t%s\t\t%v",
		r.Method,
		r.RequestURI,
		requesterIP,
		time.Since(start),
	)
}

func outputIsPublished(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	fmt.Fprintln(w, strconv.FormatBool(isPublished))
	// log request by who(IP address)
	requesterIP := r.RemoteAddr

	log.Printf(
		"%s\t\t%s\t\t%s\t\t%v",
		r.Method,
		r.RequestURI,
		requesterIP,
		time.Since(start),
	)
}
