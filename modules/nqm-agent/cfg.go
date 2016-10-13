package main

import (
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	log "github.com/Sirupsen/logrus"
)

type AgentConfig struct {
	PushURL string `json:"pushURL"`
}

type HbsConfig struct {
	RPCServer string        `json:"RPCServer"`
	Interval  time.Duration `json:"interval"`
}

type JSONConfigFile struct {
	Agent        *AgentConfig `json:"agent"`
	Hbs          *HbsConfig   `json:"hbs"`
	Hostname     string       `json:"hostname"`
	IPAddress    string       `json:"ipAddress"`
	ConnectionID string       `json:"connectionID"`
}

type GeneralConfig struct {
	Hostname     string
	IPAddress    string
	ConnectionID string
}

var (
	jsonConfig    atomic.Value // for *JSONConfigFile
	hbsResp       atomic.Value // for receiving model.NqmTaskResponse
	generalConfig *GeneralConfig
)

func publicIP() (string, error) {
	output, err := exec.Command("dig", "+short", "myip.opendns.com", "@resolver1.opendns.com").Output()
	if err != nil {
		return "UNKNOWN", err
	}
	ipStr := strings.TrimSpace(string(output))
	return ipStr, nil
}

func JSONConfig() *JSONConfigFile {
	return jsonConfig.Load().(*JSONConfigFile)
}

func SetJSONConfig(c *JSONConfigFile) {
	jsonConfig.Store(c)
}

func HBSResp() model.NqmTaskResponse {
	return hbsResp.Load().(model.NqmTaskResponse)
}

func SetHBSResp(r model.NqmTaskResponse) {
	hbsResp.Store(r)
}

func hostname() string {
	hostname := JSONConfig().Hostname
	if hostname != "" {
		log.Println("Hostname set in config: [", hostname, "]")
		return hostname
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln("os.Hostname() ERROR:", err)
	}
	// hostname -s
	// -s, --short
	// Display the short host name. This is the host name cut at the first dot.
	hostname = strings.Split(hostname, ".")[0]
	log.Println("Hostname not set in config, using system's hostname...succeeded: [", hostname, "]")

	return hostname
}

func ip() string {
	ip := JSONConfig().IPAddress
	if ip != "" {
		log.Println("IP set in config: [", ip, "]")
		return ip
	}

	ip, err := publicIP()
	if err != nil {
		log.Println("IP not set in config, getting public IP...failed:", err)
	} else {
		log.Println("IP not set in config, getting public IP...succeeded: [", ip, "]")
	}
	return ip
}

func connectionID() string {
	connectionID := JSONConfig().ConnectionID
	if connectionID != "" {
		log.Println("ConnectionID set in config: [", connectionID, "]")
		return connectionID
	}

	// Logically it shouldn't happen because ConnectionID is alwasy generated
	// after Hostname and IPAddress are set.
	if GetGeneralConfig().Hostname == "" || GetGeneralConfig().IPAddress == "" {
		log.Fatalln("ConnectionID not set in config, generating...failed!")
	}

	connectionID = GetGeneralConfig().Hostname + "@" + GetGeneralConfig().IPAddress
	log.Println("ConnectionID not set in config, generating...succeeded: [", connectionID, "]")
	return connectionID
}

func jsonUnmarshaller() *JSONConfigFile {
	var c = JSONConfigFile{
		Agent: &AgentConfig{},
		Hbs:   &HbsConfig{},
	}
	err := vipercfg.Config().Unmarshal(&c)
	if err != nil {
		log.Fatal("Parsing configuration file [", vipercfg.Config().GetString("config"), "] failed:", err)
	}
	log.Println("Reading configuration file [", vipercfg.Config().GetString("config"), "] succeeded")
	return &c
}

func GetGeneralConfig() *GeneralConfig {
	return generalConfig
}

func InitGeneralConfig(cfgFilePath string) {
	var cfg GeneralConfig
	generalConfig = &cfg

	SetJSONConfig(jsonUnmarshaller())
	SetHBSResp(model.NqmTaskResponse{})

	cfg.Hostname = hostname()
	if cfg.IPAddress = ip(); cfg.IPAddress == "UNKNOWN" {
		log.Fatalln("IP can't be \"UNKNOWN\"")
	}
	cfg.ConnectionID = connectionID()
}
