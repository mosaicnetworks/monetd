package network

import (
	"github.com/mosaicnetworks/monetd/src/common"
)

//Wrapper code to use common version of message.
//Ongoing call common.Message directly
func message(a ...interface{}) (n int, err error) {
	return common.Message(a)
}
