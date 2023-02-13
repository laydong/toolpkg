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

func Result(c *gin.Context, code int, data interface{}, msg string) {
	// 开始时间
	c.JSON(
		http.StatusOK,
		Response{
			code,
			data,
			msg,
			c.GetString(RequestIdKey),
		})
}

func Ok(c *gin.Context) {
	Result(c, http.StatusOK, map[string]interface{}{}, "操作成功")
}

func OkWithMessage(c *gin.Context, message string) {
	Result(c, http.StatusOK, map[string]interface{}{}, message)
}

func OkWithData(c *gin.Context, data interface{}) {
	Result(c, http.StatusOK, data, "操作成功")
}

func OkWithDetailed(c *gin.Context, data interface{}, message string) {
	Result(c, http.StatusOK, data, message)
}

func Fail(c *gin.Context) {
	Result(c, http.StatusBadRequest, map[string]interface{}{}, "操作失败")
}

func FailWithMessage(c *gin.Context, message string) {
	Result(c, http.StatusBadRequest, map[string]interface{}{}, message)
}

func FailWithDetailed(c *gin.Context, data interface{}, message string) {
	Result(c, http.StatusBadRequest, data, message)
}

func FailAuthMessage(c *gin.Context, message string) {
	Result(c, http.StatusUnauthorized, map[string]interface{}{}, message)
}

func FailAuthsMessage(c *gin.Context, message string) {
	Result(c, http.StatusForbidden, map[string]interface{}{}, message)
}

func FailNotMessage(c *gin.Context, message string) {
	Result(c, http.StatusNotFound, map[string]interface{}{}, message)
}
