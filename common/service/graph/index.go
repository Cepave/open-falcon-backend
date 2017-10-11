package graph

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	ojson "github.com/Cepave/open-falcon-backend/common/json"
	gt "gopkg.in/h2non/gentleman.v2"

	oHttp "github.com/Cepave/open-falcon-backend/common/http"
	"github.com/Cepave/open-falcon-backend/common/http/client"
)

// Configurations for constructing "GraphService"
type GraphServiceConfig struct {
	*oHttp.RestfulClientConfig
}

type VacuumIndexConfig struct {
	BeforeDays int
}

type ResultOfVacuumIndex struct {
	BeforeTime   ojson.JsonTime `json:"before_time"`
	AffectedRows *struct {
		Endpoints uint64 `json:"endpoints"`
		Tags      uint64 `json:"tags"`
		Counters  uint64 `json:counters`
	} `json:"affected_rows"`
}

func (self *ResultOfVacuumIndex) GetBeforeTime() time.Time {
	return time.Time(self.BeforeTime)
}
func (self *ResultOfVacuumIndex) String() string {
	return fmt.Sprintf(
		"Before Time: [%s]. Vacuumed rows: endpoints[%d]; Tags[%d]; Counters[%d].",
		self.GetBeforeTime(),
		self.AffectedRows.Endpoints,
		self.AffectedRows.Tags,
		self.AffectedRows.Counters,
	)
}

// Main service for calling of RESTful service on MySqlApi
type GraphService interface {
	// Vacuums index of graph
	VacuumIndex(*VacuumIndexConfig) *ResultOfVacuumIndex
}

func NewGraphService(config *GraphServiceConfig) GraphService {
	newClient := oHttp.NewApiService(config.RestfulClientConfig).NewClient()

	return &graphServiceImpl{
		vacuumIndex: newClient.Post().AddPath("/api/v1/graph/endpoint-index/vacuum"),
	}
}

type graphServiceImpl struct {
	vacuumIndex *gt.Request
}

func (self *graphServiceImpl) VacuumIndex(config *VacuumIndexConfig) *ResultOfVacuumIndex {
	req := self.vacuumIndex.Clone().SetQuery("for_days", strconv.Itoa(config.BeforeDays))

	result := &ResultOfVacuumIndex{}
	client.ToGentlemanResp(
		client.ToGentlemanReq(req).SendAndStatusMustMatch(http.StatusOK),
	).MustBindJson(result)

	return result
}
