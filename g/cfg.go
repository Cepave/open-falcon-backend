package g

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/toolkits/file"
)

type HttpConfig struct {
	Debug   bool   `json:"debug"`
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type QueueConfig struct {
	Sms        string `json:"sms"`
	Mail       string `json:"mail"`
	QQ         string `json:"qq"`
	Serverchan string `json:"serverchan"`
}

type FalconPortalConfig struct {
	Addr string `json:"addr"`
	Idle int    `json:"idle"`
	Max  int    `json:"max"`
}

type RedisConfig struct {
	Addr                string   `json:"addr"`
	MaxIdle             int      `json:"maxIdle"`
	HighQueues          []string `json:"highQueues"`
	LowQueues           []string `json:"lowQueues"`
	UserSmsQueue        string   `json:"userSmsQueue"`
	UserMailQueue       string   `json:"userMailQueue"`
	UserQQQueue         string   `json:"userQQQueue"`
	UserServerchanQueue string   `json:"userServerchanQueue"`
}

type ApiConfig struct {
	Portal string `json:"portal"`
	Uic    string `json:"uic"`
	Links  string `json:"links"`
}

type UicConfig struct {
	Addr string `json:"addr"`
	Idle int    `json:"idle"`
	Max  int    `json:"max"`
}

type ShortcutConfig struct {
	FalconPortal     string `json:"falconPortal"`
	FalconDashboard  string `json:"falconDashboard"`
	GrafanaDashboard string `json:"grafanaDashboard"`
	FalconAlarm      string `json:"falconAlarm"`
	FalconUIC        string `json:"falconUIC"`
}
type GlobalConfig struct {
	Debug        bool                `json:"debug"`
	UicToken     string              `json:"uicToken"`
	Http         *HttpConfig         `json:"http"`
	FalconPortal *FalconPortalConfig `json:"falcon_portal"`
	Queue        *QueueConfig        `json:"queue"`
	Redis        *RedisConfig        `json:"redis"`
	Api          *ApiConfig          `json:"api"`
	Shortcut     *ShortcutConfig     `json:"shortcut"`
	Uic          *UicConfig          `json:"uic"`
	RedirectUrl  string              `json:"redirectUrl"`
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

	configLock.Lock()
	defer configLock.Unlock()
	config = &c
	log.Println("read config file:", cfg, "successfully")
}
