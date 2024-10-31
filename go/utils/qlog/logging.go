package qlog

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mhaqqiw/sdk/go/qconstant"
	"github.com/mhaqqiw/sdk/go/qentity"
	"github.com/mhaqqiw/sdk/go/utils/qmodule"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

var Config qentity.Monitoring
var DisableTrace bool

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
		return strings.TrimPrefix(relativePath, "/")
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

func LogPrint(typeLog string, identifier string, trace string, err string) {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05 -0700")

	if typeLog == "" {
		typeLog = qconstant.ERROR
	}
	if DisableTrace {
		trace = ""
	}
	log.Printf("[%s][%s][%s] - %s -> [%s] %s\n", formattedTime, typeLog, identifier, trace, typeLog, strings.TrimSpace(err))
	if Config.NRConfig.IsEnabled {
		//TODO: send NR metrics
	}
}

func InitNRConfig(name string, key string, isForward bool) (*newrelic.Application, error) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(name),
		newrelic.ConfigLicense(key),
		newrelic.ConfigEnabled(true),
		newrelic.ConfigAppLogForwardingEnabled(isForward),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func Middleware(router *gin.Engine, app *newrelic.Application) {
	router.Use(nrgin.Middleware(app))
	router.Use(func(ctx *gin.Context) {
		traceId := ctx.Query("trace_id")
		txn, ok := ctx.Value("newRelicTransaction").(*newrelic.Transaction)
		if ok && traceId != "" {
			txn.AddAttribute("trace_id", traceId)
		} else if ok && traceId == "" {
			uuid, _ := qmodule.GenerateUUIDV1()
			txn.AddAttribute("trace_id", uuid)
		}
		ctx.Next()
	})
}
