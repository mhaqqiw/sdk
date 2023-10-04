package http

import (
	"fmt"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/mhaqqiw/sdk/go/constant"
	"github.com/mhaqqiw/sdk/go/entity"
	logging "github.com/mhaqqiw/sdk/go/utils/log"
	"github.com/mhaqqiw/sdk/go/utils/module"
	"github.com/mitchellh/mapstructure"
)

func Return(c *gin.Context, statusCode int, metadata entity.Metadata, message interface{}, monitoring entity.Monitoring) {
	start, isExist := c.Get("start")
	if !isExist {
		logging.LogPrint(constant.INFO, "c.Get", logging.Trace(), "request didn't set timestamp", monitoring)
	}
	a := time.Time{}
	mapstructure.Decode(start, &a)

	if statusCode >= 1 && statusCode < 300 {
		var res = entity.Response{Status: "ok", Code: statusCode, Message: message, ProcessTime: module.CountElapsed(a), Version: metadata.Version}
		c.IndentedJSON(int(res.Code), res)
	} else if statusCode >= 300 && statusCode < 400 {
		c.Redirect(statusCode, fmt.Sprintf("%v", message))
	} else if statusCode >= 400 && statusCode < 600 {
		var res = entity.Response{Status: "error", Code: statusCode, Message: message, ProcessTime: module.CountElapsed(a), Version: metadata.Version}
		c.AbortWithStatusJSON(statusCode, res)
	} else {
		var res = entity.Response{Status: "error", Code: 500, Message: fmt.Sprintf("http status code: %v is not listed", statusCode), ProcessTime: module.CountElapsed(a), Version: metadata.Version}
		c.AbortWithStatusJSON(statusCode, res)
	}
}
