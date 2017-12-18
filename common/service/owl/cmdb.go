package owl

import (
	"net/http"

	"github.com/satori/go.uuid"
	gt "gopkg.in/h2non/gentleman.v2"

	oHttp "github.com/Cepave/open-falcon-backend/common/http"
	"github.com/Cepave/open-falcon-backend/common/http/client"
	oJson "github.com/Cepave/open-falcon-backend/common/json"
	model "github.com/Cepave/open-falcon-backend/common/model/owl"
	"github.com/juju/errors"
)

type CmdbServiceConfig struct {
	*oHttp.RestfulClientConfig
}

type CmdbService interface {
	StartSyncJob() (uuid.UUID, error)
	GetJobStatus(uuid.UUID) *model.SyncCmdbJobInfo
}

func NewCmdbService(config CmdbServiceConfig) CmdbService {
	newClient := oHttp.NewApiService(config.RestfulClientConfig).NewClient()

	return &cmdbServiceImpl{
		startSyncJob: newClient.Post().AddPath("/api/v1/cmdb/sync"),
		getJobStatus: newClient.Get().AddPath("/api/v1/cmdb/sync"),
	}
}

type cmdbServiceImpl struct {
	startSyncJob *gt.Request
	getJobStatus *gt.Request
}

func (s *cmdbServiceImpl) StartSyncJob() (uuid uuid.UUID, err error) {
	req := s.startSyncJob.Clone()

	client.ToGentlemanReq(req).SendAndMustMatch(
		func(resp *gt.Response) error {
			switch resp.StatusCode {
			case http.StatusOK:
				var jsonResult = &struct {
					Uuid oJson.Uuid `json:"sync_id"`
				}{}

				client.ToGentlemanResp(resp).MustBindJson(jsonResult)
				uuid = jsonResult.Uuid.ToUuid()

				return nil
			case http.StatusBadRequest:
				var jsonResult = &struct {
					ErrorCode    int        `json:"error_code"`
					ErrorMessage string     `json:"error_message"`
					LastSyncId   oJson.Uuid `json:"last_sync_id"`
				}{}

				client.ToGentlemanResp(resp).MustBindJson(jsonResult)

				if jsonResult.ErrorCode == 1 {
					err = &model.UnableToLockCmdbSyncJob{
						Uuid: jsonResult.LastSyncId,
					}

					return nil
				}

				return errors.Errorf("Bad Request. Error Code: [%d]. Message: [%s]", jsonResult.ErrorCode, jsonResult.ErrorMessage)
			}

			return errors.Errorf(client.ToGentlemanResp(resp).ToDetailString())
		},
	)

	return
}

func (s *cmdbServiceImpl) GetJobStatus(uuid uuid.UUID) *model.SyncCmdbJobInfo {
	req := s.getJobStatus.Clone().
		AddPath("/" + uuid.String())

	resp := client.ToGentlemanReq(req).SendAndStatusMustMatch(http.StatusOK)

	jobInfo := &model.SyncCmdbJobInfo{}
	client.ToGentlemanResp(resp).MustBindJson(jobInfo)

	return jobInfo
}
