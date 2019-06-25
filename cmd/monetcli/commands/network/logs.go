package network

import "fmt"

func message(a ...interface{}) (n int, err error) {
	if verboseLogging {
		n, err = fmt.Println(a...)
		return n, err
	}

	return 0, nil
}
