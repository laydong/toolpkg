package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/laydong/toolpkg/logx"
	"net/http"
	"strconv"
)

// NoMethodHandle 处理未知请求方式
func NoMethodHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			map[string]string{
				"code":       strconv.Itoa(400),
				"data":       "",
				"msg":        "当前请求方式不存在",
				"request_id": c.GetString(logx.RequestIdKey),
			})
		c.Abort()
		return
	}
}
