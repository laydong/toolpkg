package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/laydong/toolpkg/logx"
	"net/http"
)

// NotRouter 处理未知请求路由
func NotRouter() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			map[string]string{
				"code":       "400",
				"data":       "",
				"msg":        "当前请求不存在",
				"request_id": c.GetString(logx.RequestIdKey),
			})
		c.Abort()
		return
	}
}
