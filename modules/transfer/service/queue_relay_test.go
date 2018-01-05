package service

import (
	cmodel "github.com/Cepave/open-falcon-backend/common/model"

	"github.com/Cepave/open-falcon-backend/modules/transfer/g"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("RelayStation", func() {
	Context("Dispatch for various type of delegatee", func() {
		var stationFactory = &RelayStationFactory{
			stationBase: &RelayStation{
				Any: []RelayDelegatee{
					new(probeDelegatee), new(probeDelegatee),
				},
				Exclusive: []RelayDelegatee{
					&probeDelegatee{inclusiveName: "pc1"},
					&probeDelegatee{inclusiveName: "pc2"},
				},
				Otherwise: []RelayDelegatee{
					&probeDelegatee{},
				},
			},
		}
		var testedStation *RelayStation

		BeforeEach(func() {
			testedStation = stationFactory.Build()

			testedStation.Dispatch(&cmodel.MetaData{Endpoint: "zz1"})
			testedStation.Dispatch(&cmodel.MetaData{Endpoint: "zz2"})
			testedStation.Dispatch(&cmodel.MetaData{Endpoint: "pc1"})
			testedStation.Dispatch(&cmodel.MetaData{Endpoint: "pc1"})
			testedStation.Dispatch(&cmodel.MetaData{Endpoint: "pc2"})
			testedStation.Dispatch(&cmodel.MetaData{Endpoint: "pc2"})
			testedStation.Dispatch(&cmodel.MetaData{Endpoint: "pc2"})
			testedStation.Dispatch(&cmodel.MetaData{Endpoint: "ga3"})
			testedStation.Dispatch(&cmodel.MetaData{Endpoint: "ga4"})
		})
		AfterEach(func() {
			testedStation = nil
		})

		It("Every delegatee of [ANY] should accept all of the metrics", func() {
			assertAcceptedSize(testedStation.Any[0], 9)
			assertAcceptedSize(testedStation.Any[1], 9)
		})
		It("Every delegatee of [EXCLUSIVE] should accept expected metrics", func() {
			assertAcceptedSize(testedStation.Exclusive[0], 2)
			assertAcceptedSize(testedStation.Exclusive[1], 3)
		})
		It("Every delegatee of [OTHERWISE] should accept metrics which are not accepted by exclusive", func() {
			assertAcceptedSize(testedStation.Otherwise[0], 4)
		})
	})

	Context("Relay to assigned delegatee", func() {
		var stationFactory = &RelayStationFactory{
			stationBase: &RelayStation{
				Any: []RelayDelegatee{
					new(probeDelegatee), new(probeDelegatee),
				},
				Exclusive: []RelayDelegatee{
					&probeDelegatee{inclusiveName: "rc1"},
					&probeDelegatee{inclusiveName: "rc2"},
				},
				Otherwise: []RelayDelegatee{
					&probeDelegatee{},
				},
			},
		}
		var testedStation *RelayStation

		BeforeEach(func() {
			testedStation = stationFactory.Build()

			testedStation.Dispatch(&cmodel.MetaData{Endpoint: "a1"})
			testedStation.Dispatch(&cmodel.MetaData{Endpoint: "rc1"})
			testedStation.Dispatch(&cmodel.MetaData{Endpoint: "rc1"})
			testedStation.Dispatch(&cmodel.MetaData{Endpoint: "rc2"})
			testedStation.Dispatch(&cmodel.MetaData{Endpoint: "a3"})
		})
		AfterEach(func() {
			testedStation = nil
		})

		It("The RelayToQueue() should get called on delegatee", func() {
			testedStation.RelayToQueue()

			assertRelaySize(testedStation.Any[0], 5)
			assertRelaySize(testedStation.Any[1], 5)
			assertRelaySize(testedStation.Exclusive[0], 2)
			assertRelaySize(testedStation.Exclusive[1], 1)
			assertRelaySize(testedStation.Otherwise[0], 2)
		})
	})

	Context("clone() method", func() {
		sampleFactory := &RelayStationFactory{
			stationBase: &RelayStation{
				Any:       []RelayDelegatee{new(probeDelegatee)},
				Exclusive: []RelayDelegatee{&probeDelegatee{inclusiveName: "cc1"}},
				Otherwise: []RelayDelegatee{new(probeDelegatee)},
			},
		}

		It("Every type of delegatees should be CLONED, another station should be untouched", func() {
			station1 := sampleFactory.Build()
			station2 := sampleFactory.Build()

			station1.Dispatch(&cmodel.MetaData{Endpoint: "cc1"})
			station1.Dispatch(&cmodel.MetaData{Endpoint: "nothing.key"})

			assertAcceptedSize(station1.Any[0], 2)
			assertAcceptedSize(station1.Exclusive[0], 1)
			assertAcceptedSize(station1.Otherwise[0], 1)

			assertAcceptedSize(station2.Any[0], 0)
			assertAcceptedSize(station2.Exclusive[0], 0)
			assertAcceptedSize(station2.Otherwise[0], 0)
		})
	})
})

var _ = Describe("genericRelayPool(implementation of RelayDelegatee)", func() {
	Context("The relay target should get fed metrics", func() {
		target1 := &counterOfTarget{}
		target2 := &counterOfTarget{}

		testedPool := &genericRelayPool{
			relayTargets: &[]func([]*cmodel.MetaData){target1.accept, target2.accept},
		}

		It("Every target should get fed metrics", func() {
			testedPool.Accept(&cmodel.MetaData{})
			testedPool.Accept(&cmodel.MetaData{})
			testedPool.Accept(&cmodel.MetaData{})

			Expect(testedPool.RelayToQueue()).To(Equal(3))

			Expect(target1.count).To(Equal(3))
			Expect(target2.count).To(Equal(3))
		})
	})

	Context("Clone() method", func() {
		pool1 := &genericRelayPool{
			relayTargets: &[]func([]*cmodel.MetaData){
				func([]*cmodel.MetaData) {},
				func([]*cmodel.MetaData) {},
			},
		}
		pool2 := interface{}(pool1.Clone()).(*genericRelayPool)

		It("The metrics should not be affected", func() {
			pool1.Accept(&cmodel.MetaData{})

			Expect(pool1.metrics).To(HaveLen(1))
			Expect(pool2.metrics).To(HaveLen(0))
			Expect(*pool2.relayTargets).To(HaveLen(2))
		})
	})
})

var _ = Describe("stringMapRelayPool(implementation of RelayDelegatee)", func() {
	Context("Every different target should get different mapToMetrics", func() {
		target1 := &counterOfTarget{}
		target2 := &counterOfTarget{}

		testedPool := &stringMapRelayPool{
			stringify: func(metric *cmodel.MetaData) string {
				return metric.Endpoint
			},
			mapToTargets: map[string]func([]*cmodel.MetaData){
				"ep1": target1.accept,
				"ep2": target2.accept,
			},
			mapToMetrics: map[string][]*cmodel.MetaData{
				"ep1": make([]*cmodel.MetaData, 0),
				"ep2": make([]*cmodel.MetaData, 0),
			},
		}

		It("Target should accept expected number of mapToMetrics", func() {
			testedPool.Accept(&cmodel.MetaData{Endpoint: "ep1"})
			testedPool.Accept(&cmodel.MetaData{Endpoint: "ep1"})
			testedPool.Accept(&cmodel.MetaData{Endpoint: "ep2"})
			testedPool.Accept(&cmodel.MetaData{Endpoint: "ep2"})
			testedPool.Accept(&cmodel.MetaData{Endpoint: "ep2"})

			testedPool.RelayToQueue()

			Expect(target1.count).To(Equal(2))
			Expect(target2.count).To(Equal(3))
		})
	})

	Context("Clone() method", func() {
		pool1 := &stringMapRelayPool{
			stringify: func(metric *cmodel.MetaData) string {
				return metric.Endpoint
			},
			mapToTargets: map[string]func([]*cmodel.MetaData){
				"cp1": new(counterOfTarget).accept,
				"cp2": new(counterOfTarget).accept,
			},
			mapToMetrics: map[string][]*cmodel.MetaData{
				"cp1": make([]*cmodel.MetaData, 0),
				"cp2": make([]*cmodel.MetaData, 0),
			},
		}
		pool2 := interface{}(pool1.Clone()).(*stringMapRelayPool)

		It("Map should get cloned", func() {
			pool1.Accept(&cmodel.MetaData{Endpoint: "cp1"})
			pool1.Accept(&cmodel.MetaData{Endpoint: "cp1"})
			pool1.Accept(&cmodel.MetaData{Endpoint: "cp2"})
			pool1.Accept(&cmodel.MetaData{Endpoint: "cp2"})
			pool1.Accept(&cmodel.MetaData{Endpoint: "cp2"})

			Expect(pool1.mapToMetrics["cp1"]).To(HaveLen(2))
			Expect(pool1.mapToMetrics["cp2"]).To(HaveLen(3))
			Expect(pool2.mapToTargets).To(HaveLen(2))
			Expect(pool2.mapToMetrics["cp1"]).To(HaveLen(0))
			Expect(pool2.mapToMetrics["cp2"]).To(HaveLen(0))
		})
	})
})

var _ = Describe("filteredRelayPool(implementation of RelayDelegatee)", func() {
	Context("Filter mechanism", func() {
		targetPool := &probeDelegatee{}
		testedPool := &filteredRelayPool{
			RelayDelegatee: targetPool,
			filter: func(metric *cmodel.MetaData) bool {
				return metric.Endpoint == "Added"
			},
		}

		It("The number of accepted metrics should be 2(accepts 4)", func() {
			testedPool.Accept(&cmodel.MetaData{Endpoint: "Added"})
			testedPool.Accept(&cmodel.MetaData{Endpoint: "No-Added"})
			testedPool.Accept(&cmodel.MetaData{Endpoint: "Added"})
			testedPool.Accept(&cmodel.MetaData{Endpoint: "No-Added"})

			assertAcceptedSize(targetPool, 2)
		})
	})

	Context("Clone() method", func() {
		pool1 := &filteredRelayPool{
			RelayDelegatee: &probeDelegatee{},
			filter:         func(*cmodel.MetaData) bool { return true },
		}
		pool2 := interface{}(pool1.Clone()).(*filteredRelayPool)

		It("RelayDelegatee should be cloned", func() {
			pool1.Accept(&cmodel.MetaData{})
			pool1.Accept(&cmodel.MetaData{})

			assertAcceptedSize(pool1.RelayDelegatee, 2)
			assertAcceptedSize(pool2.RelayDelegatee, 0)
		})
	})
})

var _ = Describe("stageRelayPool(implementation of RelayDelegatee)", func() {
	Context("Clone() method", func() {
		pool1 := &stageRelayPool{}
		pool2 := interface{}(pool1.Clone()).(*stageRelayPool)

		It("The metrics should not be affected by source", func() {
			pool1.Accept(&cmodel.MetaData{SourceMetric: &cmodel.MetricValue{}})
			pool1.Accept(&cmodel.MetaData{SourceMetric: &cmodel.MetricValue{}})

			Expect(pool1.metrics).To(HaveLen(2))
			Expect(pool2.metrics).To(HaveLen(0))
		})
	})
})

var _ = Describe("NewRelayFactoryByGlobalConfig()", func() {
	var sampleConfig *g.GlobalConfig

	BeforeEach(func() {
		sampleConfig = &g.GlobalConfig{
			Judge:    &g.JudgeConfig{},
			Graph:    &g.GraphConfig{},
			Tsdb:     &g.TsdbConfig{},
			Influxdb: &g.InfluxdbConfig{},
			NqmRest:  &g.NqmRestConfig{},
			Staging:  &g.StagingConfig{},
		}
	})
	AfterEach(func() {
		sampleConfig = nil
	})

	Context("Generic metrics(otherwise RelayDelegatees)", func() {
		DescribeTable("The number of targets should be as expected",
			func(configSetup func(*g.GlobalConfig), expectedNumber int) {
				configSetup(sampleConfig)
				testedFactory := NewRelayFactoryByGlobalConfig(sampleConfig)

				if expectedNumber == 0 {
					Expect(testedFactory.stationBase.Otherwise).To(HaveLen(0))
					return
				}

				Expect(testedFactory.stationBase.Otherwise).To(HaveLen(1))
				testedPool := interface{}(testedFactory.stationBase.Otherwise[0]).(*genericRelayPool)
				Expect(*testedPool.relayTargets).To(HaveLen(expectedNumber))
			},
			Entry("Graph and Judge", func(config *g.GlobalConfig) {
				config.Judge.Enabled = true
				config.Graph.Enabled = true
			}, 2),
			Entry("All of supported queue are enabled", func(config *g.GlobalConfig) {
				config.Judge.Enabled = true
				config.Graph.Enabled = true
				config.Tsdb.Enabled = true
				config.Influxdb.Enabled = true
			}, 4),
			Entry("Tsdb", func(config *g.GlobalConfig) {
				config.Tsdb.Enabled = true
			}, 1),
			Entry("Nothing", func(config *g.GlobalConfig) {}, 0),
		)
	})

	Context("NQM metrics(exclusive RelayDelegatees)", func() {
		DescribeTable("The number of targets should be as expected",
			func(enabled bool) {
				sampleConfig.NqmRest.Enabled = enabled
				testedFactory := NewRelayFactoryByGlobalConfig(sampleConfig)

				if !enabled {
					Expect(testedFactory.stationBase.Exclusive).To(HaveLen(0))
					return
				}

				Expect(testedFactory.stationBase.Exclusive).To(HaveLen(1))
				testedPool := interface{}(testedFactory.stationBase.Exclusive[0]).(*stringMapRelayPool)

				Expect(testedPool.mapToTargets).To(And(
					HaveKey("nqm-fping"), HaveKey("nqm-tcpconn"), HaveKey("nqm-tcpping"),
				))
				Expect(testedPool.mapToMetrics).To(And(
					HaveKey("nqm-fping"), HaveKey("nqm-tcpconn"), HaveKey("nqm-tcpping"),
				))
			},
			Entry("Enabled NQM", true),
			Entry("Disabled NQM", false),
		)
	})
})

type counterOfTarget struct {
	count int
}

func (t *counterOfTarget) accept(metrics []*cmodel.MetaData) {
	t.count = len(metrics)
}

func assertAcceptedSize(testedDelegatee interface{}, expectedSize int) {
	nativeDelegatee := testedDelegatee.(*probeDelegatee)
	ExpectWithOffset(1, nativeDelegatee.acceptedSize).To(Equal(expectedSize))
}
func assertRelaySize(testedDelegatee interface{}, expectedSize int) {
	nativeDelegatee := testedDelegatee.(*probeDelegatee)
	ExpectWithOffset(1, nativeDelegatee.relaySize).To(Equal(expectedSize))
}

type probeDelegatee struct {
	acceptedSize  int
	relaySize     int
	inclusiveName string
}

func (p *probeDelegatee) Accept(metric *cmodel.MetaData) bool {
	if p.inclusiveName != "" {
		if p.inclusiveName == metric.Endpoint {
			p.acceptedSize++
			return true
		}

		return false
	}

	p.acceptedSize++
	return true
}
func (p *probeDelegatee) RelayToQueue() int {
	p.relaySize = p.acceptedSize
	return p.acceptedSize
}
func (p *probeDelegatee) Clone() RelayDelegatee {
	newSelf := *p
	return &newSelf
}
