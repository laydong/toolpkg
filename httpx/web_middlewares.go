package httpx

import (
	"bytes"
	"github.com/laydong/toolpkg"
	"github.com/laydong/toolpkg/utils"
	"io/ioutil"
	"strings"
)

// 不需要打印入参和出参的路由
// 不需要打印入参和出参的前缀
// 不需要打印入参和出参的后缀
type logParams struct {
	NoLogParams       map[string]string
	NoLogParamsPrefix []string
	NoLogParamsSuffix []string
}

// NoLogParamsRules 不想打印的路由分组
var NoLogParamsRules logParams

// CheckNoLogParams 判断是否需要打印入参出参日志, 不需要打印返回true
func CheckNoLogParams(origin string) bool {
	if len(NoLogParamsRules.NoLogParams) > 0 {
		if _, ok := NoLogParamsRules.NoLogParams[origin]; ok {
			return true
		}
	}

	if len(NoLogParamsRules.NoLogParamsPrefix) > 0 {
		for _, v := range NoLogParamsRules.NoLogParamsPrefix {
			if strings.HasPrefix(origin, v) {
				return true
			}
		}
	}

	if len(NoLogParamsRules.NoLogParamsSuffix) > 0 {
		for _, v := range NoLogParamsRules.NoLogParamsSuffix {
			if strings.HasSuffix(origin, v) {
				return true
			}
		}
	}

	return false
}

// ginInterceptor 记录框架出入参, 开启链路追踪
func ginInterceptor(ctx *WebContext) {
	w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: ctx.Writer}
	ctx.Writer = w
	if toolpkg.ParamLog() && !CheckNoLogParams(ctx.Request.RequestURI) {
		requestData, _ := ctx.GetRawData()
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestData))
		ctx.InfoF("%s", string(requestData),
			ctx.Field("header", utils.GetString(ctx.Request.Header)),
			ctx.Field("path", ctx.Request.RequestURI),
			ctx.Field("protocol", protocol),
			ctx.Field("title", "入参"))
	}

	ctx.Next()

	if toolpkg.ParamLog() && !CheckNoLogParams(ctx.Request.RequestURI) {
		ctx.InfoF("%s", w.body.String(), ctx.Field("title", "出参"))
	}
	ctx.SpanFinish(ctx.TopSpan)
}
