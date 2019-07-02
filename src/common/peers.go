package common

import (
	"github.com/mosaicnetworks/babble/src/peers"
)

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
