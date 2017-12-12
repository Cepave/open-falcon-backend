package owl

import (
	"net/http"
	"strconv"
	"time"

	gt "gopkg.in/h2non/gentleman.v2"

	oHttp "github.com/Cepave/open-falcon-backend/common/http"
	"github.com/Cepave/open-falcon-backend/common/http/client"
	"github.com/Cepave/open-falcon-backend/common/json"
)

type ClearLogServiceConfig struct {
	*oHttp.RestfulClientConfig
}

type ClearLogService interface {
	ClearLogEntries(int) *ResultOfClearLogEntries
}

func NewClearLogService(config ClearLogServiceConfig) ClearLogService {
	newClient := oHttp.NewApiService(config.RestfulClientConfig).NewClient()

	return &clearLogServiceImpl{
		clearLogEntries: newClient.Post().AddPath("/api/v1/owl/task/log/clear"),
	}
}

type clearLogServiceImpl struct {
	clearLogEntries *gt.Request
}

type ResultOfClearLogEntries struct {
	BeforeTime   json.JsonTime `json:"before_time"`
	AffectedRows int           `json:"affected_rows"`
}

func (r *ResultOfClearLogEntries) GetBeforeTime() time.Time {
	return time.Time(r.BeforeTime)
}

// ClearLogEntries panics if error happens.
func (s *clearLogServiceImpl) ClearLogEntries(forDays int) *ResultOfClearLogEntries {
	req := s.clearLogEntries.Clone().
		SetQueryParams(map[string]string{
			"for_days": strconv.Itoa(forDays),
		})

	result := &ResultOfClearLogEntries{}

	client.ToGentlemanResp(
		client.ToGentlemanReq(req).SendAndStatusMustMatch(http.StatusOK),
	).MustBindJson(result)

	return result
}
