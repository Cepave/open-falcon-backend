package service

import (
	gt "gopkg.in/h2non/gentleman.v2"

	oHttp "github.com/Cepave/open-falcon-backend/common/http"
	"github.com/Cepave/open-falcon-backend/common/http/client"
	"github.com/Cepave/open-falcon-backend/common/model"
	apiModel "github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
)

type MysqlApiServiceConfig struct {
	*oHttp.RestfulClientConfig
}

type MysqlApiService interface {
	GetHealth() *model.MysqlApi
}

func NewMysqlApiService(config *MysqlApiServiceConfig) MysqlApiService {
	newClient := oHttp.NewApiService(config.RestfulClientConfig).NewClient()
	return &mysqlApiServiceImpl{
		Config:    config,
		getHealth: newClient.Get().AddPath("/health"),
	}
}

type mysqlApiServiceImpl struct {
	Config    *MysqlApiServiceConfig
	getHealth *gt.Request
}

func (s *mysqlApiServiceImpl) GetHealth() *model.MysqlApi {
	req := s.getHealth.Clone()
	view := &model.MysqlApi{
		Address: s.Config.HttpClientConfig.Url,
	}

	resp, err := client.ToGentlemanReq(req).SendAndStatusMatch(200)
	if err != nil {
		view.Message = err.Error()
		return view
	}

	// Handle the response body
	health := &apiModel.HealthView{}
	if err := client.ToGentlemanResp(resp).BindJson(health); err != nil {
		view.Message = err.Error()
		return view
	}
	view.Response = health

	return view
}
