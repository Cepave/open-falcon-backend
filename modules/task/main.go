package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/viper"

	oHttp "github.com/Cepave/open-falcon-backend/common/http"
	client "github.com/Cepave/open-falcon-backend/common/http/client"
	"github.com/Cepave/open-falcon-backend/common/logruslog"
	oos "github.com/Cepave/open-falcon-backend/common/os"
	"github.com/Cepave/open-falcon-backend/common/vipercfg"

	"github.com/Cepave/open-falcon-backend/modules/task/collector"
	"github.com/Cepave/open-falcon-backend/modules/task/cron"
	"github.com/Cepave/open-falcon-backend/modules/task/database"
	"github.com/Cepave/open-falcon-backend/modules/task/g"
	"github.com/Cepave/open-falcon-backend/modules/task/http"
	"github.com/Cepave/open-falcon-backend/modules/task/index"
	"github.com/Cepave/open-falcon-backend/modules/task/proc"
)

func main() {
	vipercfg.Parse()
	vipercfg.Bind()

	if vipercfg.Config().GetBool("version") {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}
	if vipercfg.Config().GetBool("vg") {
		fmt.Println(g.VERSION, g.COMMIT)
		os.Exit(0)
	}

	// global config
	vipercfg.Load()
	g.ParseConfig(vipercfg.Config().GetString("config"))
	logruslog.Init()

	viperObj := vipercfg.Config()

	/**
	 * Variables of services
	 */
	var cronService *cron.TaskCronService
	// :~)

	// proc
	proc.Start()

	// graph index
	index.Start()
	// collector
	collector.Start()

	// http
	http.Start()

	/**
	 * Initializes APIs to databases
	 */
	database.InitMySqlApi(buildRestfulClientConfig(viperObj))
	// :~)

	/**
	 * Initializes cron services from Viper configuration and starts it
	 */
	cronService = cron.NewCronServices(buildTaskCronConfig(viperObj))
	cronService.Start()
	// :~)

	oos.HoldingAndWaitSignal(
		func(signal os.Signal) {
			cronService.Stop()
		},
		os.Interrupt, os.Kill,
		syscall.SIGTERM,
	)
}

func buildRestfulClientConfig(viperObj *viper.Viper) *oHttp.RestfulClientConfig {
	url := viperObj.GetString("mysql_api.host")

	if resource := viperObj.GetString("mysql_api.resource"); resource != "" {
		url += "/" + resource
	}

	httpClientConfig := client.NewDefaultConfig()
	httpClientConfig.Url = url

	return &oHttp.RestfulClientConfig{
		HttpClientConfig: httpClientConfig,
		FromModule:       "Task",
	}
}
func buildTaskCronConfig(viperObj *viper.Viper) *cron.TaskCronConfig {
	return &cron.TaskCronConfig{
		VacuumQueryObjects: &cron.VacuumQueryObjectsConf{
			Cron:    viperObj.GetString("cron.vacuum_query_objects.schedule"),
			ForDays: viperObj.GetInt("cron.vacuum_query_objects.for_days"),
			Enable:  viperObj.GetBool("cron.vacuum_query_objects.enable"),
		},
		VacuumGraphIndex: &cron.VacuumGraphIndexConf{
			Cron:    viperObj.GetString("cron.vacuum_graph_index.schedule"),
			ForDays: viperObj.GetInt("cron.vacuum_graph_index.for_days"),
			Enable:  viperObj.GetBool("cron.vacuum_graph_index.enable"),
		},
		ClearTaskLogEntries: &cron.ClearTaskLogEntriesConf{
			Cron:    viperObj.GetString("cron.clear_task_log_entries.schedule"),
			ForDays: viperObj.GetInt("cron.clear_task_log_entries.for_days"),
			Enable:  viperObj.GetBool("cron.clear_task_log_entries.enable"),
		},
	}
}
