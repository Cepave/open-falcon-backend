package owl

import (
	gormExt "github.com/Cepave/open-falcon-backend/common/gorm"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
)

func GetISPsByName(name string) []string {
	var result []*owlModel.Isp
	var queryResults = DbFacade.GormDb.Model(&owlModel.Isp{}).
		Select(`
		isp_name
	`).
		Where(`
		isp_name LIKE ?
		`,
		name+"%",
	)
	gormExt.ToDefaultGormDbExt(queryResults.Find(&result))

	var owlIspNames = []string{}
	for _, v := range result {
		owlIspNames = append(owlIspNames, v.Name)
	}

	return owlIspNames
}
