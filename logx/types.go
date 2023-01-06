package logx

import (
	"github.com/laydong/toolpkg"
	"go.uber.org/zap"
	"time"
)

const (
	//DefaultAppName       = "app"               // 默认应用名称
	//DefaultAppMode       = "dev"               // 默认应用环境
	//DefaultLogType       = "file"              // 默认日志类型
	//DefaultLogPath       = "/home/logs/app/"   // 默认文件目录
	//DefaultChildPath     = "/file-%s.log"      // 默认子目录
	//DefaultRotationSize  = 32 * 1024 * 1024    // 默认大小为32M
	//DefaultRotationCount = 0                   // 默认不限制
	//DefaultRotationTime  = 24 * time.Hour      // 默认每天轮转一次
	//DefaultNoBuffWrite   = false               // 不不开启无缓冲写入
	//DefaultMaxAge        = 90 * 24 * time.Hour // 默认保留90天

	LevelInfo   = "info"
	LevelWarn   = "warn"
	LevelError  = "error"
	MessageType = "message_type" //日志类型
)

var (
	HeaderAppName    = "app-name"
	KeyPath          = "path"
	KeyTitle         = "title"
	KeyOriginAppName = "origin_app_name"
)

type Config struct {
	appName       string        // 应用名
	appMode       string        // 应用环境
	logType       string        // 日志类型
	logPath       string        // 日志主路径
	childPath     string        // 日志子路径+文件名
	RotationSize  int64         // 单个文件大小
	RotationCount uint          // 可以保留的文件个数
	NoBuffWrite   bool          // 设置无缓冲日志写入
	RotationTime  time.Duration // 日志分割的时间
	MaxAge        time.Duration // 日志最大保留的天数
}

type LogOptionFunc func(*Config)

type Field = zap.Field

var (
	Sugar *zap.Logger

	defaultLogLevel = zap.NewAtomicLevel()
	DefaultConfig   = &Config{
		appName:       toolpkg.DefaultAppName,
		appMode:       toolpkg.DefaultAppMode,
		logType:       toolpkg.DefaultLogType,
		logPath:       toolpkg.DefaultLogPath,
		childPath:     toolpkg.DefaultChildPath,
		RotationSize:  int64(toolpkg.DefaultRotationSize),
		RotationCount: uint(toolpkg.DefaultRotationCount),
		NoBuffWrite:   toolpkg.DefaultNoBuffWrite,
		RotationTime:  toolpkg.DefaultRotationTime,
		MaxAge:        toolpkg.DefaultMaxAge,
	}
)
