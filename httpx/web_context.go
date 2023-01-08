package httpx

import (
	"github.com/gin-gonic/gin"
	"github.com/laydong/toolpkg/alarmx"
	"github.com/laydong/toolpkg/logx"
	"github.com/laydong/toolpkg/tracex"
	"github.com/laydong/toolpkg/utils"
	uuid "github.com/satori/go.uuid"
)

// WebHandlerFunc http请求的处理者
type WebHandlerFunc func(*WebContext)

// WebContext http 的context
// WebContext 继承了 gin.Context, 并且扩展了日志功能
type WebContext struct {
	*gin.Context
	*logx.LogContext
	*tracex.TraceContext
	*alarmx.AlarmContext
}

const ginFlag = "__gin__gin"

// NewWebContext 创建 http contextx
func NewWebContext(ginContext *gin.Context) *WebContext {
	obj, existed := ginContext.Get(ginFlag)
	if existed {
		return obj.(*WebContext)
	}

	logId := ginContext.GetHeader(utils.RequestIdKey)
	if logId == "" {
		logId = utils.Md5(uuid.NewV4().String())
		ginContext.Request.Header.Set(utils.RequestIdKey, logId)
		ginContext.Set(utils.RequestIdKey, logId)
	}

	tmp := &WebContext{
		Context:      ginContext,
		LogContext:   logx.NewLogContext(logId),
		TraceContext: tracex.NewTraceContext(ginContext.Request.RequestURI, ginContext.Request.Header),
	}
	ginContext.Set(ginFlag, tmp)

	return tmp
}
