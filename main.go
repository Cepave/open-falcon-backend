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


type SamplingTarget struct {
    Address     string
    ISP         string
    Province    string
    City        string
    ServerRoom  string
}


type Statistic struct {
    Timestamp      int64
    
    /**
     * Value could be:
     *     Packet Loss - int
     *     Transmission Time - float64
     */
    Value          interface{}
}

func main() {
    urlPush := "http://10.20.30.40:1988/v1/push"
    //connectionID := "nqm-agent@10.20.30.40"

	var samplingTargetList []SamplingTarget    
    samplingTargetList = append(samplingTargetList, SamplingTarget{"203.208.150.145", "ChinaTelecom", "Zhejiang", "Hangzhou", "Room-1"})
    samplingTargetList = append(samplingTargetList, SamplingTarget{"203.208.146.33", "ChinaTelecom", "Zhejiang", "Hangzhou", "Room-1"})
    samplingTargetList = append(samplingTargetList, SamplingTarget{"203.208.232.53", "ChinaTelecom", "Zhejiang", "Hangzhou", "Room-2"})
    samplingTargetList = append(samplingTargetList, SamplingTarget{"211.160.177.145", "Chinanet", "Zhejiang", "Wenzhou", "Room-1"})
    samplingTargetList = append(samplingTargetList, SamplingTarget{"211.160.177.150", "Chinanet", "Zhejiang", "Wenzhou", "Room-2"})
    samplingTargetList = append(samplingTargetList, SamplingTarget{"211.160.177.70", "Chinanet", "Zhejiang", "Wenzhou", "Room-2"})
    samplingTargetList = append(samplingTargetList, SamplingTarget{"61.190.194.110", "ChinaTelecom", "Anhui", "Hefei", "Room-1"})
    samplingTargetList = append(samplingTargetList, SamplingTarget{"61.190.194.114", "ChinaTelecom", "Anhui", "Hefei", "Room-2"})
    samplingTargetList = append(samplingTargetList, SamplingTarget{"61.190.194.222", "ChinaTelecom", "Anhui", "Hefei", "Room-2"})
    
    statisticsOfPacketsSent, statisticsOfPacketsReceived, statisticsOfTransmissionTime := Probe(samplingTargetList)
    
    var params [] ParamToAgent
    params = append(params, FormParams("packets-sent" ,statisticsOfPacketsSent)...)
    params = append(params, FormParams("packets-received", statisticsOfPacketsReceived)...)
    params = append(params, FormParams("transmission-time", statisticsOfTransmissionTime)...)
    
    
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

func FormParams(metric string, statistics map[SamplingTarget]Statistic) []ParamToAgent {
    var params [] ParamToAgent
    
    for samplingTarget, statistic := range statistics {
        endpoint := "nqm-agent-1@10.20.30.40"
        value := statistic.Value
        counterType := "GAUGE"
        tags := "target="+samplingTarget.Address+
                ",isp="+samplingTarget.ISP+
                ",province="+samplingTarget.Province+
                ",city="+samplingTarget.City+
                ",tag="+samplingTarget.ServerRoom
        timestamp := statistic.Timestamp
        step := int64(60)
        
        param := ParamToAgent{metric, endpoint, value, counterType, tags, timestamp, step}
        params = append(params, param)
    }
    return params
}

func Probe(samplingTargetList []SamplingTarget) (map[SamplingTarget]Statistic, map[SamplingTarget]Statistic, map[SamplingTarget]Statistic) {
    statisticsOfPacketsSent := make(map[SamplingTarget]Statistic)
    statisticsOfPacketsReceived := make(map[SamplingTarget]Statistic)
    statisticsOfTransmissionTime := make(map[SamplingTarget]Statistic)
    for _, samplingTarget := range samplingTargetList {
        fpingCommand := exec.Command("fping", "-p", "20", "-i", "10", "-c", "4", "-q", "-a", samplingTarget.Address)
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
            delete(statisticsOfPacketsSent, samplingTarget)
            delete(statisticsOfPacketsReceived, samplingTarget)
            delete(statisticsOfTransmissionTime, samplingTarget)
            continue
        }
        xmt, err := strconv.Atoi(parsedFpingResult[4])
        xmtStatistic := Statistic{time.Now().Unix(), xmt}
        statisticsOfPacketsSent[samplingTarget] = xmtStatistic

        rcv, err := strconv.Atoi(parsedFpingResult[5])
        rcvStatistic := Statistic{time.Now().Unix(), rcv}
        statisticsOfPacketsReceived[samplingTarget] = rcvStatistic
        
        tt, err := strconv.ParseFloat(parsedFpingResult[11],64)
        ttStatistic := Statistic{time.Now().Unix(), tt}
        statisticsOfTransmissionTime[samplingTarget] = ttStatistic
    }
    
    return statisticsOfPacketsSent, statisticsOfPacketsReceived, statisticsOfTransmissionTime
}