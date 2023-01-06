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
	DefaultAppName       = "app"               // 默认应用名称
	DefaultAppMode       = "dev"               // 默认应用环境
	DefaultRunMode       = "debug"             // 默认运行环境
	DefaultAppVersion    = "1.0.0"             // 默认版本
	DefaultLogType       = "file"              // 默认日志类型
	DefaultLogPath       = "/home/logs/app/"   // 默认文件目录
	DefaultChildPath     = "/file-%s.log"      // 默认子目录
	DefaultRotationSize  = 32 * 1024 * 1024    // 默认大小为32M
	DefaultRotationCount = 0                   // 默认不限制
	DefaultRotationTime  = 24 * time.Hour      // 默认每天轮转一次
	DefaultNoBuffWrite   = false               // 不不开启无缓冲写入
	DefaultMaxAge        = 90 * 24 * time.Hour // 默认保留90天
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

type AppConf struct {
	AppName       string        `json:"app_name"`       // 默认应用名称
	AppMode       string        `json:"app_mode"`       // 默认应用环境
	RunMode       string        `json:"run_mode"`       // 默认运行环境
	AppVersion    string        `json:"app_version"`    // 默认版本
	LogType       string        `json:"log_type"`       // 默认日志类型
	LogPath       string        `json:"log_path"`       // 默认文件目录
	ChildPath     string        `json:"child_path"`     // 默认子目录
	RotationSize  int           `json:"rotation_size"`  // 默认大小为32M
	RotationCount int           `json:"rotation_count"` // 默认不限制
	RotationTime  time.Duration `json:"rotation_time"`  // 默认每天轮转一次
	NoBuffWrite   bool          `json:"no_buff_write"`  // 不不开启无缓冲写入
	MaxAge        time.Duration `json:"max_age"`        // 默认保留90天
}

// SetAppConf 更新配置文件
func SetAppConf(conf AppConf) {
	if conf.AppName != "" {
		DefaultAppName = conf.AppName
	}
	if conf.AppMode != "" {
		DefaultAppMode = conf.AppMode
	}
	if conf.RunMode != "" {
		DefaultRunMode = conf.RunMode
	}
	if conf.AppVersion != "" {
		DefaultAppVersion = conf.AppVersion
	}
	if conf.LogType != "" {
		DefaultLogType = conf.LogType
	}
	if conf.LogPath != "" {
		DefaultLogPath = conf.LogPath
	}
	if conf.ChildPath != "" {
		DefaultChildPath = conf.ChildPath
	}
	if conf.RotationSize > 0 {
		DefaultRotationSize = conf.RotationSize
	}

	if conf.RotationCount > 0 {
		DefaultRotationCount = conf.RotationCount
	}
	if conf.RotationTime > 0 {
		DefaultRotationTime = conf.RotationTime
	}
	if conf.NoBuffWrite != false {
		DefaultNoBuffWrite = true
	}
	if conf.MaxAge > 0 {
		DefaultMaxAge = conf.MaxAge
	}
}
