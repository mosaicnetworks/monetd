package common_test

import (
	"testing"

	"github.com/mosaicnetworks/monetd/src/common"
)

func TestGetMyIP(t *testing.T) {

	result := common.GetMyIP()

	if result == "" {
		t.Error("\nFail: GetMyIP returned an empty string\n")
	} else {
		t.Logf(" => %s", result)
	}

}
