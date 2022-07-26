package log

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/laydong/toolpkg"
	"github.com/laydong/toolpkg/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"time"
)

// InitLog 初始化日志文件 logPath= /home/logs/app/appName/childPath
func InitLog(options ...LogOptionFunc) {
	for _, f := range options {
		f(DefaultConfig)
	}

	Sugar = initSugar(DefaultConfig)
}

func initSugar(lc *Config) *zap.Logger {
	loglevel := zapcore.InfoLevel
	defaultLogLevel.SetLevel(loglevel)

	logPath := fmt.Sprintf("%s/%s/%s", lc.logPath, lc.appName, fmt.Sprintf(lc.childPath, time.Now().Format("2006-01-02")))

	var core zapcore.Core
	//打印至文件中
	if lc.logType == "file" {
		configs := zap.NewProductionEncoderConfig()
		configs.FunctionKey = "func"
		configs.EncodeTime = timeEncoder

		//w := zapcore.AddSync(GetWriter(logPath, lc))

		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(configs),
			GetWriteSyncer(logPath, lc),
			defaultLogLevel,
		)
		log.Printf("[glogs_sugar] log success")
	} else {
		// 打印在控制台
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), defaultLogLevel)
		log.Printf("[glogs_sugar] log success")
	}

	filed := zap.Fields(zap.String("app_name", lc.appName), zap.String("app_mode", lc.appMode))
	//return zap.New(core, filed, zap.AddCaller(), zap.AddCallerSkip(3))
	return zap.New(core, filed)
	//Sugar = logger.Sugar()

}

func InfoApi(c *gin.Context, template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	fields = append(fields, zap.Any("datetime", time.Now().Format(toolpkg.TimeFormat)))
	fields = append(fields, zap.String(toolpkg.RequestIdKey, utils.GetRequestIdKey(c)))
	fields = append(fields, zap.String(MessageType, "api_log"))
	writers(Sugar, LevelInfo, msg, fields...)
}

func InfoSdk(c *gin.Context, template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	fields = append(fields, zap.Any("datetime", time.Now().Format(toolpkg.TimeFormat)))
	fields = append(fields, zap.String(toolpkg.RequestIdKey, utils.GetRequestIdKey(c)))
	fields = append(fields, zap.String(MessageType, "sdk_log"))
	writers(Sugar, LevelInfo, msg, fields...)
}

func InfoDB(c *gin.Context, template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	fields = append(fields, zap.Any("datetime", time.Now().Format(toolpkg.TimeFormat)))
	fields = append(fields, zap.String(toolpkg.RequestIdKey, utils.GetRequestIdKey(c)))
	fields = append(fields, zap.String(MessageType, "db_log"))
	writers(Sugar, LevelInfo, msg, fields...)
}

func InfoF(c *gin.Context, template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	fields = append(fields, zap.Any("datetime", time.Now().Format(toolpkg.TimeFormat)))
	fields = append(fields, zap.String(toolpkg.RequestIdKey, utils.GetRequestIdKey(c)))
	fields = append(fields, zap.String(MessageType, "info_log"))
	writers(Sugar, LevelInfo, msg, fields...)
}

func WarnF(c *gin.Context, title string, template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	fields = append(fields, zap.Any("datetime", time.Now().Format(toolpkg.TimeFormat)))
	fields = append(fields, zap.String(toolpkg.RequestIdKey, utils.GetRequestIdKey(c)))
	fields = append(fields, zap.String(MessageType, "warn_log"))
	writer(nil, Sugar, LevelWarn, msg, title, fields...)
}

func ErrorF(c *gin.Context, template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	fields = append(fields, zap.Any("datetime", time.Now().Format(toolpkg.TimeFormat)))
	fields = append(fields, zap.String(toolpkg.RequestIdKey, utils.GetRequestIdKey(c)))
	fields = append(fields, zap.String(MessageType, "err_log"))
	writers(Sugar, LevelError, msg, fields...)
}

func Warn(template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	writer(nil, Sugar, LevelWarn, msg, LevelWarn, fields...)
}

func Error(template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	writer(nil, Sugar, LevelError, msg, LevelError, fields...)
}

func Info(template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	writers(Sugar, LevelInfo, msg, fields...)
}

func InfoLogId(logId, template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	fields = append(fields, zap.Any("datetime", time.Now().Format(toolpkg.TimeFormat)))
	fields = append(fields, zap.String(toolpkg.RequestIdKey, logId))
	fields = append(fields, zap.String(MessageType, "info_log"))
	writers(Sugar, LevelError, msg, fields...)
}

func WarnLogId(logId, template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	fields = append(fields, zap.Any("datetime", time.Now().Format(toolpkg.TimeFormat)))
	fields = append(fields, zap.String(toolpkg.RequestIdKey, logId))
	fields = append(fields, zap.String(MessageType, "warn_log"))
	writers(Sugar, LevelError, msg, fields...)
}

func ErrorLogId(logId, template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	fields = append(fields, zap.Any("datetime", time.Now().Format(toolpkg.TimeFormat)))
	fields = append(fields, zap.String(toolpkg.RequestIdKey, logId))
	fields = append(fields, zap.String(MessageType, "err_log"))
	writers(Sugar, LevelError, msg, fields...)
}
