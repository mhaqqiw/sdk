package qhttp

import (
	"encoding/json"
	"fmt"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/mhaqqiw/sdk/go/qconstant"
	"github.com/mhaqqiw/sdk/go/qentity"
	"github.com/mhaqqiw/sdk/go/utils/qlog"
	"github.com/mhaqqiw/sdk/go/utils/qmodule"
	"github.com/mitchellh/mapstructure"
)

var Metadata qentity.Metadata

func Return(c *gin.Context, statusCode int, message interface{}) {
	start, isExist := c.Get("start")
	if !isExist {
		qlog.LogPrint(qconstant.INFO, "c.Get", qlog.Trace(), "request didn't set timestamp")
	}
	a := time.Time{}
	if start != nil {
		mapstructure.Decode(start, &a)
	}

	var rawMessage json.RawMessage
	var err error
	switch v := message.(type) {
	case string:
		rawMessage = json.RawMessage([]byte(fmt.Sprintf(`"%s"`, v)))
	case []byte:
		rawMessage = json.RawMessage(v)
	default:
		rawMessage, err = json.Marshal(v) // Best-effort conversion to JSON
		if err != nil {
			qlog.LogPrint(qconstant.ERROR, "json.Marshal", qlog.Trace(), err.Error())
			rawMessage = json.RawMessage(`{"error":"failed to marshal message"}`)
		}
	}

	if statusCode >= 1 && statusCode < 300 {
		res := qentity.Response{Status: "ok", Code: statusCode, Message: rawMessage, ProcessTime: qmodule.CountElapsed(a)}
		c.IndentedJSON(res.Code, res)
	} else if statusCode >= 300 && statusCode < 400 {
		c.Redirect(statusCode, fmt.Sprintf("%v", message))
	} else if statusCode >= 400 && statusCode < 600 {
		res := qentity.Response{Status: "error", Code: statusCode, Message: rawMessage, ProcessTime: qmodule.CountElapsed(a)}
		c.AbortWithStatusJSON(statusCode, res)
	} else {
		res := qentity.Response{Status: "error", Code: 500, Message: json.RawMessage(fmt.Sprintf(`{"error":"http status code: %v is not listed"}`, statusCode)), ProcessTime: qmodule.CountElapsed(a)}
		c.AbortWithStatusJSON(500, res)
	}
}
