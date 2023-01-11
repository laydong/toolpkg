package toolpkg

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/laydong/toolpkg/logx"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// MiddlewareApiLog 记录框架出入参, 开启链路追踪
func MiddlewareApiLog(c *gin.Context) {
	start := time.Now()
	traceId := c.GetHeader(XtraceKey)
	if traceId == "" {
		traceId = uuid.New().String()
		c.Header(XtraceKey, traceId)
		c.Header(RequestIdKey, traceId)
	}
	c.Set(XtraceKey, traceId)
	c.Set(RequestIdKey, traceId)
	var request Req
	request.Method = c.Request.Method
	request.URL = c.Request.URL.String()
	request.Path = c.Request.URL.Path
	request.Query = c.Request.URL.RawQuery
	request.Headers = c.Request.Header
	request.IP = c.ClientIP()
	var body []byte
	body, _ = c.GetRawData()
	// 将原body塞回去
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	c.Next()
	request.Body = string(body)
	// 自行处理日志
	logx.InfoApi(c, "API请求日志",
		zap.Any("request", request),
		zap.Any("Resp_body", blw.body.String()),
		zap.Any("status_code", c.Writer.Status()),
		zap.Any("run_time", fmt.Sprintf("%.3fms", float64(time.Since(start).Nanoseconds())/1e6)),
		zap.Any("error", strings.TrimRight(c.Errors.ByType(gin.ErrorTypePrivate).String(), "\n")),
		zap.Any("source", c.GetHeader("app_name")),
	)
}

type Req struct {
	URL     string              `json:"url"`
	Method  string              `json:"method"`
	IP      string              `json:"ip"`
	Path    string              `json:"path"`
	Headers map[string][]string `json:"headers"`
	Query   interface{}         `json:"query"`
	Body    interface{}         `json:"body"`
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

type AppConf struct {
	AppName       string        `json:"app_name"`       // 默认应用名称
	AppMode       string        `json:"app_mode"`       // 默认应用环境
	LogType       string        `json:"log_type"`       // 默认日志类型
	LogPath       string        `json:"log_path"`       // 默认文件目录
	ChildPath     string        `json:"child_path"`     // 默认子目录
	RotationSize  int           `json:"rotation_size"`  // 默认大小为32M
	RotationCount int           `json:"rotation_count"` // 默认不限制
	RotationTime  time.Duration `json:"rotation_time"`  // 默认每天轮转一次
	NoBuffWrite   bool          `json:"no_buff_write"`  // 不不开启无缓冲写入
	MaxAge        time.Duration `json:"max_age"`        // 默认保留90天
}

// InitLog 初始化日志服务
func InitLog(conf AppConf) {

	defaultConfig := logx.Config{
		AppName:       logx.DefaultAppName,
		AppMode:       logx.DefaultAppMode,
		LogType:       logx.DefaultLogType,
		LogPath:       logx.DefaultLogPath,
		ChildPath:     logx.DefaultChildPath,
		RotationSize:  int64(logx.DefaultRotationSize),
		RotationCount: uint(logx.DefaultRotationCount),
		NoBuffWrite:   logx.DefaultNoBuffWrite,
		RotationTime:  logx.DefaultRotationTime,
		MaxAge:        logx.DefaultMaxAge,
	}
	if conf.AppName != "" {
		defaultConfig.AppName = conf.AppName
	}
	if conf.AppMode != "" {
		defaultConfig.AppMode = conf.AppMode
	}
	if conf.LogType != "" {
		defaultConfig.LogType = conf.LogType
	}
	if conf.LogPath != "" {
		defaultConfig.LogPath = conf.LogPath
	}
	if conf.ChildPath != "" {
		defaultConfig.ChildPath = conf.ChildPath
	}
	if conf.RotationSize > 0 {
		defaultConfig.RotationSize = int64(conf.RotationSize)
	}

	if conf.RotationCount > 0 {
		defaultConfig.RotationCount = uint(conf.RotationCount)
	}
	if conf.RotationTime > 0 {
		defaultConfig.RotationTime = conf.RotationTime
	}
	if conf.NoBuffWrite != false {
		defaultConfig.NoBuffWrite = true
	}
	if conf.MaxAge > 0 {
		defaultConfig.MaxAge = conf.MaxAge
	}
	logx.InitLog(&defaultConfig)
}

// GetNewGinContext 获取新的上下文
func GetNewGinContext() *gin.Context {
	ctx := new(gin.Context)
	uid := uuid.New().String()
	ctx.Request = &http.Request{
		Header: make(map[string][]string),
	}
	ctx.Request.Header.Set(XtraceKey, uid)
	ctx.Request.Header.Set(RequestIdKey, uid)
	ctx.Set(RequestIdKey, uid)
	ctx.Set(XtraceKey, uid)
	return ctx
}
