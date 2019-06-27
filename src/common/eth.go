package common

import (
	"log"
	"regexp"
)

//IsValidAddress checks if is a valid ethereum style address
func IsValidAddress(v string) bool {

	re := regexp.MustCompile("^(0[xX]){0,1}[0-9a-fA-F]{40}$")
	return re.MatchString(v)
}

//GetNodeSafeLabel converts a free format string into a node label friendly format
//Anything other than an alphanumeric is converted to _
func GetNodeSafeLabel(moniker string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	return reg.ReplaceAllString(moniker, "_")
}
