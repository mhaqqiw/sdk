package qlog

import (
	"context"
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

var (
	Config       qentity.Monitoring
	DisableTrace bool
)

const (
	StackCallerDefault  = 2 // default caller that call Tracer outside this package
	StackCallerExternal = 3 // external caller (outside package) eg: qlog.ErrorCtx, qlog.InfoCtx, qlog.DebugCtx
	TRACK_ID            = "track-id"
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
		return strings.TrimPrefix(relativePath, "/")
	}

	// If the absolute path does not start with the current working directory, return the original path
	return absolutePath
}

func Trace(stackCallers ...int) string {
	stackLevel := StackCallerDefault

	if len(stackCallers) > 0 {
		stackLevel = stackCallers[0]
	}

	pc := make([]uintptr, 15)
	runtime.Callers(stackLevel, pc)
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

func ErrorCtx(ctx context.Context, err error) {
	trace := Trace(StackCallerExternal)

	// get track_id from context
	trackID, _ := ctx.Value(TRACK_ID).(string)
	if trackID == "" {
		trackID, _ = qmodule.GenerateUUIDV1()
	}

	LogPrint(qconstant.ERROR, trackID, trace, err.Error())
}

func InfoCtx(ctx context.Context, message string) {
	trace := Trace(StackCallerExternal)

	// get track_id from context
	trackID, _ := ctx.Value(TRACK_ID).(string)
	if trackID == "" {
		trackID, _ = qmodule.GenerateUUIDV1()
	}

	LogPrint(qconstant.INFO, trackID, trace, message)
}

func DebugCtx(ctx context.Context, message string) {
	trace := Trace(StackCallerExternal)

	// get track_id from context
	trackID, _ := ctx.Value(TRACK_ID).(string)
	if trackID == "" {
		trackID, _ = qmodule.GenerateUUIDV1()
	}

	LogPrint(qconstant.DEBUG, trackID, trace, message)
}

func Error(err error) {
	trace := Trace(StackCallerExternal)
	uuid, _ := qmodule.GenerateUUIDV1()

	LogPrint(qconstant.ERROR, uuid, trace, err.Error())
}

func Info(message string) {
	trace := Trace(StackCallerExternal)
	uuid, _ := qmodule.GenerateUUIDV1()

	LogPrint(qconstant.INFO, uuid, trace, message)
}

func Debug(message string) {
	trace := Trace(StackCallerExternal)
	uuid, _ := qmodule.GenerateUUIDV1()

	LogPrint(qconstant.DEBUG, uuid, trace, message)
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
