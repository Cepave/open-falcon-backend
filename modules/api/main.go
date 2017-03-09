package main

import (
	"github.com/Cepave/open-falcon-backend/modules/api/app/controller"
	"github.com/Cepave/open-falcon-backend/modules/api/config"
	"github.com/Cepave/open-falcon-backend/modules/api/graph"
	log "github.com/Sirupsen/logrus"
	yaag_gin "github.com/masato25/yaag/gin"
	"github.com/spf13/viper"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/masato25/yaag.v1/yaag"
)

func initGraph() {
	graph.Start(viper.GetStringMapString("graphs.cluster"))
}

func main() {
	viper.AddConfigPath(".")
	viper.SetConfigName("cfg")
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
		routes.Use(yaag_gin.Document())
	}
	initGraph()
	//start gin server
	controller.StartGin(viper.GetString("web_port"), routes)
}
