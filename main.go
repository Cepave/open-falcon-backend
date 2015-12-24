package main
import (
	"fmt"
	"net/http"
	"encoding/json"
	"time"
    "bytes"
	"os/exec"
    "strings"
    "strconv"
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

type SamplingIPAddress struct {
    Address    string
    View       int
}

func main() {
    urlPush := "http://10.20.30.40:1988/v1/push"
	params := make([]*ParamToAgent, 0)
	samplingFrame := make([]*SamplingIPAddress, 0)
    
    samplingFrame = append(samplingFrame, &SamplingIPAddress{"8.8.8.8", 2})
	fpingCommand := exec.Command("fping", "-p", "20", "-i", "10", "-c", "4", "-q", "-a", "8.8.8.8")
    cmdOutput, err := fpingCommand.CombinedOutput()
    fpingResult := string(cmdOutput)
    fmt.Print(fpingResult)
    
    parsedFpingResult := strings.FieldsFunc(fpingResult, func(r rune) bool {
        switch r {
            case ':', '/', '=', ',':
                return true
        }
        return false
    })
    fmt.Println(parsedFpingResult)
    view := strconv.Itoa(samplingFrame[0].View)
    
    metric := "transmission-time"
    endpoint := "nqm-endpoint"
    value := parsedFpingResult[11]
    counterType := "GAUGE"
    tags := "view="+view+",ip="+samplingFrame[0].Address
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