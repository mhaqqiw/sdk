package qlog

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/mhaqqiw/sdk/go/qconstant"
	"github.com/mhaqqiw/sdk/go/qentity"
)

func Trace() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	_, line := f.FileLine(pc[0])
	s := f.Name()
	index := strings.LastIndex(s, ".")
	if index > 0 {
		s = s[:index]
	}
	return fmt.Sprintf("%s:%d - %s", s, line, f.Name())
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
