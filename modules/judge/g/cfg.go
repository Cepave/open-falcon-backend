package g

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/toolkits/file"
)

type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type RpcConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type HbsConfig struct {
	Servers  []string `json:"servers"`
	Timeout  int64    `json:"timeout"`
	Interval int64    `json:"interval"`
}

type RedisConfig struct {
	Dsn          string `json:"dsn"`
	MaxIdle      int    `json:"maxIdle"`
	ConnTimeout  int    `json:"connTimeout"`
	ReadTimeout  int    `json:"readTimeout"`
	WriteTimeout int    `json:"writeTimeout"`
}

type AlarmConfig struct {
	Enabled             bool         `json:"enabled"`
	MinInterval         int64        `json:"minInterval"`
	QueuePattern        string       `json:"queuePattern"`
	AllowReSet          bool         `json:"allow_reset"`
	EventsStoreFilePath string       `json:"events_store_file_path"`
	Redis               *RedisConfig `json:"redis"`
}

type GlobalConfig struct {
	Debug     bool         `json:"debug"`
	DebugHost string       `json:"debugHost"`
	RootDir   string       `json:"root_dir"`
	Remain    int          `json:"remain"`
	Http      *HttpConfig  `json:"http"`
	Rpc       *RpcConfig   `json:"rpc"`
	Hbs       *HbsConfig   `json:"hbs"`
	Alarm     *AlarmConfig `json:"alarm"`
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
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	//support develop mode
	if c.RootDir == "" {
		c.RootDir = filepath.Dir(os.Args[0])
	}

	configLock.Lock()
	defer configLock.Unlock()

	config = &c

	log.Println("read config file:", cfg, "successfully")
}
