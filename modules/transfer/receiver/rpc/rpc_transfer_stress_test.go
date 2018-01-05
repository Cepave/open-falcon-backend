package rpc

import (
	"fmt"
	"net/rpc"
	"time"

	rd "github.com/Pallinder/go-randomdata"

	cmodel "github.com/Cepave/open-falcon-backend/common/model"
	trpc "github.com/Cepave/open-falcon-backend/common/testing/jsonrpc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	number_of_metrics       = 0
	number_of_threads       = 16
	number_of_agents        = 1024
	times_of_measure        = 3
	sec_interval_of_metrics = 30 // Seconds
)

var _ = Describe("Stressing test on receiving metrics", func() {
	var availableThreads chan bool

	BeforeEach(func() {
		jsonRpcSkipper.Skip()

		if number_of_metrics == 0 {
			Skip(`Skip pressure because "number_of_metrics" is 0`)
		}

		availableThreads = make(chan bool, number_of_threads)
	})

	AfterEach(func() {
		close(availableThreads)
	})

	Measure(fmt.Sprintf("%d metrics over %d agents", number_of_metrics, number_of_agents), func(b Benchmarker) {
		b.Time("runtime", func() {
			for i := 0; i < number_of_agents; i++ {
				availableThreads <- true
				go func(agentNumber int) {
					defer GinkgoRecover()
					defer func() {
						<-availableThreads
					}()

					resp, err := sendMetrics(number_of_metrics)

					Expect(err).To(Succeed())
					Expect(resp.Message).To(Equal("ok"))
					Expect(resp.Invalid).To(Equal(0))
					Expect(resp.Total).To(Equal(number_of_metrics))
				}(i)
			}

			/**
			 * Waiting for all of the go routines has completed
			 */
			for i := 0; i < number_of_threads; i++ {
				availableThreads <- true
			}
			// :~)
		})
	}, times_of_measure)
})

func sendMetrics(numberOfMetrics int) (*cmodel.TransferResponse, error) {
	var (
		reply *cmodel.TransferResponse
		err   error
	)

	ginkgoClient := &trpc.GinkgoJsonRpc{}
	ginkgoClient.OpenClient(
		func(client *rpc.Client) {
			reply = &cmodel.TransferResponse{}
			err = client.Call("Transfer.Update", buildMetrics(numberOfMetrics), reply)
		},
	)

	return reply, err
}

var (
	sampleMetrics = []string{
		"cpu.idle", "cpu.busy", "disk.out.peak", "disk.in.peak",
		"net.out.bytes", "net.in.bytes", "net.drop.bytes",
		"io.out.peak", "io.in.peak",
	}

	valueRanges = [][]int{
		{1, 100},
		{500, 2500},
		{1, 50000},
		{10000, 500000},
	}
)

func buildMetrics(numberOfMetrics int) []*cmodel.MetricValue {
	metrics := make([]*cmodel.MetricValue, numberOfMetrics)

	endpoint := fmt.Sprintf("pc%03d.gg%03d.net.tw", rd.Number(1, 999), rd.Number(1, 999))
	metric := rd.StringSample(sampleMetrics...)
	valueRange := valueRanges[rd.Number(len(valueRanges))]
	step := int64(rd.Number(1, 6) * 5)

	startTime := time.Now().Add(time.Duration(-sec_interval_of_metrics*numberOfMetrics) * time.Second)

	for i := 0; i < numberOfMetrics; i++ {
		metrics[i] = &cmodel.MetricValue{
			Endpoint: endpoint, Metric: metric,
			Step: step, Type: "GAUGE", Tags: "",
			Timestamp: startTime.Add(time.Duration(sec_interval_of_metrics*i) * time.Second).Unix(),
			Value:     rd.Number(valueRange[0], valueRange[1]),
		}
	}

	return metrics
}
