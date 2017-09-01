package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/controller"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/graph"
	jconf "github.com/Cepave/open-falcon-backend/modules/f2e-api/lambda_extends/conf"
	"github.com/gin-gonic/gin"
	yaagGin "github.com/masato25/yaag/gin"
	"github.com/masato25/yaag/yaag"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Version   = "<UNDEFINED>"
	GitCommit = "<UNDEFINED>"
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
		fmt.Printf("version %s, build %s\n", Version, GitCommit)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	var err error
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	viper.Set("lambda_extends.root_dir", pwd)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./api/config")
	cfg = strings.Replace(cfg, ".json", "", 1)
	viper.SetConfigName(cfg)

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = config.InitLog(viper.GetString("log_level"))
	if err != nil {
		log.Fatal(err)
	}
	err = config.InitDB(viper.GetBool("db.db_debug"))
	if err != nil {
		log.Fatalf("db conn failed with error %s", err.Error())
	}
	config.Init()
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
	// init graph rpc pools
	initGraph()

	// init redis connection for external alarm Use
	config.StartRedis(
		viper.GetBool("redis.enable"),
		viper.GetString("redis.address"),
		// if not passwrod set, keep blank
		viper.GetString("redis.password"),
		viper.GetString("redis.default_bucket"),
	)

	//inject lambda web
	if viper.GetBool("lambda_extends.enable") {
		jconf.ReadConf()
	}

	//start gin server
	log.Debugf("will start with port: %v, test_mode: %v", viper.GetString("web_port"), viper.GetBool("test_mode"))
	if viper.GetBool("test_mode") {
		go func() {
			ginTest := controller.StartGin(viper.GetString("web_port"), routes, viper.GetBool("test_mode"))
			ginTest.Run(viper.GetString("web_port"))
		}()
	} else {
		go controller.StartGin(viper.GetString("web_port"), routes, viper.GetBool("test_mode"))
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		os.Exit(0)
	}()
	select {}
}
