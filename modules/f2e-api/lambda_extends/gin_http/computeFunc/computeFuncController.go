package computeFunc

import (
	"strings"

	"github.com/Cepave/open-falcon-backend/modules/f2e-api/lambda_extends/conf"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/lambda_extends/gin_http/openFalcon"
	"github.com/gin-gonic/gin"
	_ "github.com/robertkrimen/otto/underscore"
)

func GetTestData(c *gin.Context) {
	c.JSON(200, gin.H{
		"data": getFakeData(),
	})
}

func GetAvaibleFun(c *gin.Context) {
	c.JSON(200, gin.H{
		"funcations": conf.GetAvaibleFun(),
	})
}

func getParamsFromHTTP(funcParams []string, c *gin.Context) (tmpparams map[string]interface{}) {
	tmpparams = map[string]interface{}{}
	for _, params := range funcParams {
		ss := strings.Split(params, ":")
		paramsKey := ss[0]
		paramset := c.DefaultQuery(paramsKey, "")
		if paramset != "" {
			tmpparams[paramsKey] = paramset
		}
	}
	return
}

func Compute(c *gin.Context) {
	funcName := c.DefaultQuery("funcName", "")
	if funcName == "" {
		c.JSON(400, gin.H{
			"msg": "Get params fun error",
		})
	}
	funcInstance := GetFuncSetup(funcName)
	if funcInstance.FuncationName == "" {
		c.JSON(400, gin.H{
			"msg": "Not found this compute method",
		})
	}
	vm := initJSvM()
	tmpparams := getParamsFromHTTP(funcInstance.Params, c)
	source := c.DefaultQuery("source", "real")
	if source == "real" {
		vm.Set("input", openFalcon.QDataGet(c))
	} else {
		vm.Set("input", getFakeData())
	}
	vm = SetParamsToJSVM(tmpparams, funcInstance.Params, vm)
	vm.Run(funcInstance.Codes)
	output, err := vm.Get("output")
	if err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"compted_data": output.String(),
		"funcName":     funcName,
		"paramsGot":    tmpparams,
	})
}
