package service

import (
	"time"

	rd "github.com/Pallinder/go-randomdata"
	nlist "github.com/toolkits/container/list"
	tproc "github.com/toolkits/proc"

	cmodel "github.com/Cepave/open-falcon-backend/common/model"

	"github.com/Cepave/open-falcon-backend/modules/transfer/g"
	"github.com/Cepave/open-falcon-backend/modules/transfer/proc"
	"github.com/Cepave/open-falcon-backend/modules/transfer/sender"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("RecvMetricValues()", func() {
	Context("Every supported transportation got metrics", func() {
		BeforeEach(func() {
			judgeConfig := &g.JudgeConfig{
				Enabled: true, Replicas: 1,
				Cluster: map[string]string{ "test-judge-1": "10.20.30.1:6080" },
			}
			graphConfig := &g.GraphConfig{
				Enabled: true, Replicas: 1, Cluster: map[string]string{ "test-graph-1": "10.20.30.2:6070" },
				ClusterList: map[string]*g.ClusterNode {
					"test-graph-1": &g.ClusterNode{ Addrs: []string{ "10.20.30.2:6070" } },
				},
			}

			g.SetConfig(&g.GlobalConfig{
				Judge: judgeConfig,
				Graph: graphConfig,
				Tsdb: &g.TsdbConfig{ Enabled: true },
				Influxdb: &g.InfluxdbConfig{ Enabled: true },
				NqmRest: &g.NqmRestConfig { Enabled: true },
				Staging: &g.StagingConfig { Enabled: true, Filters: []string{ "pc02" } },
			})

			sender.SetNodeRings(judgeConfig, graphConfig)

			sender.JudgeQueues["test-judge-1"] = nlist.NewSafeListLimited(16)
			sender.GraphQueues["test-graph-110.20.30.2:6070"] = nlist.NewSafeListLimited(16)
			sender.InfluxdbQueues["default"] = nlist.NewSafeListLimited(16)

			sender.TsdbQueue = nlist.NewSafeListLimited(16)
			sender.NqmIcmpQueue = nlist.NewSafeListLimited(16)
			sender.NqmTcpQueue = nlist.NewSafeListLimited(16)
			sender.NqmTcpconnQueue = nlist.NewSafeListLimited(16)
			sender.StagingQueue = nlist.NewSafeListLimited(16)

			DefaultRelayStationFactory = NewRelayFactoryByGlobalConfig(g.Config())
		})
		AfterEach(func() {
			g.SetConfig(nil)
			sender.SetNodeRings(nil, nil)

			sender.JudgeQueues = make(map[string]*nlist.SafeListLimited)
			sender.GraphQueues = make(map[string]*nlist.SafeListLimited)
			sender.InfluxdbQueues = make(map[string]*nlist.SafeListLimited)

			sender.TsdbQueue = nil
			sender.NqmIcmpQueue = nil
			sender.NqmTcpQueue = nil
			sender.NqmTcpconnQueue = nil
			sender.StagingQueue = nil

			proc.RecvCnt       = tproc.NewSCounterQps("RecvCnt")
			proc.RpcRecvCnt    = tproc.NewSCounterQps("RpcRecvCnt")

			DefaultRelayStationFactory = nil
		})

		It("Every supported queue should have the expected number of metrics", func() {
			reply := &cmodel.TransferResponse{}
			err := RecvMetricValues(
				[]*cmodel.MetricValue {
					/**
					 * Three metrics
					 */
					{
						Endpoint: "pc01.it.cepave.com", Metric: "m01", Step: 30, Type: "GAUGE", Tags: "",
						Value: 11, Timestamp: time.Now().Unix() + 2,
					},
					{
						Endpoint: "pc02.it.cepave.com", Metric: "m01", Step: 30, Type: "GAUGE", Tags: "",
						Value: 12, Timestamp: time.Now().Unix() + 4,
					},
					{
						Endpoint: "pc02.it.cepave.com", Metric: "m01", Step: 30, Type: "GAUGE", Tags: "",
						Value: 13, Timestamp: time.Now().Unix() + 6,
					},
					// :~)
					/**
					 * 1 fping metrics
					 */
					{
						Endpoint: "pc02.it.cepave.com", Metric: "nqm-fping", Step: 10, Type: "GAUGE", Tags: "",
						Value: 21, Timestamp: time.Now().Unix() + 2,
					},
					// :~)
					/**
					 * 2 TCP ping metrics
					 */
					{
						Endpoint: "pc02.it.cepave.com", Metric: "nqm-tcpping", Step: 10, Type: "GAUGE", Tags: "",
						Value: 32, Timestamp: time.Now().Unix() + 2,
					},
					{
						Endpoint: "pc02.it.cepave.com", Metric: "nqm-tcpping", Step: 10, Type: "GAUGE", Tags: "",
						Value: 43, Timestamp: time.Now().Unix() + 4,
					},
					// :~)
					/**
					 * 3 TCP conn metrics
					 */
					{
						Endpoint: "pc01.it.cepave.com", Metric: "nqm-tcpconn", Step: 10, Type: "GAUGE", Tags: "",
						Value: 29, Timestamp: time.Now().Unix() + 2,
					},
					{
						Endpoint: "pc01.it.cepave.com", Metric: "nqm-tcpconn", Step: 10, Type: "GAUGE", Tags: "",
						Value: 33, Timestamp: time.Now().Unix() + 4,
					},
					{
						Endpoint: "pc01.it.cepave.com", Metric: "nqm-tcpconn", Step: 10, Type: "GAUGE", Tags: "",
						Value: 31, Timestamp: time.Now().Unix() + 6,
					},
					// :~)
				},
				reply, "rpc",
			)

			Expect(err).To(Succeed())
			Expect(reply.Total).To(Equal(9))
			Expect(reply.Invalid).To(BeEquivalentTo(0))

			/**
			 * Asserts the counter of statistics
			 */
			Expect(proc.RecvCnt.Cnt).To(BeEquivalentTo(9))
			Expect(proc.RpcRecvCnt.Cnt).To(BeEquivalentTo(9))
			Expect(proc.HttpRecvCnt.Cnt).To(BeEquivalentTo(0))
			// :~)

			/**
			 * Asserts the length of queue
			 */
			Expect(sender.JudgeQueues["test-judge-1"].Len()).To(Equal(3))
			Expect(sender.GraphQueues["test-graph-110.20.30.2:6070"].Len()).To(Equal(3))
			Expect(sender.InfluxdbQueues["default"].Len()).To(Equal(3))

			Expect(sender.TsdbQueue.Len()).To(Equal(3))
			Expect(sender.NqmIcmpQueue.Len()).To(Equal(1))
			Expect(sender.NqmTcpQueue.Len()).To(Equal(2))
			Expect(sender.NqmTcpconnQueue.Len()).To(Equal(3))

			Expect(sender.StagingQueue.Len()).To(Equal(5))
			// :~)
		})
	})
})

var _ = Describe("checkAndRefineMetric function", func() {
	Context("Success refine metric", func() {
		It("Refined metric should be as expected one", func() {
			metricTimestamp := time.Now().Add(-5 * time.Second)

			testedResult, invalidCount := checkAndRefineMetric(
				&cmodel.MetricValue {
					Metric: "disk.sda1.free", Endpoint: "pc807.net.tw",
					Step: 8, Type: "GAUGE", Tags: "uuid=non-1,interface=scsi",
					Timestamp: metricTimestamp.Unix(),
					Value: 301,
				},
				time.Now(),
			)

			Expect(invalidCount).To(Equal(0))
			Expect(testedResult).To(PointTo(MatchFields(
				IgnoreExtras,
				Fields {
					"Metric": Equal("disk.sda1.free"),
					"Endpoint": Equal("pc807.net.tw"),
					"Timestamp": Equal(metricTimestamp.Unix()),
					"Step": BeEquivalentTo(8),
					"CounterType": Equal("GAUGE"),
					"Tags": Equal(map[string]string {
						"uuid": "non-1",
						"interface": "scsi",
					}),
					"Value": Equal(301.0),
				},
			)))
		})

		Context("Timestamp of metric is unusual", func() {
			now := time.Now()

			DescribeTable("The value of timestamp should be the start time",
				func(testedValue int64) {
					testedResult, invalidCount := checkAndRefineMetric(
						&cmodel.MetricValue {
							Metric: "something.d1", Endpoint: "gk27.net.tw",
							Step: 8, Type: "GAUGE", Tags: "tag1=v1",
							Timestamp: testedValue,
							Value: 301,
						},
						now,
					)

					Expect(invalidCount).To(Equal(0))
					Expect(testedResult.Timestamp).To(Equal(now.Unix()))
				},
				Entry("The value is < 0", int64(-1)),
				Entry("The value is more than 2 hours in future", now.Add(121 * time.Minute).Unix()),
			)
		})
	})

	Context("Failure to refine metric", func() {
		var validMetric *cmodel.MetricValue

		BeforeEach(func() {
			validMetric = &cmodel.MetricValue {
				Metric: "disk.sda1.free", Endpoint: "pc807.net.tw",
				Step: 8, Type: "GAUGE", Tags: "uuid=non-1,interface=scsi",
				Timestamp: time.Now().Unix(),
				Value: 301,
			}
		})

		DescribeTable("The invalid count should be 1",
			func(setupFunc func(*cmodel.MetricValue)) {
				setupFunc(validMetric)

				testedResult, invalidCount := checkAndRefineMetric(validMetric, time.Now())

				Expect(testedResult).To(BeNil())
				Expect(invalidCount).To(Equal(1))
			},
			Entry("metric is \"kernel.hostname\"", func(metric *cmodel.MetricValue) { metric.Metric = "kernel.hostname" }),
			Entry("metric is \"\"", func(metric *cmodel.MetricValue) { metric.Metric = "" }),
			Entry("endpoint is \"\"", func(metric *cmodel.MetricValue) { metric.Endpoint = "" }),
			Entry("type is \"NN\"", func(metric *cmodel.MetricValue) { metric.Type = "NN" }),
			Entry("value is \"\"", func(metric *cmodel.MetricValue) { metric.Value = "" }),
			Entry("value is \"not-num\"", func(metric *cmodel.MetricValue) { metric.Value = "not-num" }),
			Entry("step is \"-1\"", func(metric *cmodel.MetricValue) { metric.Step = -1 }),
			Entry("len(metric) is \"511\"", func(metric *cmodel.MetricValue) { metric.Metric = rd.RandStringRunes(511) }),
		)
	})
})

var _ = Describe("refineValue() function", func() {
	Context("Refine to valid value", func() {
		DescribeTable("Refined value should be as expected one",
			func(source interface{}, expectedValue float64) {
				testedValue, ok := refineValue(source)

				Expect(ok).To(BeTrue())
				Expect(testedValue).To(Equal(expectedValue))
			},
			Entry("string value", "87.65", 87.65),
			Entry("float64 value", 33.76, 33.76),
			Entry("int value", int(7061), 7061.0),
			Entry("int8 value", int8(3), 3.0),
			Entry("int16 value", int16(876), 876.0),
			Entry("int32 value", int32(4003), 4003.0),
			Entry("int64 value", int64(99081), 99081.0),
		)
	})

	Context("Refine to invalid value", func() {
		DescribeTable("The success should be false",
			func(source interface{}) {
				testedValue, ok := refineValue(source)

				Expect(ok).To(BeFalse())
				Expect(testedValue).To(Equal(0.0))
			},
			Entry("in-convertible string", "ki99"),
			Entry("in-convertible type(bool)", true),
			Entry("in-convertible type(struct)", &struct { A int } { 33 }),
		)
	})
})
