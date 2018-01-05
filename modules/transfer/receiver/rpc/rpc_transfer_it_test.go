package rpc

import (
	"net/rpc"
	"time"

	cmodel "github.com/Cepave/open-falcon-backend/common/model"
	trpc "github.com/Cepave/open-falcon-backend/common/testing/jsonrpc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Sends metrics by RPC protocol", jsonRpcSkipper.PrependBeforeEach(func() {
	Context("Sending metrics", func() {
		It("There should be no error", func() {
			ginkgoClient := &trpc.GinkgoJsonRpc{}
			ginkgoClient.OpenClient(
				func(client *rpc.Client) {
					reply := &cmodel.TransferResponse{}

					err := client.Call(
						"Transfer.Update",
						[]*cmodel.MetricValue{
							{
								Endpoint: "pc01.it.cepave.com", Metric: "m01", Step: 30, Type: "GAUGE", Tags: "",
								Value: 11, Timestamp: time.Now().Unix() + 2,
							},
							{
								Endpoint: "pc01.it.cepave.com", Metric: "m01", Step: 30, Type: "GAUGE", Tags: "",
								Value: 12, Timestamp: time.Now().Unix() + 4,
							},
							{
								Endpoint: "pc01.it.cepave.com", Metric: "m01", Step: 30, Type: "GAUGE", Tags: "",
								Value: 13, Timestamp: time.Now().Unix() + 6,
							},
						},
						reply,
					)

					Expect(err).To(Succeed())
					GinkgoT().Logf("Reply from \"Transfer.Update\": %#v", reply)
					Expect(reply).To(PointTo(MatchFields(IgnoreExtras, Fields{
						"Message": Equal("ok"),
						"Invalid": Equal(0),
						"Total":   Equal(3),
					})))
				},
			)
		})
	})
}))
