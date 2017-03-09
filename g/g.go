package g

import "path/filepath"

var Modules map[string]bool
var BinOf map[string]string
var cfgOf map[string]string
var ModuleApps map[string]string
var logpathOf map[string]string
var PidOf map[string]string
var AllModulesInOrder []string

func init() {
	//	dirs, _ := ioutil.ReadDir("./modules")
	//	for _, dir := range dirs {
	//		Modules[dir.Name()] = true
	//	}
	Modules = map[string]bool{
		"agent":      true,
		"aggregator": true,
		"alarm":      true,
		"fe":         true,
		"graph":      true,
		"hbs":        true,
		"judge":      true,
		"nodata":     true,
		"nqm-mng":    true,
		"query":      true,
		"sender":     true,
		"task":       true,
		"transfer":   true,
		"api":        true,
	}

	BinOf = map[string]string{
		"agent":      "./agent/bin/falcon-agent",
		"aggregator": "./aggregator/bin/falcon-aggregator",
		"alarm":      "./alarm/bin/falcon-alarm",
		"fe":         "./fe/bin/falcon-fe",
		"graph":      "./graph/bin/falcon-graph",
		"hbs":        "./hbs/bin/falcon-hbs",
		"judge":      "./judge/bin/falcon-judge",
		"nodata":     "./nodata/bin/falcon-nodata",
		"nqm-mng":    "./nqm-mng/bin/falcon-nqm-mng",
		"query":      "./query/bin/falcon-query",
		"sender":     "./sender/bin/falcon-sender",
		"task":       "./task/bin/falcon-task",
		"transfer":   "./transfer/bin/falcon-transfer",
		"api":        "./api/bin/falcon-api",
	}

	cfgOf = map[string]string{
		"agent":      "./agent/config/cfg.json",
		"aggregator": "./aggregator/config/cfg.json",
		"alarm":      "./alarm/config/cfg.json",
		"fe":         "./fe/config/cfg.json",
		"graph":      "./graph/config/cfg.json",
		"hbs":        "./hbs/config/cfg.json",
		"judge":      "./judge/config/cfg.json",
		"nodata":     "./nodata/config/cfg.json",
		"nqm-mng":    "./nqm-mng/config/cfg.json",
		"query":      "./query/config/cfg.json",
		"sender":     "./sender/config/cfg.json",
		"task":       "./task/config/cfg.json",
		"transfer":   "./transfer/config/cfg.json",
		"api":        "./api/config/cfg.json",
	}

	ModuleApps = map[string]string{
		"agent":      "falcon-agent",
		"aggregator": "falcon-aggregator",
		"alarm":      "falcon-alarm",
		"graph":      "falcon-graph",
		"fe":         "falcon-fe",
		"hbs":        "falcon-hbs",
		"judge":      "falcon-judge",
		"nodata":     "falcon-nodata",
		"nqm-mng":    "falcon-nqm-mng",
		"query":      "falcon-query",
		"sender":     "falcon-sender",
		"task":       "falcon-task",
		"transfer":   "falcon-transfer",
		"api":        "falcon-api",
	}

	logpathOf = map[string]string{
		"agent":      "./agent/logs/agent.log",
		"aggregator": "./aggregator/logs/aggregator.log",
		"alarm":      "./alarm/logs/alarm.log",
		"fe":         "./fe/logs/fe.log",
		"graph":      "./graph/logs/graph.log",
		"hbs":        "./hbs/logs/hbs.log",
		"judge":      "./judge/logs/judge.log",
		"nodata":     "./nodata/logs/nodata.log",
		"nqm-mng":    "./nqm-mng/logs/nqm-mng.log",
		"query":      "./query/logs/query.log",
		"sender":     "./sender/logs/sender.log",
		"task":       "./task/logs/task.log",
		"transfer":   "./transfer/logs/transfer.log",
		"api":        "./api/logs/api.log",
	}

	PidOf = map[string]string{
		"agent":      "<NOT SET>",
		"aggregator": "<NOT SET>",
		"alarm":      "<NOT SET>",
		"graph":      "<NOT SET>",
		"fe":         "<NOT SET>",
		"hbs":        "<NOT SET>",
		"judge":      "<NOT SET>",
		"nodata":     "<NOT SET>",
		"nqm-mng":    "<NOT SET>",
		"query":      "<NOT SET>",
		"sender":     "<NOT SET>",
		"task":       "<NOT SET>",
		"transfer":   "<NOT SET>",
		"api":        "<NOT SET>",
	}

	// Modules are deployed in this order
	AllModulesInOrder = []string{
		"graph",
		"hbs",
		"fe",
		"alarm",
		"sender",
		"query",
		"judge",
		"transfer",
		"nodata",
		"nqm-mng",
		"task",
		"aggregator",
		"agent",
		"api",
	}
}

func Bin(name string) string {
	p, _ := filepath.Abs(BinOf[name])
	return p
}

func Cfg(name string) string {
	p, _ := filepath.Abs(cfgOf[name])
	return p
}

func LogPath(name string) string {
	p, _ := filepath.Abs(logpathOf[name])
	return p
}

func LogDir(name string) string {
	d, _ := filepath.Abs(filepath.Dir(logpathOf[name]))
	return d
}
