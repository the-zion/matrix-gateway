package auth

import (
	"bytes"
	"fmt"
	config "github.com/go-kratos/gateway/api/gateway/config/v1"
	v1 "github.com/go-kratos/gateway/api/gateway/middleware/auth/v1"
	"github.com/go-kratos/gateway/middleware"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	reasonUnAuthorized string = `{"code":401, "message":"UNAUTHORIZED"}`
	reasonTokenExpired string = `{"code":401, "message":"TOKEN_EXPIRED"}`
)

var (
	LOG                       = log.NewHelper(log.With(log.GetLogger(), "source", "accesslog"))
	ErrMissingJwtToken        = errors.Unauthorized("TOKEN_MISSING", "JWT token is missing")
	ErrTokenInvalid           = errors.Unauthorized("TOKEN_INVALID", "Token is invalid")
	ErrTokenExpired           = errors.Unauthorized("TOKEN_EXPIRED", "JWT token has expired")
	ErrTokenParseFail         = errors.Unauthorized("TOKEN_PARSE_FAIL", "Fail to parse JWT token ")
	ErrUnSupportSigningMethod = errors.Unauthorized("UNAUTHORIZED", "Wrong signing method")
	ErrWrongContext           = errors.Unauthorized("UNAUTHORIZED", "Wrong context for middleware")
	ErrNeedTokenProvider      = errors.Unauthorized("UNAUTHORIZED", "Token provider is missing")
	ErrSignToken              = errors.Unauthorized("UNAUTHORIZED", "Can not sign token.Is the key correct?")
	ErrGetKey                 = errors.Unauthorized("UNAUTHORIZED", "Can not get key while signing token")
)

func init() {
	middleware.Register("auth", Middleware)
}

func claimsFunc() jwt.Claims {
	return &jwt.MapClaims{}
}

func newResponse(statusCode int, header http.Header, data []byte) *http.Response {
	return &http.Response{
		Header:        header,
		StatusCode:    statusCode,
		ContentLength: int64(len(data)),
		Body:          ioutil.NopCloser(bytes.NewReader(data)),
	}
}

func Middleware(c *config.Middleware) (middleware.Middleware, error) {
	options := &v1.Auth{}
	if c.Options != nil {
		if err := anypb.UnmarshalTo(c.Options, options, proto.UnmarshalOptions{Merge: true}); err != nil {
			return nil, err
		}
	}
	return func(next http.RoundTripper) http.RoundTripper {
		return middleware.RoundTripperFunc(func(req *http.Request) (reply *http.Response, err error) {
			url := req.RequestURI
			if _, ok := options.JwtCheckRouters[url]; ok {
				auths := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
				if len(auths) != 2 || !strings.EqualFold(auths[0], "Bearer") {
					return newResponse(200, req.Header, []byte(reasonUnAuthorized)), ErrMissingJwtToken
				}
				jwtToken := auths[1]
				var (
					tokenInfo *jwt.Token
					err       error
				)
				tokenInfo, err = jwt.ParseWithClaims(jwtToken, claimsFunc(), func(token *jwt.Token) (interface{}, error) {
					return []byte(options.Key), nil
				})
				if err != nil {
					ve, ok := err.(*jwt.ValidationError)
					if !ok {
						return newResponse(200, req.Header, []byte(reasonUnAuthorized)), errors.Unauthorized("UNAUTHORIZED", err.Error())
					}
					if ve.Errors&jwt.ValidationErrorMalformed != 0 {
						return nil, ErrTokenInvalid
					}
					if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
						return nil, ErrTokenExpired
					}
					return nil, ErrTokenParseFail
				}
				if !tokenInfo.Valid {
					return nil, ErrTokenInvalid
				}
				if tokenInfo.Method != jwt.SigningMethodHS256 {
					return nil, ErrUnSupportSigningMethod
				}
				token := *((tokenInfo.Claims.(jwt.Claims)).(*jwt.MapClaims))
				fmt.Println(token["uuid"])
				req.Header.Set("uuid", token["uuid"].(string))
			}
			reply, err = next.RoundTrip(req)
			return reply, err
		})
	}, nil
}
