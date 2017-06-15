package config

import (
	"log"
	"runtime"

	"github.com/spf13/viper"
)

type apiClients struct {
	KeyPairs map[string]interface{}
	Enable   bool
}

func (mine *apiClients) NameIncludes(name string) bool {
	keys := mine.Keys()
	for _, k := range keys {
		if k == name {
			return true
		}
	}
	return false
}

func (mine *apiClients) AuthToken(name string, sig string) bool {
	keys := mine.Keys()
	for _, k := range keys {
		if k == name && mine.KeyPairs[name] == sig {
			return true
		}
	}
	return false
}

func (mine *apiClients) Keys() []string {
	keys := make([]string, 0, len(mine.KeyPairs))
	for k := range mine.KeyPairs {
		keys = append(keys, k)
	}
	return keys
}

// change log:
// deprecated
const (
	VERSION = "0.0.1"
)

var ApiClient apiClients

func Init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	ApiClient.KeyPairs = viper.GetStringMap("services")
	ApiClient.Enable = viper.GetBool("enable_services")
}
