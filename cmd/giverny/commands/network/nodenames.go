package network

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mosaicnetworks/monetd/src/common"
)

var nodeNamePrefix = "node"

func getNodesWithNames(srcFile string, numNodes int, numValidators int, initialIP string) ([]node, error) {
	var rtn []node
	lastDigit := 0
	IPStem := ""
	var err error

	if initialIP != "" {
		splitIP := strings.Split(initialIP, ".")

		if len(splitIP) != 4 {
			return rtn, errors.New("malformed initial IP: " + initialIP)
		}

		lastDigit, err = strconv.Atoi(splitIP[3])
		if err != nil {
			fmt.Println("lastDigit Set to Zero")
			lastDigit = 0
		} else {
			IPStem = strings.Join(splitIP[:3], ".") + "."
		}
	}

	netaddr := ""
	// if no input file specified, then we generate node0, node1 etc
	if srcFile == "" {
		for i := 0; i < numNodes; i++ {
			if IPStem != "" {
				netaddr = IPStem + strconv.Itoa(lastDigit+i)
			}

			rtn = append(rtn, node{Moniker: nodeNamePrefix + strconv.Itoa(i),
				NetAddr: netaddr, Validator: (numValidators < 1 || i < numValidators),
				Tokens: defaultTokens, Address: "", PubKeyHex: ""})
		}
		return rtn, nil
	}

	// Read file line by line.

	file, err := os.Open(srcFile)

	if err != nil {
		common.ErrorMessage("failed opening file: ", err)
		return rtn, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	i := 1
	for scanner.Scan() {
		if IPStem != "" {
			netaddr = IPStem + strconv.Itoa(lastDigit+i-1)
		}

		moniker := strings.TrimSpace(scanner.Text())
		if moniker == "" {
			continue
		} // Ignore blank lines
		safeMoniker := common.GetNodeSafeLabel(moniker)
		if moniker != safeMoniker {
			return rtn, errors.New("node name " + moniker + " contains invalid characters")
		}

		rtn = append(rtn, node{Moniker: scanner.Text(),
			NetAddr: netaddr, Validator: (numValidators < 1 || i <= numValidators),
			Tokens: defaultTokens, Address: "", PubKeyHex: ""})

		if i >= numNodes {
			break
		}
		i++
	}

	file.Close()
	return rtn, nil
}
