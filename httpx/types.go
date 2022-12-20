package httpx

import (
	"github.com/gin-gonic/gin"
)

const (
	protocol = "http"
)

// DefaultWebServerMiddlewares 默认的Http Server中间件
// 其实应该保证TowerLogware 不panic，但是无法保证，多一个recovery来保证业务日志崩溃后依旧有访问日志
var DefaultWebServerMiddlewares = []WebHandlerFunc{
	ginInterceptor,
	ginWebHandler(gin.Recovery()),
	recovery,
}

// 拦截到错误后处理span, 记录日志, 然后panic
func recovery(ctx *WebContext) {
	defer func() {
		if err := recover(); err != nil {
			ctx.SpanFinish(ctx.TopSpan)
			ctx.ErrorF("系统错误, err: %v", err)
			panic(err)
		}
	}()
	ctx.Next()
}
