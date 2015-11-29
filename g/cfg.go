package g

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/toolkits/file"
)

type HttpConfig struct {
	Enable bool   `json:"enable"`
	Listen string `json:"listen"`
}

type GraphConfig struct {
	ConnTimeout int32             `json:"connTimeout"`
	CallTimeout int32             `json:"callTimeout"`
	MaxConns    int32             `json:"maxConns"`
	MaxIdle     int32             `json:"maxIdle"`
	Replicas    int32             `json:"replicas"`
	Cluster     map[string]string `json:"cluster"`
}

type ApiConfig struct {
	Name      string    `json:"name"`
	Token     string    `json:"token"`
	Event     string    `json:"event"`
	Map       string    `json:"map"`
	Geo       string    `json:"geo"`
	Query     string    `json:"query"`
	Dashboard string    `json:"dashboard"`
}

type DbConfig struct {
	Addr     string `json:"addr"`
	Idle     int    `json:"idle"`
	Max      int    `json:"max"`
}

type GlobalConfig struct {
	Debug    bool            `json:"debug"`
	Http     *HttpConfig     `json:"http"`
	Graph    *GraphConfig    `json:"graph"`
	Api      *ApiConfig      `json:"api"`
	Db       *DbConfig       `json:"db"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("config file not specified: use -c $filename")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file specified not found:", cfg)
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file", cfg, "error:", err.Error())
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file", cfg, "error:", err.Error())
	}

	// set config
	configLock.Lock()
	defer configLock.Unlock()
	config = &c

	log.Println("g.ParseConfig ok, file", cfg)
}
