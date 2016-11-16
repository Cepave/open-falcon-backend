package owl

import (
	gormExt "github.com/Cepave/open-falcon-backend/common/gorm"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
)

func GetISPsByName(name string) []*owlModel.Isp {
	var q = DbFacade.GormDb.Model(&owlModel.Isp{}).
		Select(`
		*
	`).
		Where(`
		isp_name LIKE ?
		`,
		name+"%",
	)

	var results []*owlModel.Isp
	gormExt.ToDefaultGormDbExt(q.Find(&results))

	return results
}
