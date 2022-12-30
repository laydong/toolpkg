package toolpkg

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

const (
	XtraceKey    = "trace-id"            //外部链路ID
	RequestIdKey = "request-id"          //链路ID
	TimeFormat   = "2006-01-02 15:04:05" //默认时间
)

const (
	PROXY_HTTP int = iota
	PROXY_SOCKS4
	PROXY_SOCKS5
	PROXY_SOCKS4A

	// CURL like OPT
	OPT_AUTOREFERER
	OPT_FOLLOWLOCATION
	OPT_CONNECTTIMEOUT
	OPT_CONNECTTIMEOUT_MS
	OPT_MAXREDIRS
	OPT_PROXYTYPE
	OPT_TIMEOUT
	OPT_TIMEOUT_MS
	OPT_COOKIEJAR
	OPT_INTERFACE
	OPT_PROXY
	OPT_REFERER
	OPT_USERAGENT

	// Other OPT
	OPT_REDIRECT_POLICY
	OPT_PROXY_FUNC
	OPT_DEBUG
	OPT_UNSAFE_TLS

	OPT_CONTEXT
)

func GetNewGinContext() *gin.Context {
	ctx := new(gin.Context)
	uid := uuid.NewV4().String()
	ctx.Request = &http.Request{
		Header: make(map[string][]string),
	}
	ctx.Request.Header.Set(XtraceKey, uid)
	ctx.Request.Header.Set(RequestIdKey, uid)
	ctx.Set(RequestIdKey, uid)
	ctx.Set(XtraceKey, uid)
	return ctx
}
