package qhttp

import (
	"fmt"
	"net/http"

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

	if statusCode >= 1 && statusCode < 300 {
		res := qentity.Response{Status: "ok", Code: statusCode, Message: message, ProcessTime: qmodule.CountElapsed(a), Version: Metadata.Version}
		c.IndentedJSON(res.Code, res)
	} else if statusCode >= 300 && statusCode < 400 {
		c.Redirect(statusCode, fmt.Sprintf("%v", message))
	} else if statusCode >= 400 && statusCode < 600 {
		res := qentity.Response{Status: "error", Code: statusCode, Message: message, ProcessTime: qmodule.CountElapsed(a), Version: Metadata.Version}
		c.AbortWithStatusJSON(statusCode, res)
	} else {
		res := qentity.Response{Status: "error", Code: http.StatusInternalServerError, Message: fmt.Sprintf("http status code: %v is not listed", statusCode), ProcessTime: qmodule.CountElapsed(a), Version: Metadata.Version}
		c.AbortWithStatusJSON(http.StatusInternalServerError, res)
	}
}
