package main
import (
	"fmt"
	"net/http"
	"encoding/json"
	"time"
    "bytes"
)

type ParamToAgent struct {
    Metric      string      `json:"metric"`
    Endpoint    string      `json:"endpoint"`
    Value       interface{} `json:"value"`       // number or string
    CounterType string      `json:"counterType"`
    Tags        string      `json:"tags"`
    Timestamp   int64       `json:"timestamp"`
    Step        int64       `json:"step"`
}
func main() {
    urlPush := "http://10.20.30.40:1988/v1/push"
	params := make([]*ParamToAgent, 0)
    metric := "transmission-time"
    endpoint := "nqm-endpoint"
    value := 100.0
    counterType := "GAUGE"
    tags := "view=1,ip=2.3.4.5"
    timestamp := time.Now().Unix()
    step := int64(30)
	
    params = append(params, &ParamToAgent{metric, endpoint, value, counterType, tags, timestamp, step})
    
    paramsBody, err := json.Marshal(params)
    if err != nil {
        fmt.Println(urlPush+", format body error,", err)    
	}

	postReq, err := http.NewRequest("POST", urlPush, bytes.NewBuffer(paramsBody))
	postReq.Header.Set("Content-Type", "application/json; charset=UTF-8")
	postReq.Header.Set("Connection", "close")
	
    httpClient := &http.Client{}
    postResp, err := httpClient.Do(postReq)
	if err != nil {
		fmt.Println(urlPush+", sending post request occurs error,", err)
	}
	defer postResp.Body.Close()
    
    fmt.Println("URL:", urlPush)
    fmt.Println("Body:", string(paramsBody))
}