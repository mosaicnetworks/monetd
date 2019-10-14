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

type testIPRecord struct {
	input  string
	input2 bool
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

func TestCheckIP(t *testing.T) {

	var tests = []testIPRecord{
		testIPRecord{input: "127.0.0.1", input2: false, output: true},
		testIPRecord{input: "10.0.0.1:1337", input2: false, output: true},
		testIPRecord{input: "127.0.0.1", input2: true, output: true},
		testIPRecord{input: ":8080", input2: false, output: true},
		testIPRecord{input: ":8080", input2: true, output: false},
		testIPRecord{input: "8.8.8.8", input2: true, output: false},
		testIPRecord{input: "192.168.1.12", input2: true, output: true},
	}

	for _, test := range tests {
		ret := common.CheckIP(test.input, test.input2)
		if ret != test.output {
			t.Errorf("\nWrong Answer: %s\nGot: %s\nExpected: %s\n",
				test.input, strconv.FormatBool(ret), strconv.FormatBool(test.output))
		} else {
			t.Logf("%s %t => %s", test.input, test.input2, strconv.FormatBool(test.output))
		}
	}
}
