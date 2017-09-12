package service

import (
	"net/http"

	gt "gopkg.in/h2non/gentleman.v2"

	oHttp "github.com/Cepave/open-falcon-backend/common/http"
	"github.com/Cepave/open-falcon-backend/common/http/client"
	"github.com/Cepave/open-falcon-backend/common/model"
	apiModel "github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	"github.com/juju/errors"
)

type MysqlApiServiceConfig struct {
	*oHttp.RestfulClientConfig
}

type MysqlApiService interface {
	GetHealth() *model.MysqlApi
}

func NewMysqlApiService(config MysqlApiServiceConfig) MysqlApiService {
	newClient := oHttp.NewApiService(config.RestfulClientConfig).NewClient()
	return &mysqlApiServiceImpl{
		getHealth: newClient.Get().AddPath("/health"),
	}
}

type mysqlApiServiceImpl struct {
	getHealth *gt.Request
}

func (s *mysqlApiServiceImpl) GetHealth() *model.MysqlApi {
	req := s.getHealth.Clone()
	resp, err := client.ToGentlemanReq(req).SendAndMatch(
		func(resp *gt.Response) error {
			switch resp.StatusCode {
			case http.StatusOK, http.StatusNotFound:
				return nil
			}

			return errors.Errorf(client.ToGentlemanResp(resp).ToDetailString())
		},
	)

	view := &model.MysqlApi{
		Address:    req.Context.Request.URL.String(),
		StatusCode: resp.StatusCode,
	}

	health := &apiModel.HealthView{}
	if e := client.ToGentlemanResp(resp).BindJson(health); e != nil {
		err = errors.Annotate(err, e.Error())
	}
	if err != nil {
		view.Message = err.Error()
	}
	view.Response = health

	return view
}
