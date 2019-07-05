package network

import (
	"strings"

	"github.com/mosaicnetworks/monetd/src/common"
)

//PeersWizard implements interactive management of the Peers in the MonetCli Config
func PeersWizard(configDir string) error {

	var err error

	tree, err := common.LoadTomlConfig(configDir)
	if err != nil {
		return err
	}

	peers := GetPeersLabelsList(tree)
	peers = append(peers, common.WizardExit)

peerloop:
	for {
		common.BannerTitle("peers")

		tree, err := common.LoadTomlConfig(configDir)
		if err != nil {
			return err
		}

		peers := GetPeersLabelsList(tree)
		peers = append(peers, common.WizardExit)

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
				selectedPeer, err = editPeer(configDir, selectedPeer)
			case common.WizardDelete:
				err = deletePeer(configDir, selectedPeer)
				break actionloop
			case common.WizardExit:
				break actionloop

			}

			if err != nil {
				return err
			}
		}
		common.ContinuePrompt()
	}

	//	fmt.Printf("%v \n", peers)

	return nil
}

func viewPeer(configDir string, peername string) error {
	tree, err := common.LoadTomlConfig(configDir)
	if err != nil {
		return err
	}

	peer := tree.GetPath([]string{"validators", peername})

	//Enhancement: reformat
	common.MessageWithType(common.MsgInformation, peer)

	return nil
}

func editPeer(configDir string, peername string) (string, error) {
	tree, err := common.LoadTomlConfig(configDir)
	if err != nil {
		return peername, err
	}

	peer := tree.GetPath([]string{"validators", peername})

	//Enhancement: reformat
	common.MessageWithType(common.MsgInformation, peer)

	netaddr := tree.GetPath([]string{"validators", peername, "ip"}).(string)
	pubkey := tree.GetPath([]string{"validators", peername, "pubkey"}).(string)
	moniker := tree.GetPath([]string{"validators", peername, "moniker"}).(string)

	netaddr = common.RequestString("Net Address", netaddr)
	pubkey = common.RequestString("Public Key", pubkey)
	moniker = common.RequestString("Moniker", moniker)

	confirm := common.RequestSelect("Save Peer Updates: ", []string{"No", "Yes"}, "No")
	if confirm == "No" {
		return peername, nil
	}

	// set label
	// set address
	// set validator

	err = tree.DeletePath([]string{"validators", peername})
	if err != nil {
		return peername, err
	}

	labelsafe := common.GetNodeSafeLabel(moniker)
	addr, err := common.PublicKeyHexToAddressHex(strings.TrimPrefix(strings.ToLower(pubkey), "0x"))
	if err != nil {
		message("invalid pubkey to address conversion: ", pubkey)
		return peername, err
	}

	tree.SetPath([]string{"validators", labelsafe, "ip"}, netaddr)
	tree.SetPath([]string{"validators", labelsafe, "pubkey"}, pubkey)
	tree.SetPath([]string{"validators", labelsafe, "moniker"}, moniker)
	tree.SetPath([]string{"validators", labelsafe, "label"}, labelsafe)
	tree.SetPath([]string{"validators", labelsafe, "address"}, addr)
	tree.SetPath([]string{"validators", labelsafe, "validator"}, true)

	err = common.SaveTomlConfig(configDir, tree)
	if err != nil {
		common.MessageWithType(common.MsgDebug, "Cannot save TOML file")
		return peername, err
	}
	return labelsafe, nil
}

func deletePeer(configDir string, peername string) error {
	tree, err := common.LoadTomlConfig(configDir)
	if err != nil {
		return err
	}

	peer := tree.GetPath([]string{"validators", peername})

	//Enhancement: reformat
	common.MessageWithType(common.MsgInformation, peer)

	confirm := common.RequestSelect("Delete Peer: ", []string{"No", "Yes"}, "No")
	if confirm == "No" {
		return nil
	}

	err = tree.DeletePath([]string{"validators", peername})
	if err != nil {
		return err
	}

	err = common.SaveTomlConfig(configDir, tree)
	if err != nil {
		common.MessageWithType(common.MsgDebug, "Cannot save TOML file")
		return err
	}

	return nil
}
