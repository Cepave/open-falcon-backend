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
    Value       interface{} `json:"value"`
    CounterType string      `json:"counterType"`
    Tags        string      `json:"tags"`
    Timestamp   int64       `json:"timestamp"`
    Step        int64       `json:"step"`
}

type SamplingUnit struct {
    Address    string
    View       int
}

type Sample struct {
    SampledUnit    SamplingUnit
    Timestamp      int64
    
    /**
     * Packet Loss - int
     * Transmission Time - float64
     */
    Value          interface{}
}

func main() {
    urlPush := "http://10.20.30.40:1988/v1/push"

	var samplingFrame []SamplingUnit    
    samplingFrame = append(samplingFrame, SamplingUnit{"8.8.8.8", 2})
    samplingFrame = append(samplingFrame, SamplingUnit{"8.8.4.4", 3})
    
    samplesOfPacketLoss, samplesOfTransmissionTime := Sampling(samplingFrame)
    
    var params [] ParamToAgent
    params = append(params, FormParams("packet-loss" ,samplesOfPacketLoss)...)
    params = append(params, FormParams("transmission-time", samplesOfTransmissionTime)...)
    
    
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

func FormParams(metric string, samples map[string]Sample) []ParamToAgent {
    var params [] ParamToAgent
    
    for address, sample := range samples {
        view := strconv.Itoa(sample.SampledUnit.View)

        endpoint := "nqm-endpoint"
        value := sample.Value
        counterType := "GAUGE"
        tags := "view="+view+",ip="+address
        timestamp := sample.Timestamp
        step := int64(30)
        
        param := ParamToAgent{metric, endpoint, value, counterType, tags, timestamp, step}
        params = append(params, param)
    }
    return params
}

func Sampling(samplingFrame []SamplingUnit) ( map[string]Sample,  map[string]Sample) {
    samplesOfPacketLoss := make(map[string]Sample)
    samplesOfTransmissionTime := make(map[string]Sample)
    for _, samplingUnit := range samplingFrame {
        fpingCommand := exec.Command("fping", "-p", "20", "-i", "10", "-c", "4", "-q", "-a", samplingUnit.Address)
        cmdOutput, err := fpingCommand.CombinedOutput()
        if err != nil {
            fmt.Println("error occured:")
            fmt.Printf("%s", err)
        }
        fpingResult := string(cmdOutput)
        fmt.Print(fpingResult)
        
        parsedFpingResult := strings.FieldsFunc(fpingResult, func(r rune) bool {
            switch r {
                case ' ', '\n', ':', '/', '%', '=', ',':
                    return true
            }
            return false
        })
        fmt.Println(parsedFpingResult)
        if len(parsedFpingResult) != 13 {
            delete(samplesOfPacketLoss, samplingUnit.Address)
            delete(samplesOfTransmissionTime, samplingUnit.Address)
            continue
        }
        xmt, err := strconv.Atoi(parsedFpingResult[4])
        rcv, err := strconv.Atoi(parsedFpingResult[5])
        samplePL := Sample{samplingUnit, time.Now().Unix(), xmt - rcv}
        samplesOfPacketLoss[samplingUnit.Address] = samplePL
        
        tt, err := strconv.ParseFloat(parsedFpingResult[11],64)
        sampleTT := Sample{samplingUnit, time.Now().Unix(), tt}
        samplesOfTransmissionTime[samplingUnit.Address] = sampleTT
    }
    
    return samplesOfPacketLoss, samplesOfTransmissionTime
}