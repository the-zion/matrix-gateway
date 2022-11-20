package healthcheck

import (
	config "github.com/go-kratos/gateway/api/gateway/config/v1"
	"github.com/go-kratos/gateway/middleware"
	"net/http"
)

func init() {
	middleware.Register("health", Middleware)
}

func Middleware(c *config.Middleware) (middleware.Middleware, error) {
	return func(next http.RoundTripper) http.RoundTripper {
		return middleware.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			if req.RequestURI == "/health" {
				return &http.Response{StatusCode: 200, Header: req.Header}, nil
			}
			return next.RoundTrip(req)
		})
	}, nil
}
