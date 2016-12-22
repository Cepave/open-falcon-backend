package g

import (
	"encoding/json"
	"fmt"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/toolkits/file"
)

type HttpConfig struct {
	Enabled    bool   `json:"enabled"`
	Listen     string `json:"listen"`
	Cookie     string `json:"cookie"`
	ViewPath   string `json:"view_path"`
	StaticPath string `json:"static_path"`
}

type TimeoutConfig struct {
	Conn  int64 `json:"conn"`
	Read  int64 `json:"read"`
	Write int64 `json:"write"`
}

type CacheConfig struct {
	Enabled bool           `json:"enabled"`
	Redis   string         `json:"redis"`
	Idle    int            `json:"idle"`
	Max     int            `json:"max"`
	Timeout *TimeoutConfig `json:"timeout"`
}

type UicConfig struct {
	Addr string `json:"addr"`
	Idle int    `json:"idle"`
	Max  int    `json:"max"`
}

type GraphDBConfig struct {
	Addr           string `json:"addr"`
	Idle           int    `json:"idle"`
	Max            int    `json:"max"`
	Limit          int    `json:"limit"`
	LimitHostGroup int    `json:"limitHostGroup"`
}

type FalconPortalConfig struct {
	Addr  string `json:"addr"`
	Idle  int    `json:"idle"`
	Max   int    `json:"max"`
	Limit int    `json:"limit"`
}

type BossConfig struct {
	Addr string `json:"addr"`
	Idle int    `json:"idle"`
	Max  int    `json:"max"`
}

type ShortcutConfig struct {
	FalconPortal     string `json:"falconPortal"`
	FalconDashboard  string `json:"falconDashboard"`
	GrafanaDashboard string `json:"grafanaDashboard"`
	FalconAlarm      string `json:"falconAlarm"`
}

type LdapConfig struct {
	Enabled    bool     `json:"enabled"`
	Addr       string   `json:"addr"`
	BindDN     string   `json:"bindDN"`
	BaseDN     string   `json:"baseDN"`
	BindPasswd string   `json:"bindPasswd"`
	UserField  string   `json:"userField"`
	Attributes []string `json:"attributes"`
}

type ApiConfig struct {
	Key      string `json:"key"`
	Redirect string `json:"redirect"`
	Login    string `json:"login"`
	Access   string `json:"access"`
	Role     string `json:"role"`
	Logout   string `json:"logout"`
	Name     string `json:"name"`
	Token    string `json:"token"`
	Map      string `json:"map"`
	Contact  string `json:"contact"`
}

type GraphConfig struct {
	ConnTimeout int32             `json:"connTimeout"`
	CallTimeout int32             `json:"callTimeout"`
	MaxConns    int32             `json:"maxConns"`
	MaxIdle     int32             `json:"maxIdle"`
	Replicas    int32             `json:"replicas"`
	Cluster     map[string]string `json:"cluster"`
}

type GrpcConfig struct {
	Enabled bool `json:"enabled"`
	Port    int  `json:"port"`
}

type MqConfig struct {
	Enabled  bool   `json:"enabled"`
	Queue    string `json:"queue"`
	Consumer string `json:"consumer"`
}

type GlobalConfig struct {
	Log          string              `json:"log"`
	Company      string              `json:"company"`
	Cache        *CacheConfig        `json:"cache"`
	Http         *HttpConfig         `json:"http"`
	Salt         string              `json:"salt"`
	CanRegister  bool                `json:"canRegister"`
	Ldap         *LdapConfig         `json:"ldap"`
	Uic          *UicConfig          `json:"uic"`
	GraphDB      *GraphDBConfig      `json:"graphdb"`
	BossDB       *BossConfig         `json:"boss"`
	FalconPortal *FalconPortalConfig `json:"falcon_portal"`
	Shortcut     *ShortcutConfig     `json:"shortcut"`
	Api          *ApiConfig          `json:"api"`
	Graph        *GraphConfig        `json:"graph"`
	Grpc         *GrpcConfig         `json:"grpc"`
	Mq           *MqConfig           `json:"mq"`
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

func ParseConfig(cfg string) error {
	if cfg == "" {
		return fmt.Errorf("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		return fmt.Errorf("config file %s is nonexistent", cfg)
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		return fmt.Errorf("read config file %s fail %s", cfg, err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		return fmt.Errorf("parse config file %s fail %s", cfg, err)
	}

	configLock.Lock()
	defer configLock.Unlock()

	config = &c

	log.Println("read config file:", cfg, "successfully")
	return nil
}
