package tracex

import (
	"github.com/opentracing/opentracing-go"
	zipkinOt "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinHttp "github.com/openzipkin/zipkin-go/reporter/http"
	"log"
)

func newZkTracer(serviceName, serviceEndpoint, addr string, mod float64) opentracing.Tracer {

	// set up a span reporter
	reporter := zipkinHttp.NewReporter(addr)

	// create our local service endpoint
	endpoint, err := zipkin.NewEndpoint(serviceName, serviceEndpoint)
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	// set up our sampling strategy
	sampler, err := zipkin.NewBoundarySampler(mod, 100)
	if err != nil {
		log.Fatalf("unable to set sampling strategy: %+v\n", err)
	}

	// initialize our tracer
	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint), zipkin.WithSampler(sampler))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	// use zipkin-go-opentracing to wrap our tracer
	t := zipkinOt.Wrap(nativeTracer)

	log.Printf("[glogs_trace] zipkin success")
	return t
}
