package common

import (
	"github.com/mosaicnetworks/babble/src/peers"
)

//CheckPeersAddress verifies that a given peer is in the peers file
//In particular that the configured babble listener is in exactly the same form
//within the peers file
func CheckPeersAddress(peerfile string, babbleListen string) (bool, error) {
	currentPeers, err := peers.NewJSONPeerSet(peerfile, true).PeerSet()
	if err != nil {
		MessageWithType(MsgError, "Error loading peers.json: ", err)
		return false, err
	}

	bMatch := false
	for _, p := range currentPeers.Peers {
		if p.NetAddr == babbleListen {
			MessageWithType(MsgInformation, "Babble Gossip Endpoint is in the Peers File: ", p.NetAddr)
			bMatch = true
		} else {
			MessageWithType(MsgDebug, "Non matching peer: ", p.NetAddr)
		}
	}

	return bMatch, nil

}
