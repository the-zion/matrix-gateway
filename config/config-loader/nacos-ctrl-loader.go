package ctrlloader

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

type NacosCtrlConfigLoader struct {
	CtrlConfigLoader
}

func NewNacosConfigLoader(rawCtrlService, dstPath string) *NacosCtrlConfigLoader {
	cl := &NacosCtrlConfigLoader{}
	cl.ctrlService = prepareCtrlService(rawCtrlService)
	cl.dstPath = dstPath
	cl.hostname = cl.getHostname()
	cl.advertiseAddr = cl.getAdvertiseAddr()
	return cl
}

func (c *NacosCtrlConfigLoader) Load(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			c.nextCtrlService = true
		}
	}()

	cfgBytes, err := c.load(ctx)
	if err != nil {
		return err
	}

	tmpPath := fmt.Sprintf("%s.%s.tmp", c.dstPath, uuid.New().String())
	if err := ioutil.WriteFile(tmpPath, cfgBytes, 0644); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, c.dstPath); err != nil {
		return err
	}
	return nil
}

func (c *NacosCtrlConfigLoader) load(ctx context.Context) ([]byte, error) {
	params := url.Values{}
	params.Set("dataId", "gateway")
	params.Set("group", "matrix")
	LOG.Infof("%s is requesting config from %s with params: %+v", c.hostname, c.ctrlService, params)
	api, err := c.urlfor("/nacos/v1/cs/configs", params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, api, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}
	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *NacosCtrlConfigLoader) choseCtrlService() string {
	if c.nextCtrlService {
		c.ctrlServiceIdx = (c.ctrlServiceIdx + 1) % len(c.ctrlService)
		c.nextCtrlService = false
		return c.ctrlService[c.ctrlServiceIdx]
	}
	return c.ctrlService[c.ctrlServiceIdx]
}

func (c *NacosCtrlConfigLoader) getAdvertiseAddr() string {
	advAddr := os.Getenv("ADVERTISE_ADDR")
	if advAddr != "" {
		return advAddr
	}
	advDevice := os.Getenv("ADVERTISE_DEVICE")
	if advDevice == "" {
		advDevice = "eth0"
	}
	advAddr, err := c.getIPInterface(advDevice)
	if err != nil {
		LOG.Errorf("%q There was a problem with the IP %+v", c.hostname, err)
		return ""
	}
	LOG.Infof("%s uses IP %s\n", c.hostname, advAddr)
	return advAddr
}

func (c *NacosCtrlConfigLoader) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel
	for {
		if err := c.Load(ctx); err != nil {
			// logging
			continue
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second * 5):
		}
	}
}
