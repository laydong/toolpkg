package glogs

import (
	cloudutlis "cloud-utlis"
	"cloud-utlis/utils"
	"fmt"
	"github.com/gin-gonic/gin"
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

	logPath := fmt.Sprintf("%s/%s/%s", lc.logPath, lc.appName, lc.childPath)

	var core zapcore.Core
	//打印至文件中
	if lc.logType == "file" {
		configs := zap.NewProductionEncoderConfig()
		configs.FunctionKey = "func"
		configs.EncodeTime = timeEncoder

		w := zapcore.AddSync(GetWriter(logPath, lc))

		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(configs),
			w,
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

func Info(template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	writers(Sugar, LevelInfo, msg, fields...)
}

func InfoF(c *gin.Context, template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	fields = append(fields, zap.Any("datetime", time.Now().Format(cloudutlis.TimeFormat)))
	fields = append(fields, zap.String(cloudutlis.RequestIdKey, utils.GetRequestIdKey(c)))
	fields = append(fields, zap.String(MessageType, "api_log"))
	writers(Sugar, LevelInfo, msg, fields...)
}

func InfoSdk(c *gin.Context, template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	fields = append(fields, zap.Any("datetime", time.Now().Format(cloudutlis.TimeFormat)))
	fields = append(fields, zap.String(cloudutlis.RequestIdKey, utils.GetRequestIdKey(c)))
	fields = append(fields, zap.String(MessageType, "sdk_log"))
	writers(Sugar, LevelInfo, msg, fields...)
}

func InfoDB(c *gin.Context, template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	fields = append(fields, zap.Any("datetime", time.Now().Format(cloudutlis.TimeFormat)))
	fields = append(fields, zap.String(cloudutlis.RequestIdKey, utils.GetRequestIdKey(c)))
	fields = append(fields, zap.String(MessageType, "db_log"))
	writers(Sugar, LevelInfo, msg, fields...)
}

func Warn(template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	writer(nil, Sugar, LevelWarn, msg, LevelWarn, fields...)
}
func WarnF(c *gin.Context, title string, template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	fields = append(fields, zap.Any("datetime", time.Now().Format(cloudutlis.TimeFormat)))
	fields = append(fields, zap.String(cloudutlis.RequestIdKey, utils.GetRequestIdKey(c)))
	fields = append(fields, zap.String(MessageType, "api_log"))
	writer(nil, Sugar, LevelWarn, msg, title, fields...)
}

func Error(template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	writer(nil, Sugar, LevelError, msg, LevelError, fields...)
}

func ErrorF(c *gin.Context, template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	fields = append(fields, zap.Any("datetime", time.Now().Format(cloudutlis.TimeFormat)))
	fields = append(fields, zap.String(cloudutlis.RequestIdKey, utils.GetRequestIdKey(c)))
	fields = append(fields, zap.String(MessageType, "api_log"))
	writers(Sugar, LevelError, msg, fields...)
}
