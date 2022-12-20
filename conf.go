package toolpkg

import (
	"time"
)

var (
	envAppName             = "default-app"
	envAppMode             = "dev"
	envRunMode             = "debug"
	envAppVersion          = "1.0.0"
	envHttpListen          = "0.0.0.0:80"
	envGrpcListen          = "0.0.0.0:10082"
	envAppUrl              = "http://127.0.0.1:80"
	envParamLog            = true
	envLogPath             = "/home/logs/app"
	envLogType             = "console"
	envLogMaxAge           = 7 * 24 * time.Hour
	envLogMaxCount uint    = 30
	envTraceType           = ""
	envTraceAddr           = ""
	envTraceMod    float64 = 0
	envLocalIP             = ""
)

// SetAppName 设置app名称
func SetAppName(appName string) {
	if appName != "" {
		envAppName = appName
	}
}

// AppName 返回当前app名称
func AppName() string {
	return envAppName
}

// SetAppMode 设置当前的环境
func SetAppMode(appMode string) {
	if appMode != "" {
		envAppMode = appMode
	}
}

// AppMode 返回当前的环境
func AppMode() string {
	return envAppMode
}

// SetRunMode 设置运行模式
func SetRunMode(runMode string) {
	if runMode != "" {
		envRunMode = runMode
	}
}

// RunMode 返回当前的运行模式
func RunMode() string {
	return envRunMode
}

// SetAppVersion 设置app的版本号
func SetAppVersion(appVersion string) {
	if appVersion != "" {
		envAppVersion = appVersion
	}
}

// AppVersion 返回app的版本号
func AppVersion() string {
	return envAppVersion
}

// SetHttpListen 设置http监听地址
func SetHttpListen(httpListen string) {
	if httpListen != "" {
		envHttpListen = httpListen
	}
}

// HttpListen 获取http监听地址
func HttpListen() string {
	return envHttpListen
}

// SetGrpcListen 设置rpc监听地址
func SetGrpcListen(grpcListen string) {
	if grpcListen != "" {
		envGrpcListen = grpcListen
	}
}

// GrpcListen 返回rpc监听地址
func GrpcListen() string {
	return envGrpcListen
}

// SetAppUrl 设置app_url
func SetAppUrl(appUrl string) {
	if appUrl != "" {
		envAppUrl = appUrl
	}
}

// AppUrl 返回当前app_url
func AppUrl() string {
	return envAppUrl
}

// SetParamLog 设置是否打印入参和出参
func SetParamLog(ParamLog bool) {
	envParamLog = ParamLog
}

// ParamLog 返回是否打印入参和出参
func ParamLog() bool {
	return envParamLog
}

// SetLogPath 设置日志路径
func SetLogPath(path string) {
	if path != "" {
		envLogPath = path
	}
}

// LogPath 返回日志基本路径
func LogPath() string {
	return envLogPath
}

// SetLogType 设置日志类型
func SetLogType(tp string) {
	if tp != "" {
		envLogType = tp
	}
}

// LogType 返回日志类型
func LogType() string {
	return envLogType
}

// SetLogMaxAge 设置日志默认保留7天
func SetLogMaxAge(maxAge int) {
	if maxAge != 0 {
		envLogMaxAge = time.Duration(maxAge) * 24 * time.Hour
	}
}

// LogMaxAge 返回日志默认保留7天
func LogMaxAge() time.Duration {
	return envLogMaxAge
}

// SetLogMaxCount 设置日志默认限制为30个
func SetLogMaxCount(count int) {
	if count != 0 {
		envLogMaxCount = uint(count)
	}
}

// LogMaxCount 返回日志默认限制为30个
func LogMaxCount() uint {
	return envLogMaxCount
}

func SetTraceType(v string) {
	envTraceType = v
}

func TraceType() string {
	return envTraceType
}

func SetTraceAddr(v string) {
	envTraceAddr = v
}

func TraceAddr() string {
	return envTraceAddr
}

func SetTraceMod(v float64) {
	envTraceMod = v
}

func TraceMod() float64 {
	return envTraceMod
}
