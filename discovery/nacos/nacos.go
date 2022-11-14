package nacos

import (
	"github.com/go-kratos/gateway/discovery"
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func init() {
	discovery.Register("nacos", New)
}

func New(dsn *url.URL) (registry.Discovery, error) {
	host := strings.Split(dsn.Host, ":")
	ip := host[0]
	port, err := strconv.Atoi(host[1])
	if err != nil {
		return nil, err
	}

	namespace := os.Getenv("NACOS_NAMESPACE")
	if namespace == "" {
		namespace = "public"
	}

	username := os.Getenv("NACOS_USERNAME")
	if username == "" {
		username = "nacos"
	}

	password := os.Getenv("NACOS_PASSWORD")
	if password == "" {
		password = "nacos"
	}

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(ip, uint64(port)),
	}

	cc := constant.ClientConfig{
		UpdateCacheWhenEmpty: true,
		NotLoadCacheAtStart:  true,
		RotateTime:           "1h",
		MaxAge:               3,
		LogLevel:             "warn",
		NamespaceId:          namespace,
		TimeoutMs:            5000,
		Username:             username,
		Password:             password,
	}

	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ServerConfigs: sc,
			ClientConfig:  &cc,
		},
	)

	if err != nil {
		return nil, err
	}

	return nacos.New(client), nil
}
