package qlog

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/mhaqqiw/sdk/go/qconstant"
	"github.com/mhaqqiw/sdk/go/qentity"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func getRelativePath(absolutePath string) string {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		// Handle the error, e.g., log it or use a default path
		return absolutePath
	}

	// Check if the absolute path starts with the current working directory
	if strings.HasPrefix(absolutePath, currentDir) {
		// Get the relative path by removing the current working directory from the beginning
		relativePath := absolutePath[len(currentDir):]

		// Remove the leading slash if present
		if strings.HasPrefix(relativePath, "/") {
			relativePath = relativePath[1:]
		}

		return relativePath
	}

	// If the absolute path does not start with the current working directory, return the original path
	return absolutePath
}

func Trace() string {
	pc := make([]uintptr, 15)
	runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc)

	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		if !strings.Contains(frame.File, "runtime/") {
			relativePath := getRelativePath(frame.File)
			return fmt.Sprintf("%s:%d - %s", relativePath, frame.Line, frame.Function)
		}
	}

	return "unable to get trace"
}

func LogPrint(typeLog string, identifier string, trace string, err string, monitoring qentity.Monitoring) {
	if typeLog == "" {
		typeLog = qconstant.ERROR
	}
	log.Printf("[%s][%s] - %s \n\t [%s] %s\n", typeLog, identifier, trace, typeLog, strings.TrimSpace(err))
	if monitoring.NRConfig.IsEnabled {
		//TODO: send NR metrics
	}
}

func InitNRConfig(name string, key string, isForward bool) (*newrelic.Application, error) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(name),
		newrelic.ConfigLicense(key),
		newrelic.ConfigAppLogForwardingEnabled(isForward),
	)
	if err != nil {
		return nil, err
	}
	return app, nil
}
