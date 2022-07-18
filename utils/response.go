package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code      int         `json:"code"`
	Data      interface{} `json:"data"`
	Msg       string      `json:"msg"`
	RequestID string      `json:"request_id"`
}

func Result(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(
		http.StatusOK,
		Response{
			code,
			data,
			msg,
			c.GetString("request_id"),
		})
}

func Ok(c *gin.Context) {
	Result(http.StatusOK, map[string]interface{}{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(http.StatusOK, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(http.StatusOK, data, "操作成功", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(http.StatusOK, data, message, c)
}

func Fail(c *gin.Context) {
	Result(http.StatusBadRequest, map[string]interface{}{}, "操作失败", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(http.StatusBadRequest, map[string]interface{}{}, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(http.StatusBadRequest, data, message, c)
}

func FailAuthMessage(message string, c *gin.Context) {
	Result(http.StatusUnauthorized, map[string]interface{}{}, message, c)
}

func FailNotMessage(message string, c *gin.Context) {
	Result(http.StatusNotFound, map[string]interface{}{}, message, c)
}
