package toolpkg

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"
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

var (
	DefaultAppName               = "app"               // 默认应用名称
	DefaultAppMode               = "dev"               // 默认应用环境
	DefaultRunMode               = "debug"             // 默认运行环境
	DefaultAppVersion            = "1.0.0"             // 默认版本
	DefaultLogType               = "file"              // 默认日志类型
	DefaultLogPath               = "/home/logs/app/"   // 默认文件目录
	DefaultChildPath             = "/file-%s.log"      // 默认子目录
	DefaultRotationSize          = 32 * 1024 * 1024    // 默认大小为32M
	DefaultRotationCount         = 0                   // 默认不限制
	DefaultRotationTime          = 24 * time.Hour      // 默认每天轮转一次
	DefaultNoBuffWrite           = false               // 不不开启无缓冲写入
	DefaultMaxAge                = 90 * 24 * time.Hour // 默认保留90天
	DefaultTraceType             = "zipkin"            //记录类型 zipkin 和 jaeger
	DefaultTraceAddr             = ""
	DefaultTraceMod      float64 = 0
)

// GetNewGinContext 获取新的上下文
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
