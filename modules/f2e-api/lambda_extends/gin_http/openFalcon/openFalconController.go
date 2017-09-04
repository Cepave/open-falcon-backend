package openFalcon

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"strconv"

	"fmt"

	cmodel "github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/graph"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/lambda_extends/model"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func GetEndpoints(c *gin.Context) {
	enps := model.EndpointQuery("")
	c.JSON(200, gin.H{
		"status": "ok",
		"data": map[string][]string{
			"endpoints": enps,
		},
	})
}

func BatchQuery(startTs int64, endTs int64, consolFun string, step int, counter string, endpoints []string) (result []cmodel.GraphQueryResponse) {
	if viper.GetBool("test_mode") {
		return getFakeData()
	}
	result = []cmodel.GraphQueryResponse{}
	inputs := []cmodel.GraphQueryParam{}
	for _, enp := range endpoints {
		q := cmodel.GraphQueryParam{
			Start:     startTs,
			End:       endTs,
			ConsolFun: consolFun,
			Step:      step,
			Endpoint:  enp,
			Counter:   counter,
		}
		inputs = append(inputs, q)
	}
	result, errors := graph.QueryBatch(inputs)
	if len(errors) != 0 {
		for _, err := range errors {
			log.Errorf("lambda query got error :%v", err.Error())
		}
	}
	return
}

func QueryOnce(startTs int64, endTs int64, consolFun string, step int, counter string, endpoints []string) (result []cmodel.GraphQueryResponse) {
	if viper.GetBool("test_mode") {
		return getFakeData()
	}
	result = []cmodel.GraphQueryResponse{}
	for _, enp := range endpoints {
		q := cmodel.GraphQueryParam{
			Start:     startTs,
			End:       endTs,
			ConsolFun: consolFun,
			Step:      step,
			Endpoint:  enp,
			Counter:   counter,
		}
		res, _ := graph.QueryOne(q)
		log.Debug(fmt.Sprintf("%v, %v, %v", res.Counter, res.Endpoint, len(res.Values)))
		if res != nil {
			result = append(result, *res)
		} else {
			log.Errorf("graph.QueryOne return nil response [counter: %s, endpoint: %s]", res.Counter, res.Endpoint)
		}
	}
	return
}

func QDataGet(c *gin.Context) []cmodel.GraphQueryResponse {
	startTmp := c.DefaultQuery("startTs", string(time.Now().Unix()-(86400)))
	startTmp2, _ := strconv.Atoi(startTmp)
	startTs := int64(startTmp2)
	endTmp := c.DefaultQuery("endTs", string(time.Now().Unix()))
	endTmp2, _ := strconv.Atoi(endTmp)
	endTs := int64(endTmp2)
	consolFun := c.DefaultQuery("consolFun", "AVERAGE")
	stepTmp := c.DefaultQuery("step", "60")
	step, _ := strconv.Atoi(stepTmp)
	counter := c.DefaultQuery("counter", "cpu.idle")
	bossQ := c.DefaultQuery("platform", "ALL")
	var endpoints []string
	if bossQ == "ALL" {
		endpoints = model.EndpointQuery("")
	} else {
		endpoints = model.BossEndpointQuery(bossQ)
	}
	var result []cmodel.GraphQueryResponse
	// set single query rpc as default query, because graph not yet update batchquery rpc version into production
	batchQ := c.DefaultQuery("batch", "0")
	if batchQ == "0" {
		result = QueryOnce(startTs, endTs, consolFun, step, counter, endpoints)
	} else {
		result = BatchQuery(startTs, endTs, consolFun, step, counter, endpoints)
	}
	log.Debug(fmt.Sprintf("%s: %d", "openfaclon query got", len(result)))
	return result
}

func QueryData(c *gin.Context) {
	result := QDataGet(c)
	c.JSON(200, gin.H{
		"status": "ok",
		"data":   result,
	})
}

func getFakeData() (result []cmodel.GraphQueryResponse) {
	root_dir := viper.GetString("lambda_extends.root_dir")
	sampleDataPath := root_dir + "/data/test_data_sample1.json"
	data, err := ioutil.ReadFile(sampleDataPath)
	if err != nil {
		log.Error(err.Error())
		return
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Error(err.Error())
		return
	}
	return result
}
