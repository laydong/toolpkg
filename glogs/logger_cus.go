package glogs

import (
	"go.uber.org/zap"
	"net/http"
)

type CusLog struct {
	Logger *zap.Logger
	Config *Config
}

// NewLogger 得到一个zap.Logger
func NewLogger(options ...LogOptionFunc) *CusLog {
	var cus = &CusLog{Config: DefaultConfig}
	for _, f := range options {
		f(cus.Config)
	}

	cus.Logger = initSugar(cus.Config)
	return cus
}

func (l *CusLog) Info(template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	writer(nil, l.Logger, LevelInfo, msg, LevelInfo, fields...)
}
func (l *CusLog) InfoF(r *http.Request, title string, template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	writer(r, l.Logger, LevelInfo, msg, title, fields...)
}

func (l *CusLog) Warn(template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	writer(nil, l.Logger, LevelWarn, msg, LevelWarn, fields...)
}
func (l *CusLog) WarnF(r *http.Request, title string, template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	writer(r, l.Logger, LevelWarn, msg, title, fields...)
}

func (l *CusLog) Error(template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	writer(nil, l.Logger, LevelError, msg, LevelError, fields...)
}
func (l *CusLog) ErrorF(r *http.Request, title string, template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	writer(r, l.Logger, LevelError, msg, title, fields...)
}
