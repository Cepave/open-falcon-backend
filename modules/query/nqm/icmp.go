package nqm

import (
	"net/http"
	"encoding/json"
	"github.com/Cepave/open-falcon-backend/modules/query/g"
	osling "github.com/Cepave/open-falcon-backend/common/sling"
	"github.com/bitly/go-simplejson"
	"github.com/dghubble/sling"
	ojson "github.com/Cepave/open-falcon-backend/common/json"
	model "github.com/Cepave/open-falcon-backend/modules/query/model/nqm"
)

// As the DSL arguments for calling JSONRPC
type IcmpDslArgs struct {
	Dsl *NqmDsl `json:"dsl"`
}

// The raw result returned from JSONRPC
type IcmpResult struct {
	grouping []int32
	metrics  *model.Metrics
}

// Used to unmarshal JSON with specific structure of IcmpResult(because of reusing struct)
func (icmpResult *IcmpResult) UnmarshalJSON(p []byte) error {
	sjson, err := simplejson.NewJson(p)
	if err != nil {
		return err
	}

	jsonObj := ojson.ToJsonExt(sjson)

	icmpResult.grouping = toGroupingIds(jsonObj.Get("grouping").MustArray())
	icmpResult.metrics = &model.Metrics{
		Max:                     jsonObj.GetExt("max").MustInt16(),
		Min:                     jsonObj.GetExt("min").MustInt16(),
		Avg:                     jsonObj.Get("avg").MustFloat64(),
		Med:                     jsonObj.GetExt("med").MustInt16(),
		Mdev:                    jsonObj.Get("mdev").MustFloat64(),
		Loss:                    jsonObj.Get("loss").MustFloat64(),
		Count:                   jsonObj.GetExt("count").MustInt32(),
		NumberOfSentPackets:     jsonObj.Get("number_of_sent_packets").MustUint64(),
		NumberOfReceivedPackets: jsonObj.Get("number_of_received_packets").MustUint64(),
		NumberOfAgents:          jsonObj.GetExt("number_of_agents").MustInt32(),
		NumberOfTargets:         jsonObj.GetExt("number_of_targets").MustInt32(),
	}

	return nil
}

func toGroupingIds(srcArray []interface{}) []int32 {
	result := make([]int32, len(srcArray))

	for i, v := range srcArray {
		number, e := v.(json.Number).Int64()
		if e != nil {
			panic(e)
		}
		result[i] = int32(number)
	}

	return result
}

var rootClient *sling.Sling

func initIcmp() {
	config := g.Config()

	rootClient = sling.New().Base(config.NqmLog.ServiceUrl)
}

// Retrieves the statistics of ICMP log by DSL
func getStatisticsOfIcmpByDsl(queryParams *NqmDsl) (result []IcmpResult, err error) {
	err = osling.ToSlintExt(
		rootClient.New().Post("/nqm/icmp/query/by-dsl").
			BodyJSON(queryParams),
	).DoReceive(http.StatusOK, &result)

	return
}
