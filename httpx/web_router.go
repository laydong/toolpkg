package httpx

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// toGinHandlerFunc 转换为Gin的HandlerFunc
func (hdlr WebHandlerFunc) toGinHandlerFunc() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ctx := NewWebContext(ginCtx)
		hdlr(ctx)
	}
}

func ginWebHandler(ginHandler gin.HandlerFunc) WebHandlerFunc {
	return func(ctx *WebContext) {
		ginHandler(ctx.Context)
	}
}

func decorateWebHandlers(handlers []WebHandlerFunc) []gin.HandlerFunc {
	_handlers := []gin.HandlerFunc{}
	for _, hdlr := range handlers {
		_handlers = append(_handlers, hdlr.toGinHandlerFunc())
	}
	return _handlers
}

var _ WebRouter = &WebRoute{}

// WebRouter interface WebRequest Router
// 它合并了 gin.IRoute 和 gin.IRoutes
type WebRouter interface {
	// Group gin.IRoute.Group
	Group(string, ...WebHandlerFunc) WebRouter

	// Use gin.IRoutes.Use
	Use(...WebHandlerFunc) WebRouter

	Any(pattern string, handlers ...WebHandlerFunc) WebRouter
	GET(pattern string, handlers ...WebHandlerFunc) WebRouter
	POST(pattern string, handlers ...WebHandlerFunc) WebRouter
	DELETE(pattern string, handlers ...WebHandlerFunc) WebRouter
	PATCH(pattern string, handlers ...WebHandlerFunc) WebRouter
	PUT(pattern string, handlers ...WebHandlerFunc) WebRouter
	OPTIONS(pattern string, handlers ...WebHandlerFunc) WebRouter
	HEAD(pattern string, handlers ...WebHandlerFunc) WebRouter

	StaticFile(relativePath, filepath string) WebRouter
	Static(relativePath, root string) WebRouter
	StaticFS(relativePath string, fs http.FileSystem) WebRouter

	// TODO
	// AutoGET(...WebHandlerFunc) WebRouter
	// AutoPOST(...WebHandlerFunc) WebRouter
	// AutoDELETE(...WebHandlerFunc) WebRouter
	// AutoPATCH(...WebHandlerFunc) WebRouter
	// AutoPUT(...WebHandlerFunc) WebRouter
	// AutoOPTIONS(...WebHandlerFunc) WebRouter
	// AutoHEAD(...WebHandlerFunc) WebRouter
	// AutoHandle(string, ...WebHandlerFunc) WebRouter
}

// WebRouterGroup interface
// type WebRouterGroup interface {
// 	WebRouter
// 	Group(string, ...WebHandlerFunc) WebRouterGroup
// }

// WebRoute struct
// 它实现了gin.IRoutes, gin.IRoute
type WebRoute struct {
	RouterGroup *gin.RouterGroup
	server      *WebServer
	root        bool
}

// Group creates a new web router group
func (wrc *WebRoute) Group(pattern string, handlers ...WebHandlerFunc) WebRouter {
	return &WebRoute{
		RouterGroup: wrc.RouterGroup.Group(pattern, decorateWebHandlers(handlers)...),
		server:      wrc.server,
		root:        false,
	}
}

// Use attachs a global middleware to the router. ie. the middleware attached though Use() will be
// included in the handlers chain for every single request. Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware.
func (wrc *WebRoute) Use(middleware ...WebHandlerFunc) WebRouter {
	wrc.RouterGroup.Use(decorateWebHandlers(middleware)...)
	return wrc.returnObject()
}

// Any 注册所有的方法
func (wrc *WebRoute) Any(pattern string, handlers ...WebHandlerFunc) WebRouter {
	wrc.RouterGroup.Any(pattern, decorateWebHandlers(handlers)...)
	return wrc.returnObject()
}

// GET 注册GET方法
func (wrc *WebRoute) GET(pattern string, handlers ...WebHandlerFunc) WebRouter {
	wrc.RouterGroup.GET(pattern, decorateWebHandlers(handlers)...)
	return wrc.returnObject()
}

// POST 注册POST方法
func (wrc *WebRoute) POST(pattern string, handlers ...WebHandlerFunc) WebRouter {
	wrc.RouterGroup.POST(pattern, decorateWebHandlers(handlers)...)
	return wrc.returnObject()
}

// DELETE 注册DELETE方法
func (wrc *WebRoute) DELETE(pattern string, handlers ...WebHandlerFunc) WebRouter {
	wrc.RouterGroup.DELETE(pattern, decorateWebHandlers(handlers)...)
	return wrc.returnObject()
}

// PATCH 注册PATCH方法
func (wrc *WebRoute) PATCH(pattern string, handlers ...WebHandlerFunc) WebRouter {
	wrc.RouterGroup.PATCH(pattern, decorateWebHandlers(handlers)...)
	return wrc.returnObject()
}

// PUT 注册PUT方法
func (wrc *WebRoute) PUT(pattern string, handlers ...WebHandlerFunc) WebRouter {
	wrc.RouterGroup.PUT(pattern, decorateWebHandlers(handlers)...)
	return wrc.returnObject()
}

// OPTIONS 注册OPTIONS方法
func (wrc *WebRoute) OPTIONS(pattern string, handlers ...WebHandlerFunc) WebRouter {
	wrc.RouterGroup.OPTIONS(pattern, decorateWebHandlers(handlers)...)
	return wrc.returnObject()
}

// HEAD 注册HEAD方法
func (wrc *WebRoute) HEAD(pattern string, handlers ...WebHandlerFunc) WebRouter {
	wrc.RouterGroup.HEAD(pattern, decorateWebHandlers(handlers)...)
	return wrc.returnObject()
}

// StaticFile 静态文件
func (wrc *WebRoute) StaticFile(relativePath, filepath string) WebRouter {
	wrc.RouterGroup.StaticFile(relativePath, filepath)
	return wrc.returnObject()
}

// Static 静态文件
func (wrc *WebRoute) Static(relativePath, root string) WebRouter {
	wrc.RouterGroup.Static(relativePath, root)
	return wrc.returnObject()
}

// StaticFS 静态文件
func (wrc *WebRoute) StaticFS(relativePath string, fs http.FileSystem) WebRouter {
	wrc.RouterGroup.StaticFS(relativePath, fs)
	return wrc.returnObject()
}

func (wrc *WebRoute) returnObject() WebRouter {
	if wrc.root {
		return wrc.server
	}
	return wrc
}
