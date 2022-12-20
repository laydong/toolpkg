package httpx

import (
	"github.com/gin-gonic/gin"
	"github.com/laydong/toolpkg/grace"
	"github.com/laydong/toolpkg/pprofx"
	"net/http"
	"runtime"
	"time"
)

// WebServer 基于http协议的服务
// 这里的实现是基于gin框架, 封装了gin的所有的方法
// gin 的核心是高效路由, 但是gin.Engine和gin.IRouter(s)的高耦合让我们无法复用, gin的作者认为它的路由就是引擎吧
type WebServer struct {
	// 重写所有的路由相关的方法
	*WebRoute
	// 继承引擎本身的其他方法
	*gin.Engine
}

// NewWebServer 创建WebServer
func NewWebServer(mode string) *WebServer {
	gin.SetMode(mode)

	server := &WebServer{
		Engine: gin.New(),
		WebRoute: &WebRoute{
			root: true,
		},
	}
	server.WebRoute.server = server
	server.WebRoute.RouterGroup = &server.Engine.RouterGroup

	if mode == "debug" {
		pprofx.Wrap(server.Engine)
	}

	return server
}

// RouterRegister 路由注册者
type RouterRegister func(*WebServer)

// type RouterRegister func(WebRouter)

// RegisterRouter 注册路由
func (webServer *WebServer) Register(rr RouterRegister) {
	rr(webServer)
}

const (
	defaultReadTimeout  = time.Second * 3
	defaultWriteTimeout = time.Second * 3
)

// RunGrace 实现Server接口
func (webServer *WebServer) RunGrace(addr string, timeouts ...time.Duration) error {
	readTimeout, writeTimeout := defaultReadTimeout, defaultWriteTimeout
	if len(timeouts) > 0 {
		readTimeout = timeouts[0]
		if len(timeouts) > 1 {
			writeTimeout = timeouts[1]
		}
	}

	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Handler:      webServer.Engine,
	}
	if runtime.GOOS == "windows" {
		return server.ListenAndServe()
	}
	return grace.Serve(server)
}

// Delims 设置模板的分解符
// 重写gin方法
func (webServer *WebServer) Delims(left, right string) *WebServer {
	webServer.Engine.Delims(left, right)
	return webServer
}

// SecureJsonPrefix sets the secureJsonPrefix used in Context.SecureJSON.
// 重写gin方法
func (webServer *WebServer) SecureJsonPrefix(prefix string) *WebServer {
	webServer.SecureJsonPrefix(prefix)
	return webServer
}

// HandleContext re-enter a contextx that has been rewritten.
// This can be done by setting c.Request.URL.Path to your new target.
// Disclaimer: You can loop yourself to death with this, use wisely.
func (webServer *WebServer) HandleContext(wc *WebContext) {
	webServer.Engine.HandleContext(wc.Context)
}

// NoRoute adds handlers for NoRoute. It return a 404 code by default.
// 重写gin方法
func (webServer *WebServer) NoRoute(handlers ...WebHandlerFunc) {
	webServer.Engine.NoRoute(decorateWebHandlers(handlers)...)
}

// NoMethod sets the handlers called when... TODO.
// 重写gin方法
func (webServer *WebServer) NoMethod(handlers ...WebHandlerFunc) {
	webServer.Engine.NoMethod(decorateWebHandlers(handlers)...)
}

// Use adds middleware to the group, see example code in github.
func (webServer *WebServer) Use(middleware ...WebHandlerFunc) WebRouter {
	return webServer.WebRoute.Use(middleware...)
}
