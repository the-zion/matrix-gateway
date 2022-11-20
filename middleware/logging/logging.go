package logging

import (
	"github.com/go-kratos/kratos/contrib/log/tencent/v2"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"os"
	"strings"
	"time"

	config "github.com/go-kratos/gateway/api/gateway/config/v1"
	"github.com/go-kratos/gateway/middleware"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	tencentLogger tencent.Logger
	logger        = log.GetLogger()
)

func init() {
	middleware.Register("logging", Middleware)
}

// Middleware is a logging middleware.
func Middleware(c *config.Middleware) (middleware.Middleware, error) {
	if tencentLogger != nil {
		tencentLogger.Close()
	}
	var err error
	switch os.Getenv("LOG_SELECT") {
	case "tencent":
		tencentLogger, err = tencent.NewLogger(
			tencent.WithEndpoint(os.Getenv("TENCENT_LOG_HOST")),
			tencent.WithAccessKey(os.Getenv("TENCENT_LOG_ACCESSKEY")),
			tencent.WithAccessSecret(os.Getenv("TENCENT_LOG_ACCESSSECRET")),
			tencent.WithTopicID(os.Getenv("TENCENT_LOG_TOPIC_ID")),
		)
		if err != nil {
			log.Fatalf("failed to new tencent logger: %v", err)
		}
		tencentLogger.GetProducer().Start()
		logger = tencentLogger
	}
	return func(next http.RoundTripper) http.RoundTripper {
		return middleware.RoundTripperFunc(func(req *http.Request) (reply *http.Response, err error) {
			startTime := time.Now()
			reply, err = next.RoundTrip(req)
			level := log.LevelInfo
			code := http.StatusBadGateway
			errMsg := ""
			if err != nil {
				level = log.LevelError
				errMsg = err.Error()
			} else {
				code = reply.StatusCode
			}
			ctx := req.Context()
			nodes, _ := middleware.RequestBackendsFromContext(ctx)

			var traceId, spanId string
			if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
				traceId = span.TraceID().String()
			}
			if span := trace.SpanContextFromContext(ctx); span.HasSpanID() {
				spanId = span.TraceID().String()
			}
			log.WithContext(ctx, log.NewFilter(logger, log.FilterFunc(
				func(level log.Level, keyvals ...interface{}) bool {
					if keyvals == nil {
						return false
					}

					if level == log.LevelError {
						return false
					}

					for i := 0; i < len(keyvals); i++ {
						if keyvals[i] == "code" && keyvals[i+1] != 200 {
							return false
						}
					}
					return true
				}))).Log(level,
				"source", "accesslog",
				"trace.id", traceId,
				"span.id", spanId,
				"host", req.Host,
				"method", req.Method,
				"scheme", req.URL.Scheme,
				"path", req.URL.Path,
				"query", req.URL.RawQuery,
				"code", code,
				"error", errMsg,
				"latency", time.Since(startTime).Seconds(),
				"backend", strings.Join(nodes, ","),
			)
			return reply, err
		})
	}, nil
}
