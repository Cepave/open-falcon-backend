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
)

// Constructs configuration of gock, which
func NewGockConfig() *GockConfig {
	port := rand.Int31n(1000) + 30000
	host := fmt.Sprintf("test-pc%03d.gock.com", rand.Int31n(999)+1)
	url := fmt.Sprintf("http://%s:%d", host, port)

	return &GockConfig{
		Host: host,
		Port: uint16(port),
		Url:  url,
		GentlemanT: &implGentlemanT{
			url: url,
		},
	}
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
	Url        string
	GentlemanT GentlemanT
}

func (c *GockConfig) NewHttpConfig() *client.HttpClientConfig {
	config := client.NewDefaultConfig()
	config.Url = c.Url
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

func (c *GockConfig) New() *gock.Request {
	logger.Infof("New Gock request: %s", c.Url)
	return gock.New(c.Url)
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
