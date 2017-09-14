package model

import (
	apiModel "github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
)

type MysqlApi struct {
	Address  string               `json:"address"`
	Message  string               `json:"message"`
	Response *apiModel.HealthView `json:"response"`
}
