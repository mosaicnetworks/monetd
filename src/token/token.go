//Package token provides library functions for manipulating token balances
//
// Taking inspiration from the SI units, we have suitable multiples
//
//  1 000 000 000 000 000 000 000 000	yotta	(Y)	10^24
//  1 000 000 000 000 000 000 000		zetta	(Z)	10^21
//  1 000 000 000 000 000 000			exa		(E)	10^18
//  1 000 000 000 000 000				peta	(P)	10^15
//  1 000 000 000 000					tera	(T)	10^12
//  1 000 000 000						giga	(G)	10^9
//  1 000 000							mega	(M)	10^6
//  1 000								kilo	(K)	10^3
//
//  NB we use a capital K for kilo, so all letters are capital.
//
//  Capital E is treated as an exponential, lower case e is a hex number.

package token

import (
	"fmt"
	"strconv"
	"strings"
)

const tokenLetters = "KMGTPEZY"

func ExpandBalance(input string) (string, error) {

	cleanInput := strings.TrimSpace(input)

	token := cleanInput[len(cleanInput)-1:]
	tokenPower := (strings.Index(tokenLetters, token) + 1) * 3

	if tokenPower == 0 {
		return input, nil
	} // Not found  ( -1 + 1 ) * 3  =  0

	last := len(cleanInput) - 1
	cleanInput = cleanInput[:last] // trim token from string.

	idx := strings.Index(cleanInput, ".")
	if idx >= 0 {
		pre := cleanInput[:idx]
		fix := cleanInput[idx+1:]
		tokenPower -= len(fix)

		cleanInput = pre + fix
	}

	format := "%0" + strconv.Itoa(tokenPower) + "d"
	return cleanInput + fmt.Sprintf(format, 0), nil

}

//  1 000 000 000 000 000 000 000 000	yotta (Y)	1024
//  1 000 000 000 000 000 000 000		zetta (Z)	1021
//  1 000 000 000 000 000 000			exa (E)		1018
//  1 000 000 000 000 000				peta (P)	1015
//  1 000 000 000 000					tera (T)	1012
//  1 000 000 000						giga (G)	109
//  1 000 000							mega (M)	106
//  1 000								kilo (k)	103
