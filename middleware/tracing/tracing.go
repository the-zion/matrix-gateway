package tracing

import (
	"context"
	"fmt"
	config "github.com/go-kratos/gateway/api/gateway/config/v1"
	v1 "github.com/go-kratos/gateway/api/gateway/middleware/tracing/v1"
	"github.com/go-kratos/gateway/middleware"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"net/http"
	"os"
	"sync"
)

var (
	hostname, _ = os.Hostname()
	//defaultTimeout     = time.Duration(10 * time.Second)
	defaultServiceName = "gateway" + "." + hostname
	defaultTracerName  = "gateway" + "." + hostname
)

var globaltp = &struct {
	provider   trace.TracerProvider
	propagator propagation.TextMapPropagator
	initOnce   sync.Once
}{}

func init() {
	middleware.Register("tracing", Middleware)
}

// Middleware is a opentelemetry middleware.
func Middleware(c *config.Middleware) (middleware.Middleware, error) {
	options := &v1.Tracing{}
	if c.Options != nil {
		if err := anypb.UnmarshalTo(c.Options, options, proto.UnmarshalOptions{Merge: true}); err != nil {
			return nil, err
		}
	}
	if globaltp.provider == nil {
		globaltp.initOnce.Do(func() {
			globaltp.provider = newTracerProvider(context.Background(), options)
			//propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
			globaltp.propagator = propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
			otel.SetTracerProvider(globaltp.provider)
			otel.SetTextMapPropagator(globaltp.propagator)
		})
	}
	tracer := otel.Tracer(defaultTracerName)
	return func(next http.RoundTripper) http.RoundTripper {
		return middleware.RoundTripperFunc(func(req *http.Request) (reply *http.Response, err error) {
			req.Header.Set("x-md-service-name", defaultServiceName)
			ctx, span := tracer.Start(
				req.Context(),
				fmt.Sprintf("%s %s", req.Method, req.URL.Path),
				trace.WithSpanKind(trace.SpanKindClient),
			)

			// attributes for each request
			span.SetAttributes(
				semconv.HTTPMethodKey.String(req.Method),
				semconv.HTTPTargetKey.String(req.URL.Path),
				semconv.NetPeerIPKey.String(req.RemoteAddr),
			)

			globaltp.propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))

			defer func() {
				if err != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
				} else {
					span.SetStatus(codes.Ok, "OK")
				}
				if reply != nil {
					span.SetAttributes(semconv.HTTPStatusCodeKey.Int(reply.StatusCode))
				}
				span.End()
			}()
			return next.RoundTrip(req.WithContext(ctx))
		})
	}, nil
}

func newTracerProvider(ctx context.Context, options *v1.Tracing) trace.TracerProvider {
	var (
		//timeout     = defaultTimeout
		serviceName = defaultServiceName
	)

	if appInfo, ok := kratos.FromContext(ctx); ok {
		serviceName = appInfo.Name()
	}

	//if options.Timeout != nil {
	//	timeout = options.Timeout.AsDuration()
	//}

	var sampler sdktrace.Sampler
	if options.SampleRatio == nil {
		sampler = sdktrace.AlwaysSample()
	} else {
		sampler = sdktrace.TraceIDRatioBased(float64(*options.SampleRatio))
	}

	//client := otlptracehttp.NewClient(
	//	otlptracehttp.WithEndpoint(options.HttpEndpoint),
	//	otlptracehttp.WithTimeout(timeout),
	//	otlptracehttp.WithInsecure(),
	//	otlptracehttp.WithURLPath(""),
	//)
	//
	//exporter, err := otlptrace.New(ctx, client)
	//if err != nil {
	//	log.Fatalf("creating OTLP trace exporter: %v", err)
	//}

	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(options.HttpEndpoint)))
	if err != nil {
		log.Fatalf("creating jaeger trace exporter: %v", err)
	}

	// attributes for all requests
	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		attribute.String("exporter", "jaeger"),
		attribute.Float64("float", 312.23),
		attribute.KeyValue{
			Key: "token", Value: attribute.StringValue(options.HttpEndpointToken),
		},
	)

	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resources),
	)
}
