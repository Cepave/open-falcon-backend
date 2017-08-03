package rdb

import (
	gormExt "github.com/Cepave/open-falcon-backend/common/gorm"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
)

func GetAgentConfig(key string) *model.AgentConfigResult {
	var result model.AgentConfigResult
	gormDbExt := gormExt.ToDefaultGormDbExt(
		DbFacade.GormDb.First(&result, "key = ?", key),
	)

	if gormDbExt.IsRecordNotFound() {
		return nil
	}
	gormDbExt.PanicIfError()

	return &result
}
