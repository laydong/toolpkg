package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/laydong/toolpkg/logx"
	"github.com/unrolled/secure"
	"net/http"
)

func LoadTls() gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "localhost:443",
		})
		err := middleware.Process(c.Writer, c.Request)
		if err != nil {
			// 如果出现错误，请不要继续
			c.JSON(
				http.StatusOK,
				map[string]string{
					"code":       "400",
					"data":       "",
					"msg":        err.Error(),
					"request_id": c.GetString(logx.RequestIdKey),
				})
			return
		}
		// 继续往下处理
		c.Next()
	}
}
