package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Cepave/common/model"
	"github.com/toolkits/file"
)

type JSONConfigFile struct {
	AgentPushURL string `json:"agentPushURL"`
	HbsRPCServer string `json:"hbsRPCServer"`
	Hostname     string `json:"hostname"`
	IPAddress    string `json:"ipAddress"`
	ConnectionID string `json:"connectionID"`
}

type GeneralConfig struct {
	JSONConfigFile
	Hostname     string
	IPAddress    string
	ConnectionID string
	ISP          string
	Province     string
	City         string
}

var (
	jsonConfig    *JSONConfigFile
	generalConfig *GeneralConfig
	lock          = new(sync.RWMutex)
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
	lock.RLock()
	defer lock.RUnlock()
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
	}

	log.Println("IP not set in config, getting public IP...succeeded: [", ip, "]")
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

func loadJSONConfig() {
	parsedPtr := flag.String("c", "cfg.json", "nqm's configuration file")
	flag.Parse()
	cfgFile := filepath.Clean(*parsedPtr)
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

	lock.Lock()
	defer lock.Unlock()

	jsonConfig = &c

	log.Println("Reading configuration file [", cfgFile, "] succeeded")
}

func GetGeneralConfig() *GeneralConfig {
	return generalConfig
}

func SetGeneralConfigByAgent(agent model.NqmAgent) {
	GetGeneralConfig().ISP = agent.IspName
	GetGeneralConfig().Province = agent.ProvinceName
	GetGeneralConfig().City = agent.CityName
}

func InitGeneralConfig() {
	var cfg GeneralConfig
	generalConfig = &cfg
	loadJSONConfig()
	cfg.AgentPushURL = getJSONConfig().AgentPushURL
	cfg.HbsRPCServer = getJSONConfig().HbsRPCServer
	cfg.Hostname = getHostname()
	cfg.IPAddress = getIP()
	cfg.ConnectionID = getConnectionID()
}
