package nqm

import (
	"encoding/json"
	"net/http"
	"github.com/Cepave/open-falcon-backend/modules/query/g"
	osling "github.com/Cepave/open-falcon-backend/common/sling"
	sjson "github.com/bitly/go-simplejson"
	"github.com/dghubble/sling"
	model "github.com/Cepave/open-falcon-backend/modules/query/model/nqm"
	"github.com/Cepave/open-falcon-backend/common/utils"
)

var rootClient *sling.Sling

// The raw result returned from JSONRPC
type IcmpResult struct {
	grouping []int32
	metrics  *model.Metrics
}

// Used to unmarshal JSON with specific structure of IcmpResult(because of reusing struct)
func (r *IcmpResult) UnmarshalJSON(p []byte) error {
	jsonObj, err := sjson.NewJson(p)
	if err != nil {
		return err
	}

	/**
	 * Unmarshal grouping on ids
	 */
	r.grouping = utils.MakeAbstractArray(jsonObj.Get("grouping").MustArray()).
		MapTo(
			func(v interface{}) interface{} {
				intV, err := v.(json.Number).Int64()
				if err != nil {
					panic(err)
				}
				return int32(intV)
			},
			utils.TypeOfInt32,
		).
		GetArray().([]int32)
	// :~)

	/**
	 * Unmarshal metric data
	 */
	newMetrics := &model.Metrics{}
	newMetrics.UnmarshalSimpleJson(jsonObj)
	r.metrics = newMetrics
	// :~)

	return nil
}

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
