package service

import (
	"strings"

	cmodel "github.com/Cepave/open-falcon-backend/common/model"

	"github.com/Cepave/open-falcon-backend/modules/transfer/g"
	"github.com/Cepave/open-falcon-backend/modules/transfer/sender"
)

var DefaultRelayStationFactory *RelayStationFactory

// This interface is used by "QueryRelay".
type RelayFunc func(mapToMetrics []*cmodel.MetaData)

// This interface used by "RelayStation"
type RelayDelegatee interface {
	Accept(*cmodel.MetaData) bool
	RelayToQueue() int
	Clone() RelayDelegatee
}

// This factory is used to build independent instance of "*RelayStation",
// which is "CLONED" by pre-defined one.
type RelayStationFactory struct {
	stationBase *RelayStation
}

func (f *RelayStationFactory) Build() *RelayStation {
	return f.stationBase.clone()
}

// Builds "*RelayStationFactory" from instance of "*g.GlobalConfig"
func NewRelayFactoryByGlobalConfig(config *g.GlobalConfig) *RelayStationFactory {
	stationBase := new(RelayStation)

	/**
	 * Sets-up the delegates for generic metrcis by enabled functions.
	 */
	genericTargets := make([]func([]*cmodel.MetaData), 0)

	if config.Graph.Enabled {
		genericTargets = append(genericTargets, sender.Push2GraphSendQueue)
	}
	if config.Judge.Enabled {
		genericTargets = append(genericTargets, sender.Push2JudgeSendQueue)
	}
	if config.Tsdb.Enabled {
		genericTargets = append(genericTargets, sender.Push2TsdbSendQueue)
	}
	if config.Influxdb.Enabled {
		genericTargets = append(genericTargets, sender.Push2InfluxdbSendQueue)
	}

	if len(genericTargets) > 0 {
		stationBase.Otherwise = append(stationBase.Otherwise, &genericRelayPool{relayTargets: &genericTargets})
	}
	// :~)

	/**
	 * Sets-up the delegates for NQM mapToMetrics.
	 */
	if config.NqmRest.Enabled {
		nqmRelayPool := buildRelayPoolForMetricMap()
		nqmRelayPool.mapToTargets = &map[string]func([]*cmodel.MetaData){
			"nqm-fping":   sender.Push2NqmIcmpSendQueue,
			"nqm-tcpconn": sender.Push2NqmTcpconnSendQueue,
			"nqm-tcpping": sender.Push2NqmTcpSendQueue,
		}
		nqmRelayPool.mapToMetrics = map[string][]*cmodel.MetaData{
			"nqm-fping":   make([]*cmodel.MetaData, 0),
			"nqm-tcpconn": make([]*cmodel.MetaData, 0),
			"nqm-tcpping": make([]*cmodel.MetaData, 0),
		}

		stationBase.Exclusive = append(stationBase.Exclusive, nqmRelayPool)
	}
	// :~)

	/**
	 * Sets-up the re-sending of mapToMetrics to staging environment.
	 */
	if config.Staging.Enabled {
		stationBase.Any = append(
			stationBase.Any,
			buildRelayPoolForEffectiveFilterOnEndpoint(new(stageRelayPool), config.Staging.Filters),
		)
	}
	// :~)

	return &RelayStationFactory{stationBase}
}

// This object handles the input metric and try to dispatch metric to
// corresponding RelayDelegatees.
type RelayStation struct {
	// All of the metrics would be fed into "Accept(*MetaData)" in this bucket
	Any []RelayDelegatee
	// If a metric fit("Accept(*MetaData)" returns "true"), it would be consumed
	Exclusive []RelayDelegatee
	// All of the metrics which are not consumed by "Exclusive" bucket, would be sent to this bucket
	Otherwise []RelayDelegatee
}

func (r *RelayStation) Dispatch(metaMetric *cmodel.MetaData) {
	for _, delegatee := range r.Any {
		delegatee.Accept(metaMetric)
	}

	for _, delegatee := range r.Exclusive {
		if delegatee.Accept(metaMetric) {
			return
		}
	}

	for _, delegatee := range r.Otherwise {
		delegatee.Accept(metaMetric)
	}
}
func (r *RelayStation) RelayToQueue() {
	r.iterateAll(func(delegatee RelayDelegatee) {
		delegatee.RelayToQueue()
	})
}
func (r *RelayStation) clone() *RelayStation {
	newStation := new(RelayStation)

	newStation.Any = cloneRelayDelegatees(r.Any)
	newStation.Exclusive = cloneRelayDelegatees(r.Exclusive)
	newStation.Otherwise = cloneRelayDelegatees(r.Otherwise)

	return newStation
}
func (r *RelayStation) iterateAll(callback func(RelayDelegatee)) {
	for _, delegatee := range r.Exclusive {
		callback(delegatee)
	}

	for _, delegatee := range r.Otherwise {
		callback(delegatee)
	}

	for _, delegatee := range r.Any {
		callback(delegatee)
	}
}

