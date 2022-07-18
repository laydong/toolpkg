package glogs

import (
	"context"
	"errors"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	rHttp "github.com/openzipkin/zipkin-go/reporter/http"
	"log"
	"net/http"
	"time"
)

// Tracer 引擎
var Tracer *zipkin.Tracer

// SpanContextKey ctx key，约定ctx的key名称
var SpanContextKey = "default_app_context_span"

type Trace struct {
	ServiceName     string // 服务名
	ServiceEndpoint string // 当前服务节点
	ZipkinAddr      string // zipkin地址
	Mod             uint64 // 采样率,0==不进行链路追踪，1==全量。值越大，采样率月底，对性能影响越小
}

// InitTrace 初始化trace
func InitTrace(serviceName, serviceEndpoint, zipkinAddr string, mod uint64) error {
	err := GetNewTrace(serviceName, serviceEndpoint, zipkinAddr, mod).InitTracer()
	if err != nil {
		log.Printf("[glogs_trace] open fail,err=%s", err.Error())
		return err
	}
	log.Printf("[glogs_trace] zipkin success")
	return nil
}

// GetNewTrace 获取配置
func GetNewTrace(serviceName, serviceEndpoint, zipkinAddr string, mod uint64) *Trace {
	return &Trace{
		ServiceName:     serviceName,
		ServiceEndpoint: serviceEndpoint,
		ZipkinAddr:      zipkinAddr,
		Mod:             mod,
	}
}

// InitTracer 初始化tracer
func (t *Trace) InitTracer() error {
	var err error
	Tracer, err = t.GetTrace()
	return err
}

// GetTrace 获取tracer
func (t *Trace) GetTrace() (*zipkin.Tracer, error) {
	if t == (&Trace{}) {
		return nil, errors.New("trace is not init")
	}
	// create a reporter to be used by the tracer
	reporter := rHttp.NewReporter(t.ZipkinAddr)
	// set-up the local endpoint for our service
	endpoint, err := zipkin.NewEndpoint(t.ServiceName, t.ServiceEndpoint)
	if err != nil {
		return nil, err
	}
	// set-up our sampling strategy
	sampler := zipkin.NewModuloSampler(t.Mod)
	if t.Mod == 0 {
		sampler = zipkin.NeverSample
	}
	// initialize the tracer
	tracer, err := zipkin.NewTracer(
		reporter,
		zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithSampler(sampler),
	)
	return tracer, err
}

// StartSpan 根据上下文创建span
func StartSpan(name string) zipkin.Span {
	ctx := context.WithValue(context.Background(), SpanContextKey, time.Now().UnixNano())
	if Tracer != nil {
		spanChild := Tracer.StartSpan(name)
		if ctx == nil {
			return spanChild
		}
		spanI := ctx.Value(SpanContextKey)
		if spanContext, ok := spanI.(model.SpanContext); ok {
			spanChild = Tracer.StartSpan(name, zipkin.Parent(spanContext))
		}
		return spanChild
	} else {
		return nil
	}
}

// StartSpanR 根据请求头创建span
func StartSpanR(r *http.Request, name string) zipkin.Span {
	if Tracer != nil {
		if r != (&http.Request{}) {
			return Tracer.StartSpan(name, zipkin.Parent(Tracer.Extract(b3.ExtractHTTP(copyRequest(r)))))
		}
		return Tracer.StartSpan(name)
	} else {
		return nil
	}
}

func StartSpanP(ctx model.SpanContext, name string) zipkin.Span {
	return Tracer.StartSpan(name, zipkin.Parent(ctx))
}

// StopSpan 结束
func StopSpan(span zipkin.Span) {
	if span == nil {
		return
	}
	span.Finish()
}

// Inject 注入span信息到请求头
func Inject(ctx context.Context, r *http.Request) error {
	injector := b3.InjectHTTP(r)
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	spanI := ctx.Value(SpanContextKey)
	if spanContext, ok := spanI.(model.SpanContext); ok {
		err := injector(spanContext)
		return err
	}
	return nil
}

func copyRequest(r *http.Request) *http.Request {
	req := &http.Request{}
	req = r
	return req
}

func GetSpanContextKey() string { return SpanContextKey }
