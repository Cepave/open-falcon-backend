package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/toolkits/file"
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
	JSONConfigFile
	hbsResp      atomic.Value // for receiving model.NqmTaskResponse
	Hostname     string
	IPAddress    string
	ConnectionID string
}

var (
	jsonConfig    *JSONConfigFile
	generalConfig *GeneralConfig
	jsonCfgLock   = new(sync.RWMutex)
)

func getBinAbsPath() string {
	bin, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.Fatalln(err)
	}
	return bin
}

func getWorkingDirAbsPath() string {
	return filepath.Dir(getBinAbsPath())
}

func getCfgAbsPath(cfgPath string) string {
	if cfgPath == "cfg.json" {
		return filepath.Join(getWorkingDirAbsPath(), cfgPath)
	}

	wd, _ := os.Getwd()
	cfgAbsPath := filepath.Join(wd, cfgPath)
	return cfgAbsPath
}

func PublicIP() (string, error) {
	output, err := exec.Command("dig", "+short", "myip.opendns.com", "@resolver1.opendns.com").Output()
	if err != nil {
		return "UNKNOWN", err
	}
	ipStr := strings.TrimSpace(string(output))
	return ipStr, nil
}

func getJSONConfig() *JSONConfigFile {
	jsonCfgLock.RLock()
	defer jsonCfgLock.RUnlock()
	return jsonConfig
}

func getHostname() string {
	hostname := getJSONConfig().Hostname
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

func getIP() string {
	ip := getJSONConfig().IPAddress
	if ip != "" {
		log.Println("IP set in config: [", ip, "]")
		return ip
	}

	ip, err := PublicIP()
	if err != nil {
		log.Println("IP not set in config, getting public IP...failed:", err)
	} else {
		log.Println("IP not set in config, getting public IP...succeeded: [", ip, "]")
	}
	return ip
}

func getConnectionID() string {
	connectionID := getJSONConfig().ConnectionID
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

func loadJSONConfig(cfgFile string) {
	cfgFile = filepath.Clean(cfgFile)
	cfgPath := getCfgAbsPath(cfgFile)

	if !file.IsExist(cfgPath) {
		log.Fatalln("Configuration file [", cfgFile, "] doesn't exist")
	}

	configContent, err := file.ToTrimString(cfgPath)
	if err != nil {
		log.Fatalln("Reading configuration file [", cfgFile, "] failed:", err)
	}

	var c JSONConfigFile
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("Parsing configuration file [", cfgFile, "] failed:", err)
	}

	jsonCfgLock.Lock()
	defer jsonCfgLock.Unlock()

	jsonConfig = &c

	log.Println("Reading configuration file [", cfgFile, "] succeeded")
}

func GetGeneralConfig() *GeneralConfig {
	return generalConfig
}

func InitGeneralConfig(cfgFilePath string) {
	var cfg GeneralConfig
	generalConfig = &cfg

	loadJSONConfig(cfgFilePath)
	cfg.Agent = getJSONConfig().Agent
	cfg.Hbs = getJSONConfig().Hbs
	cfg.hbsResp.Store(model.NqmTaskResponse{})
	cfg.Hostname = getHostname()
	if cfg.IPAddress = getIP(); cfg.IPAddress == "UNKNOWN" {
		log.Fatalln("IP can't be \"UNKNOWN\"")
	}
	cfg.ConnectionID = getConnectionID()
}
