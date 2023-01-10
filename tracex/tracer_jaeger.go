package tracex

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerCfg "github.com/uber/jaeger-client-go/config"
	jaegerLog "github.com/uber/jaeger-client-go/log"
)

func newJTracer(serviceName, addr string, mod float64) opentracing.Tracer {
	var cfg = jaegerCfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegerCfg.SamplerConfig{
			Type:  jaeger.SamplerTypeProbabilistic,
			Param: mod,
		},
		Reporter: &jaegerCfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: addr,
		},
	}

	jLogger := jaegerLog.StdLogger
	t, _, _ := cfg.NewTracer(
		jaegerCfg.Logger(jLogger),
	)

	return t
}
