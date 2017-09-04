package computeFunc

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Cepave/open-falcon-backend/modules/f2e-api/lambda_extends/conf"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/lambda_extends/model"
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

func SetOttoVM(vm *otto.Otto, pmap map[string]interface{}, key string, ptype string) {
	if value, ok := pmap[key]; ok {
		switch value.(type) {
		case string:
			vm.Set(key, value.(string))
		case int, int32, int64:
			vm.Set(key, value.(int64))
		case float32, float64:
			vm.Set(key, value.(float64))
		case bool:
			vm.Set(key, value.(bool))
		default:
			log.Errorf("no support type for function params: %v, type: %v", value, reflect.TypeOf(value))
		}
	}
}

func SetParamsToJSVM(httpParams map[string]interface{}, funcParams []string, vm *otto.Otto) *otto.Otto {
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
