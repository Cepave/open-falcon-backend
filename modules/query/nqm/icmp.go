package nqm

import (
	"net/http"
	"encoding/json"
	"github.com/Cepave/open-falcon-backend/modules/query/g"
	osling "github.com/Cepave/open-falcon-backend/common/sling"
	"github.com/bitly/go-simplejson"
	"github.com/dghubble/sling"
)

// As the DSL arguments for calling JSONRPC
type IcmpDslArgs struct {
	Dsl *NqmDsl `json:"dsl"`
}

// The raw result returned from JSONRPC
type IcmpResult struct {
	grouping []int32
	metrics  *Metrics
}

// Used to unmarshal JSON with specific structure of IcmpResult(because of reusing struct)
func (icmpResult *IcmpResult) UnmarshalJSON(p []byte) error {
	jsonObj, err := simplejson.NewJson(p)
	if err != nil {
		return err
	}

	icmpResult.grouping = toGroupingIds(jsonObj.Get("grouping").MustArray())
	icmpResult.metrics = &Metrics{
		Max:                     int16(jsonObj.Get("max").MustInt()),
		Min:                     int16(jsonObj.Get("min").MustInt()),
		Avg:                     jsonObj.Get("avg").MustFloat64(),
		Med:                     int16(jsonObj.Get("med").MustInt()),
		Mdev:                    jsonObj.Get("mdev").MustFloat64(),
		Loss:                    jsonObj.Get("loss").MustFloat64(),
		Count:                   int32(jsonObj.Get("count").MustInt()),
		NumberOfSentPackets:     uint64(jsonObj.Get("number_of_sent_packets").MustUint64()),
		NumberOfReceivedPackets: uint64(jsonObj.Get("number_of_received_packets").MustUint64()),
		NumberOfAgents:          int32(jsonObj.Get("number_of_agents").MustUint64()),
		NumberOfTargets:         int32(jsonObj.Get("number_of_targets").MustUint64()),
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
