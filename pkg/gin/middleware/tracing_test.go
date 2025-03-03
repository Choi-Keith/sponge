package middleware

import (
	"context"
	"testing"
	"time"

	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/httpcli"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
)

func TestTracing(t *testing.T) {
	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(Tracing("demo"))

	r.GET("/hello", func(c *gin.Context) {
		response.Success(c, "hello world")
	})

	go func() {
		err := r.Run(serverAddr)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Millisecond * 200)
	result := &httpcli.StdResult{}
	err := httpcli.Get(result, requestAddr+"/hello")
	assert.NoError(t, err)
	t.Log(result)
}

type propagators struct {
	embedded.TracerProvider
}

func (p *propagators) Tracer(instrumentationName string, opts ...oteltrace.TracerOption) oteltrace.Tracer {
	return &tracer{}
}

type tracer struct {
	embedded.Tracer
}

func (t *tracer) Start(ctx context.Context, spanName string, opts ...oteltrace.SpanStartOption) (context.Context, oteltrace.Span) {
	return ctx, nil
}

type tracerProvider struct {
}

func (t *tracerProvider) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {

}

func (t *tracerProvider) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	return ctx
}

func (t *tracerProvider) Fields() []string {
	return []string{}
}

func TestWithPropagators(t *testing.T) {
	cfg := &traceConfig{}
	opt := WithPropagators(&tracerProvider{})
	opt(cfg)
}

func TestWithTracerProvider(t *testing.T) {
	cfg := &traceConfig{}
	opt := WithTracerProvider(&propagators{})
	opt(cfg)
}
