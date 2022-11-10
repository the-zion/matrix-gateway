package auth

import (
	"bytes"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	config "github.com/go-kratos/gateway/api/gateway/config/v1"
	v1 "github.com/go-kratos/gateway/api/gateway/middleware/auth/v1"
	"github.com/go-kratos/gateway/middleware"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

const (
	reasonUnAuthorized string = `{"code":401, "reason":"UNAUTHORIZED"}`
	reasonTokenExpired string = `{"code":401, "reason":"TOKEN_EXPIRED"}`
)

var (
	LOG                       = log.NewHelper(log.With(log.GetLogger(), "source", "accesslog"))
	ErrMissingJwtToken        = errors.Unauthorized("TOKEN_MISSING", "JWT token is missing")
	ErrTokenInvalid           = errors.Unauthorized("TOKEN_INVALID", "Token is invalid")
	ErrTokenExpired           = errors.Unauthorized("TOKEN_EXPIRED", "JWT token has expired")
	ErrTokenParseFail         = errors.Unauthorized("TOKEN_PARSE_FAIL", "Fail to parse JWT token ")
	ErrUnSupportSigningMethod = errors.Unauthorized("UN_SUPPORT_SIGNING_METHOD", "Wrong signing method")
)

func init() {
	middleware.Register("auth", Middleware)
}

func claimsFunc() jwt.Claims {
	return &jwt.MapClaims{}
}

func newResponse(statusCode int, header http.Header, data []byte) *http.Response {
	return &http.Response{
		StatusCode:    statusCode,
		ContentLength: int64(len(data)),
		Body:          ioutil.NopCloser(bytes.NewReader(data)),
	}
}

func jwtCheck(options *v1.Auth, req *http.Request) error {
	url := req.RequestURI
	if _, ok := options.JwtCheckRouters[url]; ok {
		auths := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
		if len(auths) != 2 || !strings.EqualFold(auths[0], "Bearer") {
			//return newResponse(200, req.Header, []byte(reasonUnAuthorized)), ErrMissingJwtToken
			return ErrMissingJwtToken
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
				return errors.Unauthorized("UNAUTHORIZED", err.Error())
			}
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return ErrTokenInvalid
			}
			if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return ErrTokenExpired
			}
			return ErrTokenParseFail
		}
		if !tokenInfo.Valid {
			return ErrTokenInvalid
		}
		if tokenInfo.Method != jwt.SigningMethodHS256 {
			return ErrUnSupportSigningMethod
		}
		token := *((tokenInfo.Claims.(jwt.Claims)).(*jwt.MapClaims))
		req.Header.Set("uuid", token["uuid"].(string))
	}
	return nil
}

func casbinAuth(enforcer *casbin.Enforcer, req *http.Request) (bool, error) {
	uuid := req.Header.Get("uuid")
	url := req.RequestURI
	method := req.Method
	return enforcer.Enforce(uuid, url, method)
}

func getRequestPublicIp(req *http.Request) string {
	var ip string
	for _, ip = range strings.Split(req.Header.Get("X-Forwarded-For"), ",") {
		if ip = strings.TrimSpace(ip); ip != "" && !IsInternalIP(net.ParseIP(ip)) {
			return ip
		}
	}

	if ip = strings.TrimSpace(req.Header.Get("X-Real-Ip")); ip != "" && !IsInternalIP(net.ParseIP(ip)) {
		return ip
	}

	if ip, _, _ = net.SplitHostPort(req.RemoteAddr); !IsInternalIP(net.ParseIP(ip)) {
		return ip
	}

	return ip
}

func IsInternalIP(IP net.IP) bool {
	if IP.IsLoopback() {
		return true
	}
	if ip4 := IP.To4(); ip4 != nil {
		return ip4[0] == 10 ||
			(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) ||
			(ip4[0] == 169 && ip4[1] == 254) ||
			(ip4[0] == 192 && ip4[1] == 168)
	}
	return false
}

func getRealIp(req *http.Request) {
	realIp := getRequestPublicIp(req)
	req.Header.Set("realIp", realIp)
}

func Middleware(c *config.Middleware) (middleware.Middleware, error) {
	options := &v1.Auth{}
	if c.Options != nil {
		if err := anypb.UnmarshalTo(c.Options, options, proto.UnmarshalOptions{Merge: true}); err != nil {
			return nil, err
		}
	}

	modelBox := options.Casbin.Model
	modelByte := []byte(strings.Join(modelBox, "\n"))
	if err := ioutil.WriteFile("./model.conf", modelByte, 0644); err != nil {
		return nil, err
	}

	policyBox := options.Casbin.Policy
	policyByte := []byte(strings.Join(policyBox, "\n"))
	if err := ioutil.WriteFile("./policy.csv", policyByte, 0644); err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer("./model.conf", "./policy.csv")
	if err != nil {
		return nil, err
	}

	enforcer.AddNamedMatchingFunc("g", "KeyMatch2", util.KeyMatch2)
	enforcer.EnableLog(true)

	return func(next http.RoundTripper) http.RoundTripper {
		return middleware.RoundTripperFunc(func(req *http.Request) (reply *http.Response, err error) {
			err = jwtCheck(options, req)
			if err != nil {
				LOG.Errorf("host: %s, method: %s, requestUrl: %s, error: %v", req.Host, req.Method, req.RequestURI, err)
				if (err.(*errors.Error)).Reason == "TOKEN_EXPIRED" {
					return newResponse(401, req.Header, []byte(reasonTokenExpired)), nil
				} else {
					return newResponse(401, req.Header, []byte(reasonUnAuthorized)), nil
				}
			}

			ok, err := casbinAuth(enforcer, req)
			if err != nil {
				return nil, err
			}

			if ok == false {
				return newResponse(401, req.Header, []byte(reasonUnAuthorized)), nil
			}

			getRealIp(req)

			reply, err = next.RoundTrip(req)
			return reply, err
		})
	}, nil
}
