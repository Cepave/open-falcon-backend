package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/toolkits/file"
)

type AgentConfig struct {
	PushURL         string `json:"pushURL"`
	FpingInterval   int    `json:"fpingInterval"`
	TcppingInterval int    `json:"tcppingInterval"`
	TcpconnInterval int    `json:"tcpconnInterval"`
}

type HbsConfig struct {
	RPCServer string `json:"RPCServer"`
	Interval  int    `json:"interval"`
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
	version := flag.Bool("v", false, "show version")
	flag.Parse()
	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}
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

func InitGeneralConfig() {
	var cfg GeneralConfig
	generalConfig = &cfg
	loadJSONConfig()
	cfg.Agent = getJSONConfig().Agent
	cfg.Hbs = getJSONConfig().Hbs
	cfg.Hostname = getHostname()
	cfg.IPAddress = getIP()
	cfg.ConnectionID = getConnectionID()
}
