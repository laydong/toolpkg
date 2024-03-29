// 链路追踪

package tracex

import (
	"github.com/laydong/toolpkg/logx"
	"github.com/laydong/toolpkg/utils"
	"github.com/opentracing/opentracing-go"
	"log"
)

const (
	TraceTypeJaeger = "jaeger"
	TraceTypeZipkin = "zipkin"
)

// tracer 全局单例变量
var tracer opentracing.Tracer

// InitTrace 初始化trace
func getTracer() (opentracing.Tracer, error) {
	if tracer == nil {
		var err error
		switch logx.DefaultTraceType {
		case TraceTypeZipkin:
			tracer = newZkTracer(logx.DefaultAppName, utils.GetClientIp(), logx.DefaultTraceAddr, logx.DefaultTraceMod)
			if err != nil {
				return nil, err
			}
			log.Printf("[app] tracer success")
		case TraceTypeJaeger:
			tracer = newJTracer(logx.DefaultAppName, logx.DefaultTraceAddr, logx.DefaultTraceMod)
			if err != nil {
				return nil, err
			}
			log.Printf("[app] tracer success")
		}
	}

	return tracer, nil
}