func cloneRelayDelegatees(source []RelayDelegatee) []RelayDelegatee {
	newDelegatees := make([]RelayDelegatee, 0, len(source))

	for _, delegatee := range source {
		newDelegatees = append(newDelegatees, delegatee.Clone())
	}

	return newDelegatees
}

// This object provides mechanism on SINGLE queue of metrics, MULTIPLE RelayDelegatees.
type genericRelayPool struct {
	// This slice would be re-use across "Clone()" method
	relayTargets *[]func([]*cmodel.MetaData)
	metrics      []*cmodel.MetaData
}

func (p *genericRelayPool) Accept(metric *cmodel.MetaData) bool {
	p.metrics = append(p.metrics, metric)
	return true
}
func (p *genericRelayPool) RelayToQueue() int {
	for _, target := range *p.relayTargets {
		target(p.metrics)
	}

	return len(p.metrics)
}
func (p *genericRelayPool) Clone() RelayDelegatee {
	newSelf := *p
	return &newSelf
}

// This object provides mechanism on MULTIPLE pairs on string(key) and RelayDelegatee
type stringMapRelayPool struct {
	// This map would be re-use across "Clone()" method
	mapToTargets *map[string]func([]*cmodel.MetaData)
	mapToMetrics map[string][]*cmodel.MetaData
	stringify    func(*cmodel.MetaData) string
}

func (p *stringMapRelayPool) Accept(metric *cmodel.MetaData) bool {
	key := p.stringify(metric)

	targetQueue, ok := p.mapToMetrics[key]

	if !ok {
		return false
	}

	targetQueue = append(targetQueue, metric)
	p.mapToMetrics[key] = targetQueue

	return true
}
func (p *stringMapRelayPool) RelayToQueue() int {
	counter := 0

	for name, target := range *p.mapToTargets {
		queuedMetrics := p.mapToMetrics[name]

		counter += len(queuedMetrics)
		if len(queuedMetrics) > 0 {
			target(queuedMetrics)
		}
	}

	return counter
}
func (p *stringMapRelayPool) Clone() RelayDelegatee {
	newSelf := *p

	/**
	 * Clones the maps
	 */
	newMapToMetrics := make(map[string][]*cmodel.MetaData)
	for k, v := range p.mapToMetrics {
		newMapToMetrics[k] = v
	}
	newSelf.mapToMetrics = newMapToMetrics
	// :~)

	return &newSelf
}

// Builds "*stringMapRelayPool" on value of *MetaData.Metric
func buildRelayPoolForMetricMap() *stringMapRelayPool {
	return &stringMapRelayPool{
		stringify: func(metric *cmodel.MetaData) string {
			return metric.Metric
		},
	}
}

// Relays *MetricValue(instead of *MetaData) to certain queue
type stageRelayPool struct {
	metrics []*cmodel.MetricValue
}

func (p *stageRelayPool) Accept(metric *cmodel.MetaData) bool {
	p.metrics = append(p.metrics, metric.SourceMetric)
	return true
}
func (p *stageRelayPool) RelayToQueue() int {
	sender.Push2StagingSendQueue(p.metrics)
	return len(p.metrics)
}
func (p *stageRelayPool) Clone() RelayDelegatee {
	newSelf := *p
	return &newSelf
}

// This object filters the accepted metric, if the result from filter is true,
// calling the "Accept(*MetaData)" method of contained RelayDelegatee.
type filteredRelayPool struct {
	RelayDelegatee
	filter func(*cmodel.MetaData) bool
}

func (p *filteredRelayPool) Accept(metric *cmodel.MetaData) bool {
	if p.filter(metric) {
		p.RelayDelegatee.Accept(metric)
		return true
	}

	return false
}
func (p *filteredRelayPool) Clone() RelayDelegatee {
	newSelf := *p
	newSelf.RelayDelegatee = p.RelayDelegatee.Clone()
	return &newSelf
}

// Builds "*filteredRelayPool" by prefix of endpoint(of "*MetaData")
func buildRelayPoolForEffectiveFilterOnEndpoint(targetDelegatee RelayDelegatee, filters []string) RelayDelegatee {
	finalDelegatee := targetDelegatee

	if len(filters) > 0 &&
		!(len(filters) == 1 && filters[0] == "") {
		finalDelegatee = &filteredRelayPool{
			RelayDelegatee: targetDelegatee,
			filter: func(metric *cmodel.MetaData) bool {
				for _, filter := range filters {
					if strings.HasPrefix(metric.Endpoint, filter) {
						return true
					}
				}

				return false
			},
		}
	}

	return finalDelegatee
}
