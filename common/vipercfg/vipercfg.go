package vipercfg

import (
	log "github.com/Sirupsen/logrus"
	"github.com/chyeh/viper"
	"github.com/spf13/pflag"
)

var v = viper.New()

func Load() {
	cfgPath := v.GetString("config")
	v.SetConfigFile(cfgPath)
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalln("Error:", err)
	}
}

func Parse() {
	pflag.StringP("config", "c", "cfg.json", "configuration file")
	pflag.BoolP("version", "v", false, "show version")
	pflag.Bool("check", false, "check collector")
	pflag.BoolP("help", "h", false, "usage")
	pflag.Bool("vg", false, "show version and git commit log")
	pflag.Parse()
}

func Bind() {
	v.BindPFlag("config", pflag.Lookup("config"))
	v.BindPFlag("version", pflag.Lookup("version"))
	v.BindPFlag("check", pflag.Lookup("check"))
	v.BindPFlag("help", pflag.Lookup("help"))
	v.BindPFlag("vg", pflag.Lookup("vg"))
}

func Config() *viper.Viper {
	return v
}
