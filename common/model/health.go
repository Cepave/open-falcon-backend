package model

import (
	apiModel "github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
)

type MysqlApi struct {
	Address    string               `json:"address"`
	StatusCode int                  `json:"status_code"`
	Message    string               `json:"message"`
	Response   *apiModel.HealthView `json:"response"`
}
