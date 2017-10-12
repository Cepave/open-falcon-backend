package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/Cepave/open-falcon-backend/modules/f2e-api/lambda_extends/utils"
	"github.com/spf13/viper"
)

type Gconfig struct {
	Funcations []FunConfig
}

type FunConfig struct {
	FuncationName string   `json:"function_name"`
	FilePath      string   `json:"file_path"`
	Params        []string `json:"params"`
	Description   string   `json:"description"`
	Codes         string   `json:"-"`
}

var (
	gconfig     []*FunConfig
	configLock  = new(sync.RWMutex)
	confpath    *string
	FunctionMap map[string]*FunConfig
)

func Config() map[string]*FunConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return FunctionMap
}

func functionMapGen() {
	currentPath := viper.GetString("lambda_extends.root_dir")
	possiblePath := []string{"lambda_extends/config/js", "lambda_extends/conf/js", "lambda_extends/js", "../config/js", "../js", "f2e-api/bin"}
	f := ""
	for _, pa := range possiblePath {
		paf := fmt.Sprintf("%s/%s", currentPath, pa)
		if _, err := os.Stat(paf); err != nil {
			log.Debugf("can't not load file from: %s", paf)
		} else {
			f = paf
			break
		}
	}
	if f == "" {
		log.Fatalf("load js files got error, currentPaht: %s , please check your code tree and make is correct!", currentPath)
	} else {
		log.Info("load javascript scrips successed in " + f)
	}

	FunctionMap = map[string]*FunConfig{}
	for _, v := range gconfig {
		contain := jsFileReader(fmt.Sprintf("%s/%s", f, v.FilePath))
		v.Codes = contain
		FunctionMap[v.FuncationName] = v
	}
}

func ReadConf() {
	currentPath := viper.GetString("lambda_extends.root_dir")
	possiblePath := []string{"lambda_extends/conf/lambdaSetup.json", "lambda_extends/config/lambdaSetup.json", "../config/lambdaSetup.json", "f2e-api/config/lambdaSetup.json"}
	f := ""
	for _, pa := range possiblePath {
		paf := fmt.Sprintf("%s/%s", currentPath, pa)
		if _, err := os.Stat(paf); err != nil {
			log.Debugf("can't not load file from: %s", paf)
		} else {
			f = paf
			break
		}
	}
	if f == "" {
		log.Fatalf("lambdaSetup.json not found, currentPaht: %s", currentPath)
	} else {
		log.Info("read lambdaSetup.json successed wuth " + f)
	}

	confpath = &f
	dat, err := ioutil.ReadFile(f)
	if err != nil {
		log.Println(err)
	}
	var myconf []*FunConfig
	json.Unmarshal(dat, &myconf)
	if len(myconf) == 0 {
		log.Println("conf file is empty or format is wrong, please check it!")
	}
	gconfig = myconf
	functionMapGen()
}

func Reload() {
	configLock.RLock()
	ReadConf()
	defer configLock.RUnlock()
}

func GetFunc(key string) *FunConfig {
	return FunctionMap[key]
}

func GetAvaibleFun() []string {
	return utils.GetMapKeys(reflect.ValueOf(FunctionMap).MapKeys())
}
