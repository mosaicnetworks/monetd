package trans

import (
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/spf13/cobra"
)

//newDecodeCmd returns the command that creates a Ethereum keyfile
func newDecodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decode [payload]",
		Short: "decode a payload",
		Long: `
Decode a signed transaction payload
`,
		Args: cobra.ExactArgs(1),
		RunE: decodeTrans,
	}

	return cmd
}

func decodeTrans(cmd *cobra.Command, args []string) error {
	payload := args[0]

	common.DebugMessage(payload)

	return decodeTransactionPayload(payload)
}

func decodeTransactionPayload(rawTx string) error {

	//	rawTx := "f86d8202b28477359400825208944592d8f8d7b001e72cb26a73e4fa1806a51ac79d880de0b6b3a7640000802ca05924bde7ef10aa88db9c66dd4f5fb16b46dff2319b9968be983118b57bb50562a001b24b31010004f13d9a26b320845257a6cfc2bf819a3d55e3fc86263c5f0772"

	rawTx = strings.TrimPrefix(strings.TrimPrefix(rawTx, "0x"), "0X")
	var tx *types.Transaction

	rawTxBytes, err := hex.DecodeString(rawTx)
	if err != nil {
		return err
	}
	rlp.DecodeBytes(rawTxBytes, &tx)

	var jsonData []byte

	jsonData, err = json.Marshal(tx)
	if err != nil {
		return err
	}

	common.InfoMessage(string(jsonData))

	return nil
}
