package nacos

import (
	"github.com/go-kratos/gateway/discovery"
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"net/url"
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

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(ip, uint64(port)),
	}

	cc := constant.ClientConfig{
		UpdateCacheWhenEmpty: true,
		NotLoadCacheAtStart:  true,
		RotateTime:           "1h",
		MaxAge:               3,
		LogLevel:             "warn",
		NamespaceId:          "public",
		TimeoutMs:            5000,
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
