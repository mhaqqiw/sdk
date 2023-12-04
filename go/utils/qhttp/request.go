package qhttp

import (
	"fmt"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/mhaqqiw/sdk/go/qconstant"
	"github.com/mhaqqiw/sdk/go/qentity"
	"github.com/mhaqqiw/sdk/go/utils/qlog"
	"github.com/mhaqqiw/sdk/go/utils/qmodule"
	"github.com/mitchellh/mapstructure"
)

var metadata qentity.Metadata

func Return(c *gin.Context, statusCode int, message interface{}) {
	start, isExist := c.Get("start")
	if !isExist {
		qlog.LogPrint(qconstant.INFO, "c.Get", qlog.Trace(), "request didn't set timestamp")
	}
	a := time.Time{}
	mapstructure.Decode(start, &a)

	if statusCode >= 1 && statusCode < 300 {
		var res = qentity.Response{Status: "ok", Code: statusCode, Message: message, ProcessTime: qmodule.CountElapsed(a), Version: metadata.Version}
		c.IndentedJSON(int(res.Code), res)
	} else if statusCode >= 300 && statusCode < 400 {
		c.Redirect(statusCode, fmt.Sprintf("%v", message))
	} else if statusCode >= 400 && statusCode < 600 {
		var res = qentity.Response{Status: "error", Code: statusCode, Message: message, ProcessTime: qmodule.CountElapsed(a), Version: metadata.Version}
		c.AbortWithStatusJSON(statusCode, res)
	} else {
		var res = qentity.Response{Status: "error", Code: 500, Message: fmt.Sprintf("http status code: %v is not listed", statusCode), ProcessTime: qmodule.CountElapsed(a), Version: metadata.Version}
		c.AbortWithStatusJSON(statusCode, res)
	}
}
