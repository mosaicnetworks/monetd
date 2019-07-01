package network

import (
	"fmt"

	"github.com/mosaicnetworks/monetd/src/common"
)

//PeersWizard implements interactive management of the Peers in the MonetCli Config
func PeersWizard(configDir string) error {

	var err error
	common.Message("Peers")

	tree, err := common.LoadTomlConfig(configDir)
	if err != nil {
		return err
	}

	peers := GetPeersLabelsList(tree)
	peers = append(peers, common.WizardExit)

	refreshPeerList := false
peerloop:
	for {

		if refreshPeerList {
			tree, err := common.LoadTomlConfig(configDir)
			if err != nil {
				return err
			}

			peers := GetPeersLabelsList(tree)
			peers = append(peers, common.WizardExit)
		}

		refreshPeerList = false

		selectedPeer := common.RequestSelect("Choose a peer", peers, "")
		if selectedPeer == common.WizardExit {
			break peerloop
		}
	actionloop:
		for {
			_ = viewPeer(configDir, selectedPeer)

			action := common.RequestSelect("Choose an Action: ", []string{common.WizardEdit, common.WizardDelete, common.WizardExit}, "")

			switch action {
			case common.WizardEdit:
				err = editPeer(configDir, selectedPeer)
				refreshPeerList = true
			case common.WizardDelete:
				err = deletePeer(configDir, selectedPeer)
				refreshPeerList = true
				break actionloop
			case common.WizardExit:
				break actionloop

			}

			if err != nil {
				return err
			}
		}
	}

	fmt.Printf("%v \n", peers)

	return nil
}

func viewPeer(configDir string, peername string) error {
	tree, err := common.LoadTomlConfig(configDir)
	if err != nil {
		return err
	}

	peer := tree.GetPath([]string{"validators", peername})

	common.MessageWithType(common.MsgInformation, peer)

	return nil
	//TODO add new View Peer code
}

func editPeer(configDir string, peername string) error {
	return nil
	//TODO add new edit peer code
}

func deletePeer(configDir string, peername string) error {
	return nil
	//TODO add new delete peer code
}
