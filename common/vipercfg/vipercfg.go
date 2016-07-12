package vipercfg

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"

	"github.com/chyeh/viper"
	"github.com/spf13/pflag"
)

func rmFileExt(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func Load() {
	viper.BindPFlag("config", pflag.Lookup("c"))
	viper.BindPFlag("version", pflag.Lookup("version"))

	cfgPath := viper.GetString("config")
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalln("Configuration file [", cfgPath, "] doesn't exist")
	}

	viper.SetConfigName(rmFileExt(cfgPath))

	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("Error:", err)
	}

}
