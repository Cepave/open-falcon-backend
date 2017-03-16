package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Cepave/open-falcon-backend/modules/api/app/controller"
	"github.com/Cepave/open-falcon-backend/modules/api/config"
	"github.com/Cepave/open-falcon-backend/modules/api/graph"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/gin-gonic/gin.v1"
	yaagGin "gopkg.in/masato25/yaag.v1/gin"
	"gopkg.in/masato25/yaag.v1/yaag"
)

func initGraph() {
	graph.Start(viper.GetStringMapString("graphs.cluster"))
}

func main() {
	cfgTmp := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	flag.Parse()
	cfg := *cfgTmp
	if *version {
		fmt.Println(config.VERSION)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	viper.AddConfigPath(".")
	viper.AddConfigPath("/")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./api/config")
	cfg = strings.Replace(cfg, ".json", "", 1)
	viper.SetConfigName(cfg)

	viper.ReadInConfig()
	err := config.InitLog(viper.GetString("log_level"))
	if err != nil {
		log.Fatal(err)
	}
	err = config.InitDB(viper.GetBool("db.db_bug"))
	if err != nil {
		log.Fatalf("db conn failed with error %s", err.Error())
	}
	routes := gin.Default()
	if viper.GetBool("gen_doc") {
		yaag.Init(&yaag.Config{
			On:       true,
			DocTitle: "Gin",
			DocPath:  viper.GetString("gen_doc_path"),
			BaseUrls: map[string]string{"Production": "/api/v1", "Staging": "/api/v1"},
		})
		routes.Use(yaagGin.Document())
	}
	initGraph()
	//start gin server
	controller.StartGin(viper.GetString("web_port"), routes)
}
