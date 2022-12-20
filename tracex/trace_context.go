package tracex

import (
	"github.com/laydong/toolpkg/metautils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"log"
)

// TracerContext 链路
type TracerContext interface {
	SpanFinish(span opentracing.Span)

	// SpanStart 开启子span
	SpanStart(name string) opentracing.Span

	// SpanInject 注入请求
	SpanInject(md metautils.NiceMD)
}

func (ctx *TraceContext) SpanFinish(span opentracing.Span) {
	if span != nil {
		span.Finish()
	}
}

func (ctx *TraceContext) SpanStart(name string) opentracing.Span {
	if t, err := getTracer(); err == nil {
		if t != nil {
			return t.StartSpan(name, opentracing.FollowsFrom(ctx.TopSpan.Context()))
		}
	}
	return nil
}

// SpanInject 将span注入到request
func (ctx *TraceContext) SpanInject(md metautils.NiceMD) {
	if t, err := getTracer(); err == nil {
		if t != nil {
			err = t.Inject(ctx.TopSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(md))
			if err != nil {
				log.Printf("SpanInject, err: %s", err.Error())
			}
		}
	}
}

// TraceContext trace
type TraceContext struct {
	TopSpan opentracing.Span
}

var _ TracerContext = &TraceContext{}

// NewTraceContext new traceCtx
func NewTraceContext(name string, headers map[string][]string) *TraceContext {
	ctx := &TraceContext{}

	if t, err := getTracer(); err == nil {
		if t != nil {
			if len(headers) == 0 {
				ctx.TopSpan = t.StartSpan(name)
			} else {
				spanCtx, errno := t.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(headers))
				if errno != nil {
					ctx.TopSpan = t.StartSpan(name)
				} else {
					ctx.TopSpan = t.StartSpan(name, ext.RPCServerOption(spanCtx))
				}
			}
		}
	}

	return ctx
}
