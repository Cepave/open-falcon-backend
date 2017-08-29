package gock

import (
	"fmt"
	"math/rand"

	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugin"
	"gopkg.in/h2non/gock.v1"

	"github.com/Cepave/open-falcon-backend/common/http/client"
	mysql "github.com/Cepave/open-falcon-backend/common/service/mysqlapi"

	tHttp "github.com/Cepave/open-falcon-backend/common/testing/http"
)

// Constructs configuration of gock, which
func NewGockConfig() *GockConfig {
	port := rand.Int31n(1000) + 30000
	host := fmt.Sprintf("test-pc%03d.gock.com", rand.Int31n(999)+1)

	return newGockConfig(host, uint16(port))
}

func NewGockConfigByTestServer() *GockConfig {
	return newGockConfig(tHttp.WebTestServer.GetHost(), tHttp.WebTestServer.GetPort())
}

func newGockConfig(host string, port uint16) *GockConfig {
	newConfig := &GockConfig{
		Host: host,
		Port: port,
	}

	newConfig.GentlemanT = &implGentlemanT{
		url: newConfig.GetUrl(),
	}

	return newConfig
}

// Defines the interface used to ease testing by Gentleman library.
type GentlemanT interface {
	NewClient() *gentleman.Client
	SetupClient(*gentleman.Client) *gentleman.Client
	Plugin() plugin.Plugin
}

type GockConfig struct {
	Host       string
	Port       uint16
	GentlemanT GentlemanT
}

func (c *GockConfig) NewHttpConfig() *client.HttpClientConfig {
	config := client.NewDefaultConfig()
	config.Url = c.GetUrl()
	return config
}

func (c *GockConfig) NewMySqlApiConfig() *mysql.ApiConfig {
	return &mysql.ApiConfig{
		HttpClientConfig: c.NewHttpConfig(),
		Plugins: []plugin.Plugin{
			_gentlemanMockPlugin,
		},
	}
}

func (c *GockConfig) GetUrl() string {
	return fmt.Sprintf("http://%s:%d", c.Host, c.Port)
}

func (c *GockConfig) New() *gock.Request {
	url := c.GetUrl()

	logger.Infof("New Gock request: %s", url)
	return gock.New(url)
}

// Calls gock.Off()
func (c *GockConfig) Off() {
	gock.Off()
}

// Calls gock.EnableNetworking()
func (c *GockConfig) StartRealNetwork() {
	logger.Info("Start Gock Real Network")
	gock.EnableNetworking()
}

// Calls gock.DisableNetworking()
func (c *GockConfig) StopRealNework() {
	logger.Info("Stop Gock Real Network")
	gock.Off()
	gock.DisableNetworking()
}

type implGentlemanT struct {
	url string
}

func (t *implGentlemanT) NewClient() *gentleman.Client {
	return t.SetupClient(gentleman.New())
}
func (t *implGentlemanT) SetupClient(client *gentleman.Client) *gentleman.Client {
	client.BaseURL(t.url).Use(t.Plugin())
	return client
}
func (t *implGentlemanT) Plugin() plugin.Plugin {
	return _gentlemanMockPlugin
}

var _gentlemanMockPlugin = plugin.NewPhasePlugin("before dial", func(ctx *context.Context, h context.Handler) {
	gock.InterceptClient(ctx.Client)
	h.Next(ctx)
})
