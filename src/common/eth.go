package common

import "regexp"

//IsValidAddress checks if is a valid ethereum style address
func IsValidAddress(v string) bool {

	re := regexp.MustCompile("^(0[xX]){0,1}[0-9a-fA-F]{40}$")
	return re.MatchString(v)
}
