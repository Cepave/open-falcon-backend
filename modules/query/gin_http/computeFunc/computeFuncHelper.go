package computeFunc

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Cepave/open-falcon-backend/modules/query/conf"
	"github.com/Cepave/open-falcon-backend/modules/query/model"
	"github.com/robertkrimen/otto"
)

func getFakeData() (t []*model.Result) {
	fakedataf, err := ioutil.ReadFile("./test/realdata")
	if err != nil {
		log.Error(err.Error())
	}
	var jdata string = string(fakedataf)
	json.Unmarshal([]byte(jdata), &t)
	return
}

func GetFuncSetup(funName string) *conf.FunConfig {
	return conf.GetFunc(funName)
}

func initJSvM() *otto.Otto {
	return otto.New()
}

func SetOttoVM(vm *otto.Otto, pmap map[string]string, key string, ptype string) {
	if value, ok := pmap[key]; ok {
		switch ptype {
		case "string":
			vm.Set(key, value)
		case "int":
			intval, err := strconv.Atoi(value)
			if err != nil {
				log.Error(err.Error())
			} else {
				vm.Set(key, intval)
			}
		case "bool":
			boolVal, err := strconv.ParseBool(value)
			if err != nil {
				log.Error(err.Error())
			}
			vm.Set(key, boolVal)
		}
	}
}

func SetParamsToJSVM(httpParams map[string]string, funcParams []string, vm *otto.Otto) *otto.Otto {
	for _, params := range funcParams {
		ss := strings.Split(params, ":")
		paramsKey := ss[0]
		paramsType := ss[1]
		if httpParams[paramsKey] != "" {
			SetOttoVM(vm, httpParams, paramsKey, paramsType)
		}
	}
	return vm
}
