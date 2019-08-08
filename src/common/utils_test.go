package common_test

import (
	"strconv"
	"testing"

	"github.com/mosaicnetworks/monetd/src/common"
)

type testRecord struct {
	input  string
	output bool
}

func TestCheckMoniker(t *testing.T) {

	var tests = []testRecord{
		testRecord{input: "abcdef", output: true},
		testRecord{input: "ABCDEF", output: true},
		testRecord{input: "abc.def", output: false},
		testRecord{input: "ab_cdef", output: true},
		testRecord{input: "ab__ef", output: true},
		testRecord{input: "abcd54ef", output: true},
		testRecord{input: "1234", output: true},
		testRecord{input: "abcd ef", output: false},
		testRecord{input: "abcdef ", output: false},
		testRecord{input: " abcdef", output: false},
		testRecord{input: "a!bcdef", output: false},
	}

	for _, test := range tests {
		ret := common.CheckMoniker(test.input)
		if ret != test.output {
			t.Errorf("\nWrong Answer: %s\nGot: %s\nExpected: %s\n",
				test.input, strconv.FormatBool(ret), strconv.FormatBool(test.output))
		} else {
			t.Logf("%s => %s", test.input, strconv.FormatBool(test.output))
		}
	}
}
