package common

import (
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

// MonikerRegexp defines the set of characters that monikers can be composed of.
var MonikerRegexp = regexp.MustCompile("^[a-zA-Z0-9_]*$")

// MustPrintJSON prints the JSON encoding of the given object and
// exits the program with an error message when the marshaling fails.
func MustPrintJSON(jsonObject interface{}) error {
	str, err := json.MarshalIndent(jsonObject, "", "  ")
	if err != nil {
		return fmt.Errorf("Failed to marshal JSON object: %v", err)
	}
	fmt.Println(string(str))
	return nil
}

// CheckMoniker verifies if the moniker matches the MonikerRegexp.
func CheckMoniker(moniker string) bool {
	return MonikerRegexp.MatchString(moniker)
}

// DefaultLogger returns a default logger instance
func DefaultLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	return logger
}

//CheckIP tests whether an IP address is on a subnet
func CheckIP(ip string, portOnlyOk bool) bool {
	if len(ip) == 0 {
		return true
	}
	if ip[0] == ':' { // Port only address
		return !portOnlyOk
	}

	parts := strings.Split(ip, ":")
	trimmedIP := parts[0]

	private := false
	IP := net.ParseIP(trimmedIP)
	if IP == nil {
		return true
	} else {
		_, private24BitBlock, _ := net.ParseCIDR("10.0.0.0/8")
		_, private24BitBlock2, _ := net.ParseCIDR("127.0.0.0/8")
		_, private20BitBlock, _ := net.ParseCIDR("172.16.0.0/12")
		_, private16BitBlock, _ := net.ParseCIDR("192.168.0.0/16")
		private = private24BitBlock2.Contains(IP) || private24BitBlock.Contains(IP) || private20BitBlock.Contains(IP) || private16BitBlock.Contains(IP)
	}
	return private

}

// LogLevel converts a string to a logrus logging level constant
func LogLevel(l string) logrus.Level {
	switch l {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.DebugLevel
	}
}
