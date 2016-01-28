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
    "github.com/Cepave/common/model"
    "net/rpc"
    "sync"
    "github.com/toolkits/net"
    "log"
    "math"
    "flag"
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
    NameTag  string
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

var connectionId = flag.String("connectionId", "", "The id of connection")


func main() {

    flag.Parse()

    urlPush := "http://10.20.30.40:1988/v1/push"

    var rpcClient SingleConnRpcClient

    rpcClient.RpcServer = "10.20.30.40:6030"

    if *connectionId == "" {
        *connectionId = "nqm-agent-1@10.20.30.40"
    }

    var req = model.NqmPingTaskRequest {
        ConnectionId: *connectionId,
        Hostname: "nqm-agent-1",
        IpAddress: "10.20.30.40",
    }
    var resp model.NqmPingTaskResponse

    fmt.Println("Request: ", req)

    err := rpcClient.Call("NqmAgent.PingTask", req, &resp,)
    if err != nil {
        fmt.Printf("Call NqmAgent.PingTask error: %v", err)
    }

    fmt.Printf("Agent: %v", resp.Agent)


    if !resp.NeedPing {
        fmt.Println(resp.NeedPing)
        return
    }

	fmt.Printf("resposne: %v", resp)
    resp.Command[6]="100"
    commandTemplate := resp.Command
    fmt.Println(req)
    fmt.Println("resp.Agent.Id=", resp.Agent.Id, "resp.Agent.IspId=", resp.Agent.IspId, "resp.Agent.ProvinceId=", resp.Agent.ProvinceId, "resp.Agent.CityId=", resp.Agent.CityId)
    fmt.Println("resp.Agent.Name=", resp.Agent.Name, "resp.Agent.IspName=", resp.Agent.IspName, "resp.Agent.ProvinceName=", resp.Agent.ProvinceName, "resp.Agent.CityName=", resp.Agent.CityName)
    fmt.Println(commandTemplate)
    var samplingTargetList []SamplingTarget
    for _, target := range resp.Targets {
        fmt.Println( target.Host, target.IspName, target.ProvinceName, target.CityName, target.NameTag)
        samplingTargetList = append(samplingTargetList, SamplingTarget{target.Host, target.IspName, target.ProvinceName, target.CityName, target.NameTag})
    }
    
    statisticsOfPacketsSent, statisticsOfPacketsReceived, statisticsOfTransmissionTime := Probe(samplingTargetList, commandTemplate)
    
    var params [] ParamToAgent
    params = append(params, FormParams(resp.Agent, "packets-sent" ,statisticsOfPacketsSent)...)
    params = append(params, FormParams(resp.Agent, "packets-received", statisticsOfPacketsReceived)...)
    params = append(params, FormParams(resp.Agent, "transmission-time", statisticsOfTransmissionTime)...)
    
    
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

func FormParams(agent *model.NqmAgent, metric string, statistics map[SamplingTarget]Statistic) []ParamToAgent {
    var params [] ParamToAgent
    
    for samplingTarget, statistic := range statistics {
        endpoint := *connectionId
        value := statistic.Value
        counterType := "GAUGE"
        tags := "nqm-agent-isp="+agent.IspName+
                ",nqm-agent-province="+agent.ProvinceName+
                ",nqm-agent-city="+agent.CityName+
                ",target-ip="+samplingTarget.Address+
                ",target-isp="+samplingTarget.ISP+
                ",target-province="+samplingTarget.Province+
                ",target-city="+samplingTarget.City+
                ",target-name-tag="+samplingTarget.NameTag
        timestamp := statistic.Timestamp
        step := int64(60)
        
        param := ParamToAgent{metric, endpoint, value, counterType, tags, timestamp, step}
        params = append(params, param)
    }
    return params
}

func Probe(samplingTargetList []SamplingTarget, commandTemplate []string) (map[SamplingTarget]Statistic, map[SamplingTarget]Statistic, map[SamplingTarget]Statistic) {
    statisticsOfPacketsSent := make(map[SamplingTarget]Statistic)
    statisticsOfPacketsReceived := make(map[SamplingTarget]Statistic)
    statisticsOfTransmissionTime := make(map[SamplingTarget]Statistic)
    var targetAddressList []string
    for _, samplingTarget := range samplingTargetList {
        targetAddressList = append(targetAddressList, samplingTarget.Address)
    }
    //commandTemplate example: {"fping", "-p", "20", "-i", "10", "-c", "100", "-q", "-a"}
    commandTemplate = append(commandTemplate, targetAddressList...)
    fpingCommand := exec.Command(commandTemplate[0], commandTemplate[1:]...)
    cmdOutput, err := fpingCommand.CombinedOutput()
    if err != nil {
        fmt.Println("error occured:")
        fmt.Printf("%s", err)
    }
    fpingResults := strings.Split(string(cmdOutput),"\n")
    fpingResults = fpingResults[:len(fpingResults)-1]
    for i, result := range fpingResults {
        fmt.Print("Result ", i+1, ": ")
        fmt.Println(result)
    }
    for i, fpingResult := range fpingResults {
        parsedFpingResult := strings.FieldsFunc(fpingResult, func(r rune) bool {
            switch r {
                case ' ', '\n', ':', '/', '%', '=', ',':
                    return true
            }
            return false
        })
        fmt.Println(parsedFpingResult)
        if len(parsedFpingResult) != 13 {
            delete(statisticsOfPacketsSent, samplingTargetList[i])
            delete(statisticsOfPacketsReceived, samplingTargetList[i])
            delete(statisticsOfTransmissionTime, samplingTargetList[i])
            continue
        }
        xmt, err := strconv.Atoi(parsedFpingResult[4])
        if err != nil {
            fmt.Println("error occured:")
            fmt.Printf("%s", err)
        }
        xmtStatistic := Statistic{time.Now().Unix(), xmt}
        statisticsOfPacketsSent[samplingTargetList[i]] = xmtStatistic

        rcv, err := strconv.Atoi(parsedFpingResult[5])
        if err != nil {
            fmt.Println("error occured:")
            fmt.Printf("%s", err)
        }
        rcvStatistic := Statistic{time.Now().Unix(), rcv}
        statisticsOfPacketsReceived[samplingTargetList[i]] = rcvStatistic
        
        tt, err := strconv.ParseFloat(parsedFpingResult[11],64)
        if err != nil {
            fmt.Println("error occured:")
            fmt.Printf("%s", err)
        }
        ttStatistic := Statistic{time.Now().Unix(), tt}
        statisticsOfTransmissionTime[samplingTargetList[i]] = ttStatistic
    }
    
    return statisticsOfPacketsSent, statisticsOfPacketsReceived, statisticsOfTransmissionTime
}


type SingleConnRpcClient struct {
	sync.Mutex
	rpcClient *rpc.Client
	RpcServer string
	Timeout   time.Duration
}

func (this *SingleConnRpcClient) close() {
	if this.rpcClient != nil {
		this.rpcClient.Close()
		this.rpcClient = nil
	}
}

func (this *SingleConnRpcClient) insureConn() {
	if this.rpcClient != nil {
		return
	}

	var err error
	var retry int = 1

	for {
		if this.rpcClient != nil {
			return
		}

		this.rpcClient, err = net.JsonRpcClient("tcp", this.RpcServer, this.Timeout)
		if err == nil {
			return
		}

		log.Printf("dial %s fail: %v", this.RpcServer, err)

		if retry > 6 {
			retry = 1
		}

		time.Sleep(time.Duration(math.Pow(2.0, float64(retry))) * time.Second)

		retry++
	}
}

func (this *SingleConnRpcClient) Call(method string, args interface{}, reply interface{}) error {

	this.Lock()
	defer this.Unlock()

	this.insureConn()

	timeout := time.Duration(50 * time.Second)
	done := make(chan error)

	go func() {
		err := this.rpcClient.Call(method, args, reply)
		done <- err
	}()

	select {
	case <-time.After(timeout):
		log.Printf("[WARN] rpc call timeout %v => %v", this.rpcClient, this.RpcServer)
		this.close()
	case err := <-done:
		if err != nil {
			this.close()
			return err
		}
	}

	return nil
}
