package log

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"time"
)

// 设置应用名称,默认值default-app
func SetLogAppName(appName string) LogOptionFunc {
	return func(c *Config) {
		if appName != "" {
			c.appName = appName
		}
	}
}

// 设置环境变量,标识当前应用运行的环境,默认值dev
func SetLogAppMode(appMode string) LogOptionFunc {
	return func(c *Config) {
		if appMode != "" {
			c.appMode = appMode
		}
	}
}

// 设置日志类型,日志类型目前分为2种,console和file,默认值file
func SetLogType(logType string) LogOptionFunc {
	return func(c *Config) {
		if logType != "" {
			c.logType = logType
		}
	}
}

// 设置日志目录,这个是主目录,程序会给此目录拼接上项目名,子目录以及文件,默认值/home/logs/app
func SetLogPath(logPath string) LogOptionFunc {
	return func(c *Config) {
		if logPath != "" {
			c.logPath = logPath
		}
	}
}

// 设置子目录—+文件名,保证一个类型的文件在同一个文件夹下面便于区分、默认值glogs/%Y-%m-%d.log
func SetLogChildPath(childPath string) LogOptionFunc {
	return func(c *Config) {
		if childPath != "" {
			c.childPath = childPath
		}
	}
}

// 设置单个文件最大值byte,默认值32M
func SetLogMaxSize(size int64) LogOptionFunc {
	return func(c *Config) {
		if size > 0 {
			c.RotationSize = size
		}
	}
}

// SetLogMaxAge 设置文件最大保留时间、默认值7天
func SetLogMaxAge(maxAge time.Duration) LogOptionFunc {
	return func(c *Config) {
		if maxAge != 0 {
			c.MaxAge = maxAge
		}
	}
}

// SetRotationTime 设置文件分割时间、默认值24*time.Hour(按天分割)
func SetRotationTime(rotationTime time.Duration) LogOptionFunc {
	return func(c *Config) {
		if rotationTime != 0 {
			c.RotationTime = rotationTime
		}
	}
}

// SetRotationCount 设置保留的最大文件数量、没有默认值(表示不限制)
func SetRotationCount(n uint) LogOptionFunc {
	return func(c *Config) {
		if n != 0 {
			c.RotationCount = n
		}
	}
}

// SetNoBuffWriter 设置无缓冲写入日志，比较消耗性能
func SetNoBuffWriter() LogOptionFunc {
	return func(c *Config) {
		c.NoBuffWrite = true
	}
}

func writer(r *http.Request, logger *zap.Logger, level, msg string, title string, fields ...zap.Field) {
	if logger == nil {
		fmt.Println(msg)
		return
	}

	if r == nil {
		fields = append(fields, zap.String(KeyTitle, title))
		do(logger, level, msg, fields...)
		return
	}

	//requestID := r.Header.Get(RequestIDName)
	originAppName := r.Header.Get(HeaderAppName)
	path := r.RequestURI
	fields = append(fields, zap.String(KeyPath, path),
		//zap.String(RequestIDName, requestID),
		zap.String(KeyTitle, title),
		zap.String(KeyOriginAppName, originAppName))
	do(logger, level, msg, fields...)
	return
}

func writers(logger *zap.Logger, level, msg string, fields ...zap.Field) {
	switch level {
	case LevelInfo:
		logger.Info(msg, fields...)
	case LevelWarn:
		logger.Warn(msg, fields...)
	case LevelError:
		logger.Error(msg, fields...)
	}
}

func do(logger *zap.Logger, level, msg string, fields ...zap.Field) {
	switch level {
	case LevelInfo:
		logger.Info(msg, fields...)
	case LevelWarn:
		logger.Warn(msg, fields...)
	case LevelError:
		logger.Error(msg, fields...)
	}
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	var layout = "2006-01-02 15:04:05"
	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}

	if enc, ok := enc.(appendTimeEncoder); ok {
		enc.AppendTimeLayout(t, layout)
		return
	}

	enc.AppendString(t.Format(layout))
}

// GetWriteSyncer 按天切割按大小切割
// filename 文件名
// RotationSize 每个文件的大小
// MaxAge 文件最大保留天数
// RotationCount 最大保留文件个数
// RotationTime 设置文件分割时间
// RotationCount 设置保留的最大文件数量
func GetWriteSyncer(file string, lc *Config) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,                  // 日志文件的位置
		MaxSize:    int(lc.RotationSize),  // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: int(lc.RotationCount), // 保留旧文件的最大个数
		MaxAge:     int(lc.MaxAge),        // 保留旧文件的最大天数
		Compress:   false,                 // 是否压缩/归档旧文件
	}
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))

	//return zapcore.AddSync(lumberJackLogger)
}

// GetWriter 按天切割按大小切割
// filename 文件名
// RotationSize 每个文件的大小
// MaxAge 文件最大保留天数
// RotationCount 最大保留文件个数
// RotationTime 设置文件分割时间
// RotationCount 设置保留的最大文件数量
//func GetWriter(filename string, lc *Config) io.Writer {
//	// 生成rotatelogs的Logger 实际生成的文件名 stream-2021-5-20.log
//	// demo.log是指向最新日志的连接
//	// 保存7天内的日志，每1小时(整点)分割一第二天志
//	var options []rl.Option
//	if lc.NoBuffWrite {
//		options = append(options, rl.WithNoBuffer())
//	}
//	options = append(options,
//		rl.WithRotationSize(lc.RotationSize),
//		rl.WithRotationCount(lc.RotationCount),
//		rl.WithRotationTime(lc.RotationTime),
//		rl.WithMaxAge(lc.MaxAge),
//		rl.ForceNewFile())
//
//	hook, err := rl.New(
//		filename,
//		options...,
//	)
//
//	if err != nil {
//		panic(err)
//	}
//	return hook
//}

func dealWithArgs(tmp string, args ...interface{}) (msg string, f []zap.Field) {
	var tmpArgs []interface{}
	for _, item := range args {
		if zapField, ok := item.(zap.Field); ok {
			f = append(f, zapField)
		} else {
			tmpArgs = append(tmpArgs, item)
		}
	}
	msg = fmt.Sprintf(tmp, tmpArgs...)
	return
}

func String(key string, value interface{}) zap.Field {
	v := fmt.Sprintf("%v", value)
	return zap.String(key, v)
}
