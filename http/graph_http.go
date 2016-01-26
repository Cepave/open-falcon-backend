package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	cmodel "github.com/Cepave/common/model"
	"github.com/Cepave/query/graph"
	"github.com/Cepave/query/proc"
	"regexp"
)

type GraphHistoryParam struct {
	Start            int                     `json:"start"`
	End              int                     `json:"end"`
	CF               string                  `json:"cf"`
	EndpointCounters []cmodel.GraphInfoParam `json:"endpoint_counters"`
}

func graphQueryOne(ec cmodel.GraphInfoParam, body GraphHistoryParam, endpoint string, counter string) *cmodel.GraphQueryResponse {
	if endpoint == "" {
		endpoint = ec.Endpoint
	}
	if counter == "" {
		counter = ec.Counter
	}
	request := cmodel.GraphQueryParam{
		Start:     int64(body.Start),
		End:       int64(body.End),
		ConsolFun: body.CF,
		Endpoint:  endpoint,
		Counter:   counter,
	}
	result, err := graph.QueryOne(request)
	if err != nil {
		log.Printf("graph.queryOne fail, %v, (endpoint, counter) = (%s, %s)", err, endpoint, counter)
	}
	return result
}

func configGraphRoutes() {

	// method:post
	http.HandleFunc("/graph/history", func(w http.ResponseWriter, r *http.Request) {
		// statistics
		proc.HistoryRequestCnt.Incr()

		var body GraphHistoryParam
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&body)
		if err != nil {
			StdRender(w, "", err)
			return
		}

		if len(body.EndpointCounters) == 0 {
			StdRender(w, "", errors.New("empty_payload"))
			return
		}

		data := []*cmodel.GraphQueryResponse{}
		var result *cmodel.GraphQueryResponse
		for _, ec := range body.EndpointCounters {
			if strings.Contains(ec.Counter, "packets-sent") || strings.Contains(ec.Counter, "transmission-time") {
				// NQM Case
				var packetSentCount []cmodel.JsonFloat
				// clone a response
				result = graphQueryOne(ec, body, "", "")
				for i := range result.Values {
					result.Values[i].Value = cmodel.JsonFloat(0.0)
					packetSentCount = append(packetSentCount, cmodel.JsonFloat(0.0))
				}
				if result == nil {
					continue
				}

				resultSent := graphQueryOne(ec, body, "", "")
				if resultSent == nil {
					continue
				}
				data = append(data, resultSent)

				if strings.Contains(ec.Counter, "packets-sent") {
					counter := strings.Replace(ec.Counter, "packets-sent", "packets-received", 1)
					if resultReceived := graphQueryOne(ec, body, "", counter); resultReceived != nil {
						data = append(data, result)
						for i := range resultSent.Values {
							packetLossCount := (resultSent.Values[i].Value - resultReceived.Values[i].Value)
							result.Values[i].Value += packetLossCount
							packetSentCount[i] += resultSent.Values[i].Value
						}
					}
				} else if strings.Contains(ec.Counter, "transmission-time") {
					counter := strings.Replace(ec.Counter, "transmission-time", "packets-sent", 1)
					if resultTransmissionTime := graphQueryOne(ec, body, "", counter); resultTransmissionTime != nil {
						data = append(data, resultSent)
						for i := range resultTransmissionTime.Values {
							result.Values[i].Value += resultTransmissionTime.Values[i].Value * resultSent.Values[i].Value
							packetSentCount[i] += resultSent.Values[i].Value
						}
					}
				}

				if strings.Contains(ec.Counter, "packets-sent") {
					result.Endpoint = "all-endpoints"
					result.Counter = "packet-loss-rate"
				} else if strings.Contains(ec.Counter, "transmission-time") {
					result.Endpoint = "all-endpoints"
					result.Counter = "average"
				}
				for i := range result.Values {
					result.Values[i].Value = result.Values[i].Value / packetSentCount[i]
				}
				result.Values = result.Values
				data = append(data, result)

			} else {
				regx, _ := regexp.Compile("(\\.\\$\\s*|\\s*)$")
				endpoint := regx.ReplaceAllString(ec.Endpoint, "")
				counter := regx.ReplaceAllString(ec.Counter, "")
				result = graphQueryOne(ec, body, endpoint, counter)
				if result == nil {
					continue
				}
				data = append(data, result)
			}
		}
		log.Println("got length of data: ", len(data))

		// statistics
		proc.HistoryResponseCounterCnt.IncrBy(int64(len(data)))
		for _, item := range data {
			proc.HistoryResponseItemCnt.IncrBy(int64(len(item.Values)))
		}

		StdRender(w, data, nil)
	})

	// post, info
	http.HandleFunc("/graph/info", func(w http.ResponseWriter, r *http.Request) {
		// statistics
		proc.InfoRequestCnt.Incr()

		var body []*cmodel.GraphInfoParam
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&body)
		if err != nil {
			StdRender(w, "", err)
			return
		}

		if len(body) == 0 {
			StdRender(w, "", errors.New("empty"))
			return
		}

		data := []*cmodel.GraphFullyInfo{}
		for _, param := range body {
			if param == nil {
				continue
			}
			info, err := graph.Info(*param)
			if err != nil {
				log.Printf("graph.info fail, resp: %v, err: %v", info, err)
			}
			if info == nil {
				continue
			}
			data = append(data, info)
		}

		StdRender(w, data, nil)
	})

	// post, last
	http.HandleFunc("/graph/last", func(w http.ResponseWriter, r *http.Request) {
		// statistics
		proc.LastRequestCnt.Incr()

		var body []*cmodel.GraphLastParam
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&body)
		if err != nil {
			StdRender(w, "", err)
			return
		}

		if len(body) == 0 {
			StdRender(w, "", errors.New("empty"))
			return
		}

		data := []*cmodel.GraphLastResp{}
		for _, param := range body {
			if param == nil {
				continue
			}
			last, err := graph.Last(*param)
			if err != nil {
				log.Printf("graph.last fail, resp: %v, err: %v", last, err)
			}
			if last == nil {
				continue
			}
			data = append(data, last)
		}

		// statistics
		proc.LastRequestItemCnt.IncrBy(int64(len(data)))

		StdRender(w, data, nil)
	})

	// post, last/raw
	http.HandleFunc("/graph/last/raw", func(w http.ResponseWriter, r *http.Request) {
		// statistics
		proc.LastRawRequestCnt.Incr()

		var body []*cmodel.GraphLastParam
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&body)
		if err != nil {
			StdRender(w, "", err)
			return
		}

		if len(body) == 0 {
			StdRender(w, "", errors.New("empty"))
			return
		}

		data := []*cmodel.GraphLastResp{}
		for _, param := range body {
			if param == nil {
				continue
			}
			last, err := graph.LastRaw(*param)
			if err != nil {
				log.Printf("graph.last.raw fail, resp: %v, err: %v", last, err)
			}
			if last == nil {
				continue
			}
			data = append(data, last)
		}
		// statistics
		proc.LastRawRequestItemCnt.IncrBy(int64(len(data)))
		StdRender(w, data, nil)
	})

}
