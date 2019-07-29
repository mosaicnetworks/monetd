package common

import (
	"encoding/json"
	"fmt"
	"regexp"

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
