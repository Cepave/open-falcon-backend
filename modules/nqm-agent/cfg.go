package main

import (
	"fmt"
	"net"
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

type JSONConfig struct {
	Agent        *AgentConfig `json:"agent"`
	Hbs          *HbsConfig   `json:"hbs"`
	Hostname     string       `json:"hostname"`
	IPAddress    string       `json:"ipAddress"`
	ConnectionID string       `json:"connectionID"`
}

type Metadata struct {
	Hostname     string
	IPAddress    string
	ConnectionID string
}

var (
	jsonConfig atomic.Value // for JSONConfig
	hbsResp    atomic.Value // for receiving model.NqmTaskResponse
	metadata   = Metadata{}
)

func digPublicIP() (net.IP, error) {
	output, err := exec.Command("dig", "+short", "myip.opendns.com", "@resolver1.opendns.com").Output()
	if err != nil {
		return nil, err
	}

	ipStr := strings.TrimSpace(string(output))
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("Cannot parse IP Address from dig's output: [%s]", output)
	}

	return ip, nil
}

func Config() JSONConfig {
	return jsonConfig.Load().(JSONConfig)
}

func SetConfig(c JSONConfig) {
	jsonConfig.Store(c)
}

func HBSResp() model.NqmTaskResponse {
	return hbsResp.Load().(model.NqmTaskResponse)
}

func SetHBSResp(r model.NqmTaskResponse) {
	hbsResp.Store(r)
}

func hostname() string {
	hostname := Config().Hostname
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
	ip := net.ParseIP(Config().IPAddress)
	if ip != nil {
		log.Println("IP set in config: [", ip, "]")
		return ip.String()
	}
	log.Errorln("Invalid IP in config")

	ip, err := digPublicIP()
	if err != nil {
		log.Fatalln("Getting public IP...failed:", err)
	}
	log.Println("Getting public IP...succeeded: [", ip, "]")

	return ip.String()
}

func connectionID() string {
	connectionID := Config().ConnectionID
	if connectionID != "" {
		log.Println("ConnectionID set in config: [", connectionID, "]")
		return connectionID
	}

	// Logically it shouldn't happen because ConnectionID is alwasy generated
	// after Hostname and IPAddress are set.
	if Meta().Hostname == "" || Meta().IPAddress == "" {
		log.Fatalln("ConnectionID not set in config, generating...failed!")
	}

	connectionID = Meta().Hostname + "@" + Meta().IPAddress
	log.Println("ConnectionID not set in config, generating...succeeded: [", connectionID, "]")
	return connectionID
}

func jsonUnmarshaller() JSONConfig {
	var c = JSONConfig{
		Agent: &AgentConfig{},
		Hbs:   &HbsConfig{},
	}
	err := vipercfg.Config().Unmarshal(&c)
	if err != nil {
		log.Fatal("Parsing configuration file [", vipercfg.Config().GetString("config"), "] failed:", err)
	}
	log.Println("Reading configuration file [", vipercfg.Config().GetString("config"), "] succeeded")
	return c
}

func InitConfig() {
	c := jsonUnmarshaller()
	SetConfig(c)
}

func Meta() *Metadata {
	return &metadata
}

func GenMeta() {
	Meta().Hostname = hostname()
	Meta().IPAddress = ip()
	Meta().ConnectionID = connectionID()
}
