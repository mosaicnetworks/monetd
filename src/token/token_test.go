package token

import "testing"

type testRecord struct {
	input  string
	output string
}

func TestExpandBalance(t *testing.T) {

	var tests = []testRecord{
		testRecord{input: "1K", output: "1000"},
		testRecord{input: "1.2M", output: "1200000"},
	}

	for _, test := range tests {
		ret, _ := ExpandBalance(test.input)
		if ret != test.output {
			t.Errorf("Wrong Answer: %s\nGot: %s\nExpected: %s\n", test.input, ret, test.output)
		}
	}

}
