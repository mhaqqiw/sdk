package logging

import (
	"fmt"
	c "github/mhaqqiw/sdk/go/constant"
	"github/mhaqqiw/sdk/go/entity"
	"log"
	"runtime"
	"strings"
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

func LogPrint(typeLog string, identifier string, trace string, err string, monitoring entity.Monitoring) {
	if typeLog == "" {
		typeLog = c.ERROR
	}
	log.Printf("[%s][%s] - %s \n\t [%s] %s\n", typeLog, identifier, trace, typeLog, strings.TrimSpace(err))
	if monitoring.NRConfig.IsEnabled {
		//TODO: send NR metrics
	}
}
