package nqm

// Defines the interface for process a instance of metrics
type MetricFilter interface {
	// Checks whether or not the metrics matches conditions of filter
	IsMatch(model *Metrics) bool
}

/**
 * Macro-struct re-used by various data
 */
type Metrics struct {
	Max                     int16   `json:"max"`
	Min                     int16   `json:"min"`
	Avg                     float64 `json:"avg"`
	Med                     int16   `json:"med"`
	Mdev                    float64 `json:"mdev"`
	Loss                    float64 `json:"loss"`
	Count                   int32   `json:"count"`
	NumberOfSentPackets     uint64  `json:"number_of_sent_packets"`
	NumberOfReceivedPackets uint64  `json:"number_of_received_packets"`
	NumberOfAgents          int32   `json:"number_of_agents"`
	NumberOfTargets         int32   `json:"number_of_targets"`
}
