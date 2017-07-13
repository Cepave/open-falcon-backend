package owl

import (
	gormExt "github.com/Cepave/open-falcon-backend/common/gorm"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
)

func GetISPsByName(name string) []*owlModel.Isp {
	var q = DbFacade.GormDb.Model(&owlModel.Isp{}).
		Select(`*`).
		Where(`isp_name LIKE ?`, name+"%")

	var results []*owlModel.Isp

	gormDbExt := gormExt.ToDefaultGormDbExt(q.Find(&results))
	if gormDbExt.IsRecordNotFound() {
		return results
	}
	gormDbExt.PanicIfError()

	return results
}

func GetIspById(ispId int16) *owlModel.Isp {
	var result owlModel.Isp
	gormDbExt := gormExt.ToDefaultGormDbExt(
		DbFacade.GormDb.First(&result, ispId),
	)

	if gormDbExt.IsRecordNotFound() {
		return nil
	}
	gormDbExt.PanicIfError()

	return &result
}
