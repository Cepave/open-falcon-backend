package model

import (
	apiModel "github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
)

type MysqlApi struct {
	// Address of MySQL-API
	Address string `json:"address"`
	// Error messages during the health request
	Message string `json:"message"`
	// Response of the health request
	Response *apiModel.HealthView `json:"response"`
}
