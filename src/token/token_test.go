package token_test

import (
	"testing"

	"github.com/mosaicnetworks/monetd/src/token"
)

type testRecord struct {
	input  string
	output string
}

//  1 000 000 000 000 000 000 000 000	yotta (Y)	1024
//  1 000 000 000 000 000 000 000		zetta (Z)	1021
//  1 000 000 000 000 000 000			exa (E)		1018
//  1 000 000 000 000 000				peta (P)	1015
//  1 000 000 000 000					tera (T)	1012
//  1 000 000 000						giga (G)	109
//  1 000 000							mega (M)	106
//  1 000								kilo (k)	103

func TestExpandBalance(t *testing.T) {

	var tests = []testRecord{
		testRecord{input: "1K", output: "1000"},
		testRecord{input: "1.2M", output: "1200000"},
		testRecord{input: "1.23G", output: "1230000000"},
		testRecord{input: "1.2T", output: "1200000000000"},
		testRecord{input: "1.2P", output: "1200000000000000"},
		testRecord{input: "1.2E", output: "1200000000000000000"},
		testRecord{input: "1.2Z", output: "1200000000000000000000"},
		testRecord{input: "1.2Y", output: "1200000000000000000000000"},
		testRecord{input: "0.2K", output: "0200"},
		testRecord{input: "0x122K", output: "0x122000"},
	}

	for _, test := range tests {
		ret, _ := token.ExpandBalance(test.input)
		if ret != test.output {
			t.Errorf("\nWrong Answer: %s\nGot: %s\nExpected: %s\n", test.input, ret, test.output)
		} else {
			t.Logf("%s => %s", test.input, test.output)
		}
	}
}
