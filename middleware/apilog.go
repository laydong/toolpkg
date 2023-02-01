package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/laydong/toolpkg/logx"
	"go.uber.org/zap"
	"io/ioutil"
	"strings"
	"time"
)

// MiddlewareApiLog 记录框架出入参, 开启链路追踪
func MiddlewareApiLog(c *gin.Context) {
	start := time.Now()
	traceId := c.GetHeader(logx.XtraceKey)
	if traceId == "" {
		traceId = uuid.New().String()
		c.Header(logx.XtraceKey, traceId)
		c.Header(logx.RequestIdKey, traceId)
	}
	c.Set(logx.XtraceKey, traceId)
	c.Set(logx.RequestIdKey, traceId)
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
		zap.Any("run_time_long", time.Since(start)*time.Millisecond), //运行时长毫秒
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
