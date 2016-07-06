package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"sync"

	"github.com/Cepave/open-falcon-backend/modules/query/utils"
)

type Gconfig struct {
	Funcations []FunConfig
}

type FunConfig struct {
	FuncationName string   `json:"funcation_name"`
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
	f := "./"
	currentPath, _ := os.Getwd()
	if _, err := os.Stat(fmt.Sprintf("%s/js", currentPath)); err == nil {
		f = fmt.Sprintf("%s/js", currentPath)
	} else if _, err := os.Stat(fmt.Sprintf("%s/../js", currentPath)); err == nil {
		f = fmt.Sprintf("%s/../js", currentPath)
	} else if _, err := os.Stat(fmt.Sprintf("%s/../conf/js", currentPath)); err == nil {
		f = fmt.Sprintf("%s/../conf/js", currentPath)
	} else {
		log.Fatalf("load js files got error, currentPaht: %s , please check your code tree and make is correct!", currentPath)
	}
	FunctionMap = map[string]*FunConfig{}
	for _, v := range gconfig {
		contain := jsFileReader(fmt.Sprintf("%s/%s", f, v.FilePath))
		v.Codes = contain
		FunctionMap[v.FuncationName] = v
	}
}

func ReadConf() {
	var f string
	currentPath, _ := os.Getwd()
	if _, err := os.Stat(fmt.Sprintf("%s/conf/lambdaSetup.json", currentPath)); err == nil {
		f = fmt.Sprintf("%s/conf/lambdaSetup.json", currentPath)
	} else if _, err := os.Stat(fmt.Sprintf("%s/../conf/lambdaSetup.json", currentPath)); err == nil {
		f = fmt.Sprintf("%s/../conf/lambdaSetup.json", currentPath)
	} else {
		log.Fatalf("lambdaStup.json not found, currentPaht: %s", currentPath)
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
