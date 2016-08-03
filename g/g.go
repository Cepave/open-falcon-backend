package g

//	"io/ioutil"

var Modules map[string]bool
var ModuleBins map[string]string
var ModuleConfs map[string]string
var ModuleApps map[string]string
var AllModulesInOrder []string

func init() {
	//	dirs, _ := ioutil.ReadDir("./modules")

	//	for _, dir := range dirs {
	//		Modules[dir.Name()] = true
	//	}
	Modules = map[string]bool{
		"agent":      true,
		"nqm-agent":  true,
		"aggregator": true,
		"alarm":      true,
		"fe":         true,
		"graph":      true,
		"hbs":        true,
		"judge":      true,
		"nodata":     true,
		"query":      true,
		"sender":     true,
		"task":       true,
		"transfer":   true,
	}

	ModuleBins = map[string]string{
		"agent":      "./agent/bin/falcon-agent",
		"nqm-agent":  "./agent/bin/falcon-nqm-agent",
		"aggregator": "./aggregator/bin/falcon-aggregator",
		"alarm":      "./alarm/bin/falcon-alarm",
		"fe":         "./fe/bin/falcon-fe",
		"graph":      "./graph/bin/falcon-graph",
		"hbs":        "./hbs/bin/falcon-hbs",
		"judge":      "./judge/bin/falcon-judge",
		"nodata":     "./nodata/bin/falcon-nodata",
		"query":      "./query/bin/falcon-query",
		"sender":     "./sender/bin/falcon-sender",
		"task":       "./task/bin/falcon-task",
		"transfer":   "./transfer/bin/falcon-transfer",
	}

	ModuleConfs = map[string]string{
		"agent":      "./agent/config/cfg.json",
		"nqm-agent":  "./nqm-agent/config/cfg.json",
		"aggregator": "./aggregator/config/cfg.json",
		"alarm":      "./alarm/config/cfg.json",
		"fe":         "./fe/config/cfg.json",
		"graph":      "./graph/config/cfg.json",
		"hbs":        "./hbs/config/cfg.json",
		"judge":      "./judge/config/cfg.json",
		"nodata":     "./nodata/config/cfg.json",
		"query":      "./query/config/cfg.json",
		"sender":     "./sender/config/cfg.json",
		"task":       "./task/config/cfg.json",
		"transfer":   "./transfer/config/cfg.json",
	}

	ModuleApps = map[string]string{
		"agent":      "falcon-agent",
		"nqm-agent":  "falcon-nqm-agent",
		"aggregator": "falcon-aggregator",
		"alarm":      "falcon-alarm",
		"graph":      "falcon-graph",
		"fe":         "falcon-fe",
		"hbs":        "falcon-hbs",
		"judge":      "falcon-judge",
		"nodata":     "falcon-nodata",
		"query":      "falcon-query",
		"sender":     "falcon-sender",
		"task":       "falcon-task",
		"transfer":   "falcon-transfer",
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
		"task",
		"aggregator",
		"agent",
	}
}
