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
		regx, _ := regexp.Compile("(\\.\\$\\s*|\\s*)$")
		data := []*cmodel.GraphQueryResponse{}
		var result *cmodel.GraphQueryResponse
		isPacketLossRate := false
		for _, ec := range body.EndpointCounters {
			if strings.Contains(ec.Counter,"packet-loss-rate") {
				isPacketLossRate = true
				break
			}
		}
		isAverage := false
		for _, ec := range body.EndpointCounters {
			if strings.Contains(ec.Counter,"average") {
				isAverage = true
				break
			}
		}
		if isPacketLossRate {
			/**
			 * 下面這段，只是想先做一個跟 packets-sent 的 response 一樣的 struct
			 */
			var packetSentCount []cmodel.JsonFloat
			for _, ec := range body.EndpointCounters {
				if strings.Contains(ec.Counter,"packets-sent") {
					request := cmodel.GraphQueryParam{
						Start:     int64(body.Start),
						End:       int64(body.End),
						ConsolFun: body.CF,
						Endpoint:  ec.Endpoint,
						Counter:   ec.Counter,
					}
					result, err = graph.QueryOne(request)
					for i := range result.Values {
						result.Values[i].Value = cmodel.JsonFloat(0.0)
						packetSentCount = append(packetSentCount, cmodel.JsonFloat(0.0))
					}
					if err != nil {
						log.Printf("graph.queryOne fail, %v", err)
					}
					break
				}
			}
			
			for _, ec := range body.EndpointCounters {
				/**
				 * 此版本中，在 dashboard 查詢 packet-loss-rate 的時候，
				 * dashboard 會以 packets-sent 這個 metric 來呈現搜尋到的結果。
				 */	
				if strings.Contains(ec.Counter,"packets-sent") {
					/**
					 * 此版本中，packet-loss-rate 的 data，
					 * 必須跟 packets-sent & packets-eceived 一起看。
					 */					
					requestPacketsSent := cmodel.GraphQueryParam{
						Start:     int64(body.Start),
						End:       int64(body.End),
						ConsolFun: body.CF,
						Endpoint:  ec.Endpoint,
						Counter:   ec.Counter,
					}
					resultPacketsSent, err := graph.QueryOne(requestPacketsSent)
					if err != nil {
						log.Printf("graph.queryOne fail, %v", err)
					}
					if resultPacketsSent == nil {
						continue
					}
					data = append(data, resultPacketsSent)

					requestPacketReceived := cmodel.GraphQueryParam{
						Start:     int64(body.Start),
						End:       int64(body.End),
						ConsolFun: body.CF,
						Endpoint:  ec.Endpoint,
						Counter:   strings.Replace(ec.Counter, "packets-sent", "packets-received", 1),
					}
					resultPacketReceived, err := graph.QueryOne(requestPacketReceived)
					if err != nil {
						log.Printf("graph.queryOne fail, %v", err)
					}
					if resultPacketReceived == nil {
						continue
					}
					data = append(data, resultPacketReceived)
					for i := range resultPacketsSent.Values {
						packetLossCount := (resultPacketsSent.Values[i].Value		-
											resultPacketReceived.Values[i].Value)
						result.Values[i].Value += packetLossCount
						packetSentCount[i] += resultPacketsSent.Values[i].Value
					}
				}

			}
			
			result.Endpoint = "all-endpoints"
			result.Counter = "packet-loss-rate"
			for i := range result.Values {
				result.Values[i].Value = result.Values[i].Value/packetSentCount[i]
			}
			result.Values = result.Values
			data = append(data, result)
							
		} else if isAverage {
			/**
			 * 下面這段，只是想先製造一個跟 transmission-time 的 response 一樣的 struct
			 */
			var packetSentCount []cmodel.JsonFloat
			for _, ec := range body.EndpointCounters {
				if strings.Contains(ec.Counter,"transmission-time") {
					request := cmodel.GraphQueryParam{
						Start:     int64(body.Start),
						End:       int64(body.End),
						ConsolFun: body.CF,
						Endpoint:  ec.Endpoint,
						Counter:   ec.Counter,
					}
					result, err = graph.QueryOne(request)
					for i := range result.Values {
						result.Values[i].Value = cmodel.JsonFloat(0.0)
						packetSentCount = append(packetSentCount, cmodel.JsonFloat(0.0))
					}
					if err != nil {
						log.Printf("graph.queryOne fail, %v", err)
					}
					break
				}
			}
			for _, ec := range body.EndpointCounters {
				/**
				 * 此版本中，在 dashboard 查詢 average 的時候，
				 * dashboard 會以 transmission-time 這個 metric 來呈現搜尋到的結果。
				 */	
				if strings.Contains(ec.Counter,"transmission-time") {
					/**
					 * 此版本中，average 的 data，
					 * 必須跟 transmission-time & packets-sent 一起看。
					 */	
					requestTransmissionTime := cmodel.GraphQueryParam{
						Start:     int64(body.Start),
						End:       int64(body.End),
						ConsolFun: body.CF,
						Endpoint:  ec.Endpoint,
						Counter:   ec.Counter,
					}
					resultTransmissionTime, err := graph.QueryOne(requestTransmissionTime)
					if err != nil {
						log.Printf("graph.queryOne fail, %v", err)
					}
					if resultTransmissionTime == nil {
						continue
					}
					data = append(data, resultTransmissionTime)

					requestPacketsSent := cmodel.GraphQueryParam{
						Start:     int64(body.Start),
						End:       int64(body.End),
						ConsolFun: body.CF,
						Endpoint:  ec.Endpoint,
						Counter:   strings.Replace(ec.Counter, "transmission-time", "packets-sent", 1),
					}
					resultPacketsSent, err := graph.QueryOne(requestPacketsSent)
					if err != nil {
						log.Printf("graph.queryOne fail, %v", err)
					}
					if resultPacketsSent == nil {
						continue
					}
					data = append(data, resultPacketsSent)
					for i := range resultTransmissionTime.Values {
						result.Values[i].Value += resultTransmissionTime.Values[i].Value * resultPacketsSent.Values[i].Value
						packetSentCount[i] += resultPacketsSent.Values[i].Value
					}
				}

			}
			
			result.Endpoint = "all-endpoints"
			result.Counter = "average"
			for i := range result.Values {
				result.Values[i].Value = result.Values[i].Value/packetSentCount[i]
			}
			result.Values = result.Values
			data = append(data, result)
			
		} else {
			for _, ec := range body.EndpointCounters {
				request := cmodel.GraphQueryParam{
					Start:     int64(body.Start),
					End:       int64(body.End),
					ConsolFun: body.CF,
					Endpoint:  ec.Endpoint,
					Counter:   ec.Counter,
				}
				result, err := graph.QueryOne(request)
				if err != nil {
					log.Printf("graph.queryOne fail, %v", err)
				}
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
