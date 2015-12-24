package main
import (
	"fmt"
//	"net/http"
	"encoding/json"
	"time"
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
    urlPush := "http://172.17.0.11:1988/v1/push"
	params := make([]*ParamToAgent, 0)
    metric := "transmission-time"
    endpoint := "nqm-endpoint"
    value := 100.0
    counterType := "GAUGE"
    tags := "view=1,ip=2.3.4.5"
    timestamp := time.Now().Unix()
    step := int64(30)
	
    params = append(params, &ParamToAgent{metric, endpoint, value, counterType, tags, timestamp, step})
    paramBody, err := json.Marshal(params)
    //if err != nil {
        fmt.Println(urlPush+", format body error,", err)
        
	//}

    fmt.Println("URL:>", urlPush)
    fmt.Println("URL:>", string(paramBody))
}