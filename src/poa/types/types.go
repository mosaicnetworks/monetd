// Package types contains common types used in multiple poa packages
package types

//PeerRecord is Peers.json format with suitable export config
type PeerRecord struct {
	NetAddr   string `json:"NetAddr"`
	PubKeyHex string `json:"PubKeyHex"`
	Moniker   string `json:"Moniker"`
}

//PeerRecordList is a slice of PeerRecord and is used to pass lists of peers
//within the POA packages.
type PeerRecordList []*PeerRecord
