package owl

import (
	"fmt"

	gormExt "github.com/Cepave/open-falcon-backend/common/gorm"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
)

type ErrNotInSameHierarchy struct {
	message string
}

func (err ErrNotInSameHierarchy) Error() string {
	return err.message
}

type cityInfo struct {
	ProvinceId   int16  `db:"pv_id"`
	ProvinceName string `db:"pv_name"`
	CityId       int16  `db:"ct_id"`
	CityName     string `db:"ct_name"`
}

type city1view struct {
	Id       int16              `json:"id"`
	Name     string             `json:"name"`
	PostCode string             `json:"post_code"`
	Province *owlModel.Province `json:"province"`
}

// Checks if the hierarchy for province and city are in the same administrative region
func CheckHierarchyForCity(provinceId int16, cityId int16) error {
	if cityId == -1 {
		return nil
	}

	cityInfo := &cityInfo{}
	DbFacade.SqlxDbCtrl.Get(
		cityInfo,
		`
		SELECT pv_id, pv_name, ct_id, ct_name
		FROM owl_province
			INNER JOIN
			owl_city
			ON pv_id = ct_pv_id
				AND ct_id = ?
		`,
		cityId,
	)

	if cityInfo.ProvinceId != provinceId {
		return ErrNotInSameHierarchy{
			message: fmt.Sprintf(
				"City[ID: %d][Name: %s] belongs to province[ID: %d][Name: %s]. But got province[ID: %d]",
				cityInfo.CityId, cityInfo.CityName,
				cityInfo.ProvinceId, cityInfo.ProvinceName,
				provinceId,
			),
		}
	}

	return nil
}

func GetProvinceById(provinceId int16) *owlModel.Province {
	var result owlModel.Province

	gormDbExt := gormExt.ToDefaultGormDbExt(
		DbFacade.GormDb.First(&result, provinceId),
	)

	return gormDbExt.IfRecordNotFound(&result, (*owlModel.Province)(nil)).(*owlModel.Province)
}

func GetProvincesByName(name string) []*owlModel.Province {
	var results []*owlModel.Province

	var gormDbExt = gormExt.ToDefaultGormDbExt(
		DbFacade.GormDb.Model(&owlModel.Province{}).
			Where("pv_name LIKE ?", name+"%").
			Find(&results),
	)

	return gormDbExt.IfRecordNotFound(results, results).([]*owlModel.Province)
}

func GetCityById(cityId int16) *city1view {
	var result owlModel.City1
	gormDbExt := gormExt.ToDefaultGormDbExt(
		DbFacade.GormDb.Model(&owlModel.City1{}).
			Select(`
			ct_id, ct_name, ct_post_code, pv_id, pv_name
			`).
			Joins(`
			INNER JOIN
			owl_province
			ON ct_pv_id = pv_id
			`).
			First(&result, cityId),
	)

	view := city1view{
		Id:       result.Id,
		Name:     result.Name,
		PostCode: result.PostCode,
		Province: &owlModel.Province{
			Id:   result.ProvinceId,
			Name: result.ProvinceName,
		},
	}

	return gormDbExt.IfRecordNotFound(&view, (*city1view)(nil)).(*city1view)
}

func GetCity2ById(cityId int16) *owlModel.City2 {
	var result owlModel.City2

	gormDbExt := gormExt.ToDefaultGormDbExt(
		DbFacade.GormDb.First(&result, cityId),
	)

	return gormDbExt.IfRecordNotFound(&result, (*owlModel.City2)(nil)).(*owlModel.City2)
}

func GetCity2sByName(prefixName string) []*owlModel.City2 {
	var result []*owlModel.City2

	gormDbExt := gormExt.ToDefaultGormDbExt(
		DbFacade.GormDb.Model(&owlModel.City2{}).
			Where("ct_name LIKE ?", prefixName+"%").
			Find(&result),
	)

	return gormDbExt.IfRecordNotFound(result, result).([]*owlModel.City2)
}

func GetCitiesByName(name string) []*city1view {
	var q = DbFacade.GormDb.Model(&owlModel.City1{}).
		Select(`
		ct_id, ct_name, ct_post_code, pv_id, pv_name
		`).
		Joins(`
		INNER JOIN
		owl_province
		ON ct_pv_id = pv_id
		AND ct_name LIKE ?
		`,
			name+"%",
		)

	var results []*owlModel.City1
	gormExt.ToDefaultGormDbExt(q.Find(&results))

	var views = []*city1view{}
	for _, r := range results {
		v := &city1view{
			Id:       r.Id,
			Name:     r.Name,
			PostCode: r.PostCode,
			Province: &owlModel.Province{
				Id:   r.ProvinceId,
				Name: r.ProvinceName,
			},
		}
		views = append(views, v)
	}
	return views
}

func GetCitiesInProvinceByName(pvId int, name string) []*owlModel.City2 {
	var q = DbFacade.GormDb.Model(&owlModel.City2{}).
		Select("ct_id, ct_name, ct_post_code").
		Where(
			`
			ct_pv_id = ?
			AND ct_name LIKE ?
			`,
			pvId, name+"%",
		)

	var results []*owlModel.City2
	gormDbExt := gormExt.ToDefaultGormDbExt(q.Find(&results))

	return gormDbExt.IfRecordNotFound(
		results, results,
	).([]*owlModel.City2)
}
