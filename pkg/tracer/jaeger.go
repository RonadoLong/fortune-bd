package tracer

import (
	"github.com/uber/jaeger-client-go"
	"io"
	"time"

	"github.com/uber/jaeger-client-go/config"

	"github.com/opentracing/opentracing-go"
)

// NewTracer .
func NewTracer(servicename string, addr string) (tracer opentracing.Tracer, closer io.Closer, err error) {
	cfg := config.Configuration{
		ServiceName: servicename,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}
	var sender jaeger.Transport
	var reporter jaeger.Reporter
	sender, err = jaeger.NewUDPTransport(addr, 0)
	if err != nil {
		return
	}
	reporter = jaeger.NewRemoteReporter(sender)
	tracer, closer, err = cfg.NewTracer(config.Reporter(reporter))
	return
}
