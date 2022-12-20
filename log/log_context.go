package log

import "github.com/laydong/toolpkg/utils"

// LoggerContext 日志
type LoggerContext interface {
	InfoF(template string, args ...interface{})
	WarnF(template string, args ...interface{})
	ErrorF(template string, args ...interface{})
	Field(key string, value interface{}) Field
}

func (ctx *LogContext) InfoF(template string, args ...interface{}) {
	InfoLogId(ctx.logId, template, args...)
}

func (ctx *LogContext) WarnF(template string, args ...interface{}) {
	WarnLogId(ctx.logId, template, args...)
}

func (ctx *LogContext) ErrorF(template string, args ...interface{}) {
	ErrorLogId(ctx.logId, template, args...)
}

func (ctx *LogContext) Field(key string, value interface{}) Field {
	return String(key, value)
}

// LogContext logger
type LogContext struct {
	logId    string
	clientIP string
}

var _ LoggerContext = &LogContext{}

// NewLogContext new obj
func NewLogContext(logId string) *LogContext {
	ctx := &LogContext{
		logId:    logId,
		clientIP: utils.GetClientIp(),
	}
	return ctx
}

// GetLogId 得到LogId
func (ctx *LogContext) GetLogId() string {
	return ctx.logId
}

// GetClientIP 得到clientIP
func (ctx *LogContext) GetClientIP() string {
	return ctx.clientIP
}
