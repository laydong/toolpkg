// 链路追踪

package tracex

import (
	"github.com/laydong/toolpkg"
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
		if toolpkg.TraceMod() != 0 {
			var err error
			switch toolpkg.TraceType() {
			case TraceTypeZipkin:
				tracer = newZkTracer(toolpkg.AppName(), utils.GetClientIp(), toolpkg.TraceAddr(), toolpkg.TraceMod())
				if err != nil {
					return nil, err
				}
				log.Printf("[app] tracer success")
			case TraceTypeJaeger:
				tracer = newJTracer(toolpkg.AppName(), toolpkg.TraceAddr(), toolpkg.TraceMod())
				if err != nil {
					return nil, err
				}
				log.Printf("[app] tracer success")
			}
		}
	}

	return tracer, nil
}
