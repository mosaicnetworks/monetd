package common

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/sirupsen/logrus"
)

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

//GetNodeSafeLabel converts a free format string into a node label friendly format
//Anything other than an alphanumeric is converted to _
func GetNodeSafeLabel(moniker string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	return reg.ReplaceAllString(moniker, "_")
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
