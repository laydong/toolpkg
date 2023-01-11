package logx

import (
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

	LevelInfo    = "info"
	LevelWarn    = "warn"
	LevelError   = "error"
	MessageType  = "message_type"        //日志类型
	XtraceKey    = "trace-id"            //外部链路ID
	RequestIdKey = "request-id"          //链路ID
	TimeFormat   = "2006-01-02 15:04:05" //默认时间
)

var (
	HeaderAppName    = "app-name"
	KeyPath          = "path"
	KeyTitle         = "title"
	KeyOriginAppName = "origin_app_name"

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

var ()

type Config struct {
	AppName       string        `json:"app_name"`       // 应用名
	AppMode       string        `json:"app_mode"`       // 应用环境
	LogType       string        `json:"log_type"`       // 日志类型
	LogPath       string        `json:"log_path"`       // 日志主路径
	ChildPath     string        `json:"child_path"`     // 日志子路径+文件名
	RotationSize  int64         `json:"rotation_size"`  // 单个文件大小
	RotationCount uint          `json:"rotation_count"` // 可以保留的文件个数
	NoBuffWrite   bool          `json:"no_buff_write"`  // 设置无缓冲日志写入
	RotationTime  time.Duration `json:"rotation_time"`  // 日志分割的时间
	MaxAge        time.Duration `json:"max_age"`        // 日志最大保留的天数
}

type LogOptionFunc func(*Config)

type Field = zap.Field

var (
	Sugar *zap.Logger

	defaultLogLevel = zap.NewAtomicLevel()
	DefaultConfig   = &Config{
		AppName:       DefaultAppName,
		AppMode:       DefaultAppMode,
		LogType:       DefaultLogType,
		LogPath:       DefaultLogPath,
		ChildPath:     DefaultChildPath,
		RotationSize:  int64(DefaultRotationSize),
		RotationCount: uint(DefaultRotationCount),
		NoBuffWrite:   DefaultNoBuffWrite,
		RotationTime:  DefaultRotationTime,
		MaxAge:        DefaultMaxAge,
	}
)
