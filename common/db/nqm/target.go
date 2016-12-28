package nqm

import (
	"fmt"
	"database/sql"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	tb "github.com/Cepave/open-falcon-backend/common/textbuilder"
	sqlb "github.com/Cepave/open-falcon-backend/common/textbuilder/sql"
	"github.com/Cepave/open-falcon-backend/common/utils"
	"github.com/jmoiron/sqlx"
	sqlxExt "github.com/Cepave/open-falcon-backend/common/db/sqlx"
	"github.com/jinzhu/gorm"
	gormExt "github.com/Cepave/open-falcon-backend/common/gorm"
)

type ErrDuplicatedNqmTarget struct {
	Host string
}
func (err ErrDuplicatedNqmTarget) Error() string {
	return fmt.Sprintf("Duplicated HOST: %s", err.Host)
}

// Add and retrieve detail data of target
//
// Errors:
// 		ErrDuplicatedNqmTarget - The target is existing with the same host
//		ErrNotInSameHierarchy - The city is not belonging to the province
func AddTarget(newTarget *nqmModel.TargetForAdding) (*nqmModel.Target, error) {
	/**
	 * Checks the hierarchy over administrative region
	 */
	err := owlDb.CheckHierarchyForCity(newTarget.ProvinceId, newTarget.CityId)
	if err != nil {
		return nil, err
	}
	// :~)

	/**
	 * Executes the insertion of target and its related data
	 */
	txProcessor := &addTargetTx{
		target: newTarget,
	}

	DbFacade.NewSqlxDbCtrl().InTx(txProcessor)
	// :~)

	if txProcessor.err != nil {
		return nil, txProcessor.err
	}

	return GetTargetById(newTarget.Id), nil
}
// Update and retrieve detail data of target
func UpdateTarget(oldTarget *nqmModel.Target, updatedTarget *nqmModel.TargetForAdding) (*nqmModel.Target, error) {
	/**
	 * Checks the hierarchy over administrative region
	 */
	err := owlDb.CheckHierarchyForCity(updatedTarget.ProvinceId, updatedTarget.CityId)
	if err != nil {
		return nil, err
	}
	// :~)

	txProcessor := &updateTargetTx{
		updatedTarget: updatedTarget,
		oldTarget: oldTarget.ToTargetForAdding(),
	}

	DbFacade.NewSqlxDbCtrl().InTx(txProcessor)

	return GetTargetById(oldTarget.Id), nil
}
func GetTargetById(targetId int32) *nqmModel.Target {
	dbListTargets := DbFacade.GormDb.Model(&nqmModel.Target{}).
		Select(`
			tg_id, tg_name, tg_host, tg_probed_by_all, tg_status, tg_available, tg_comment, tg_created_ts,
			isp_id, isp_name, pv_id, pv_name, ct_id, ct_name, nt_id, nt_value,
			GROUP_CONCAT(gt.gt_id ORDER BY gt_name ASC SEPARATOR ',') AS gt_ids,
			GROUP_CONCAT(gt.gt_name ORDER BY gt_name ASC SEPARATOR '\0') AS gt_names
		`).
		Joins(`
			INNER JOIN
			owl_isp AS isp
			ON tg_isp_id = isp.isp_id
			INNER JOIN
			owl_province AS pv
			ON tg_pv_id = pv.pv_id
			INNER JOIN
			owl_city AS ct
			ON tg_ct_id = ct.ct_id
			INNER JOIN
			owl_name_tag AS nt
			ON tg_nt_id = nt.nt_id
			LEFT OUTER JOIN
			nqm_target_group_tag AS tgt
			ON tg_id = tgt.tgt_tg_id
			LEFT OUTER JOIN
			owl_group_tag AS gt
			ON tgt.tgt_gt_id = gt.gt_id
		`).
		Where("tg_id = ?", targetId).
		Group(`
			tg_id, tg_name, tg_host, tg_probed_by_all, tg_status, tg_available, tg_comment, tg_created_ts,
			isp_id, isp_name, pv_id, pv_name, ct_id, ct_name, nt_id, nt_value
		`)

	var loadedTarget = &nqmModel.Target{}
	dbListTargets = dbListTargets.Scan(loadedTarget)
	if dbListTargets.Error == gorm.ErrRecordNotFound {
		return nil
	}
	gormExt.ToDefaultGormDbExt(dbListTargets).PanicIfError()

	loadedTarget.AfterLoad()
	return loadedTarget
}
func ListTargets(query *nqmModel.TargetQuery, paging commonModel.Paging) ([]*nqmModel.Target, *commonModel.Paging) {
	var result []*nqmModel.Target

	var funcTxLoader gormExt.TxCallbackFunc = func(txGormDb *gorm.DB) commonDb.TxFinale {
		/**
		 * Retrieves the page of data
		 */
		var dbListTargets = txGormDb.Model(&nqmModel.Target{}).
			Select(`SQL_CALC_FOUND_ROWS
				tg_id, tg_name, tg_host, tg_probed_by_all, tg_status, tg_available, tg_comment, tg_created_ts,
				isp_id, isp_name, pv_id, pv_name, ct_id, ct_name, nt_id, nt_value,
				COUNT(gt.gt_id) AS gt_number,
				GROUP_CONCAT(gt.gt_id ORDER BY gt_name ASC SEPARATOR ',') AS gt_ids,
				GROUP_CONCAT(gt.gt_name ORDER BY gt_name ASC SEPARATOR '\0') AS gt_names
			`).
			Joins(`
				INNER JOIN
				owl_isp AS isp
				ON tg_isp_id = isp.isp_id
				INNER JOIN
				owl_province AS pv
				ON tg_pv_id = pv.pv_id
				INNER JOIN
				owl_city AS ct
				ON tg_ct_id = ct.ct_id
				INNER JOIN
				owl_name_tag AS nt
				ON tg_nt_id = nt.nt_id
				LEFT OUTER JOIN
				nqm_target_group_tag AS tgt
				ON tg_id = tgt.tgt_tg_id
				LEFT OUTER JOIN
				owl_group_tag AS gt
				ON tgt.tgt_gt_id = gt.gt_id
			`).
			Limit(paging.Size).
			Group(`
				tg_id, tg_name, tg_host, tg_probed_by_all, tg_status, tg_available, tg_comment, tg_created_ts,
				isp_id, isp_name, pv_id, pv_name, ct_id, ct_name, nt_id, nt_value
			`).
			Order(buildSortingClauseOfTargets(&paging)).
			Offset(paging.GetOffset())

		if query.Name != "" {
			dbListTargets = dbListTargets.Where("tg_name LIKE ?", query.Name + "%")
		}
		if query.Host != "" {
			dbListTargets = dbListTargets.Where("tg_host LIKE ?", query.Host + "%")
		}
		if query.HasIspId {
			dbListTargets = dbListTargets.Where("tg_isp_id = ?", query.IspId)
		}
		if query.HasStatusCondition {
			dbListTargets = dbListTargets.Where("tg_status = ?", query.Status)
		}
		// :~)

		gormExt.ToDefaultGormDbExt(dbListTargets.Find(&result)).PanicIfError()

		return commonDb.TxCommit
	}

	gormExt.ToDefaultGormDbExt(DbFacade.GormDb).SelectWithFoundRows(
		funcTxLoader, &paging,
	)

	/**
	 * Loads group tags
	 */
	for _, target := range result {
		target.AfterLoad()
	}
	// :~)

	return result, &paging
}

// Gets the target object or nil if the id is not existing
func GetSimpleTarget1ById(targetId int32) *nqmModel.SimpleTarget1 {
	var result nqmModel.SimpleTarget1

	if !DbFacade.SqlxDbCtrl.GetOrNoRow(
		&result,
		`
		SELECT tg_id, tg_host, tg_name,
			tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id
		FROM nqm_target
		WHERE tg_id = ?
		`,
		targetId,
	) {
		return nil
	}

	return &result
}
// Gets the targets by filter
func LoadSimpleTarget1sByFilter(filter *nqmModel.TargetFilter) []*nqmModel.SimpleTarget1 {
	var result []*nqmModel.SimpleTarget1

	/**
	 * Builds ( <expr> OR <expr> OR ... ) of SQL
	 */
	var buildRepeatOr = func(syntax string, arrayObject interface{}) tb.TextGetter {
		return tb.Surrounding(
			t.S("( "),
			tb.RepeatAndJoinByLen(t.S(syntax), sqlb.C["or"], arrayObject),
			t.S(" )"),
		)
	}
	// :~)

	/**
	 * Processes the arguments used in query
	 */
	var sqlArgs []interface{}
	sqlArgs = utils.AppendToAny(sqlArgs, filter.Name)
	sqlArgs = utils.AppendToAny(sqlArgs, filter.Host)

	sqlArgs = utils.MakeAbstractArray(sqlArgs).
		MapTo(
			utils.TypedFuncToMapper(func(v string) string {
				return v + "%"
			}),
			utils.TypeOfString,
		).GetArrayOfAny()
	// :~)

	DbFacade.SqlxDbCtrl.Select(
		&result,
		`
		SELECT tg_id, tg_name, tg_host,
			tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id
		FROM nqm_target
		` +
		sqlb.Where(
			sqlb.And(
				buildRepeatOr("tg_name LIKE ?", filter.Name),
				buildRepeatOr("tg_host LIKE ?", filter.Host),
			),
		).String(),
		sqlArgs...,
	)

	return result
}

var orderByDialectForTagets = commonModel.NewSqlOrderByDialect(
	map[string]string {
		"name": "tg_name",
		"host": "tg_host",
		"isp": "isp_name",
		"province": "pv_name",
		"city": "ct_name",
		"status": "tg_status",
		"available": "tg_available",
		"comment": "tg_comment",
		"creation_time": "tg_created_ts",
		"name_tag": "nt_value",
	},
)

func buildSortingClauseOfTargets(paging *commonModel.Paging) string {
	if len(paging.OrderBy) == 0 {
		paging.OrderBy = append(paging.OrderBy, &commonModel.OrderByEntity{ "status", commonModel.Descending })
		paging.OrderBy = append(paging.OrderBy, &commonModel.OrderByEntity{ "available", commonModel.Descending })
	}

	if len(paging.OrderBy) == 1 {
		switch paging.OrderBy[0].Expr {
		case "province":
			paging.OrderBy = append(paging.OrderBy, &commonModel.OrderByEntity{ "city", commonModel.Ascending })
		}
	}

	if paging.OrderBy[len(paging.OrderBy) - 1].Expr != "creation_time" {
		paging.OrderBy = append(paging.OrderBy, &commonModel.OrderByEntity{ "creation_time", commonModel.Descending })
	}

	querySyntax, err := orderByDialectForTagets.ToQuerySyntax(paging.OrderBy)
	gormExt.DefaultGormErrorConverter.PanicIfError(err)

	return querySyntax
}

func init() {
	originFunc := orderByDialectForTagets.FuncEntityToSyntax
	orderByDialectForTagets.FuncEntityToSyntax = func(entity *commonModel.OrderByEntity) (string, error) {
		switch entity.Expr {
		case "group_tag":
			return owlDb.GetSyntaxOfOrderByGroupTags(entity), nil
		}

		return originFunc(entity)
	}
}

type addTargetTx struct {
	target *nqmModel.TargetForAdding
	err error
}
func (targetTx *addTargetTx) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	targetTx.target.NameTagId = owlDb.BuildAndGetNameTagId(
		tx, targetTx.target.NameTagValue,
	)

	targetTx.addTarget(tx)
	if targetTx.err != nil {
		return commonDb.TxRollback
	}

	targetTx.prepareGroupTags(tx)
	return commonDb.TxCommit
}
func (targetTx *addTargetTx) addTarget(tx *sqlx.Tx) {
	txExt := sqlxExt.ToTxExt(tx)
	newTarget := targetTx.target

	addedNqmTarget := txExt.NamedExec(
		`
		INSERT INTO nqm_target(
			tg_name, tg_host, tg_comment,
			tg_probed_by_all,
			tg_status, tg_available,
			tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id
		)
		SELECT :name, :host, :comment,
			:probed_by_all,
			:status, true,
			:isp_id, :province_id, :city_id, :name_tag_id
		FROM DUAL
		WHERE NOT EXISTS (
			SELECT *
			FROM nqm_target
			WHERE tg_host = :host
		)
		`,
		map[string]interface{} {
			"name" : newTarget.Name,
			"host" : newTarget.Host,
			"status" : newTarget.Status,
			"probed_by_all" : newTarget.ProbedByAll,
			"isp_id" : newTarget.IspId,
			"province_id" : newTarget.ProvinceId,
			"city_id" : newTarget.CityId,
			"name_tag_id" : newTarget.NameTagId,
			"comment" : sql.NullString {
				newTarget.Comment,
				newTarget.Comment != "",
			},
		},
	)

	/**
	 * Rollback if the NQM target is existing(duplicated by host)
	 */
	if commonDb.ToResultExt(addedNqmTarget).RowsAffected() == 0 {
		targetTx.err = ErrDuplicatedNqmTarget{ newTarget.Host }
		return
	}
	// :~)

	txExt.Get(
		&newTarget.Id,
		`
		SELECT tg_id FROM nqm_target
		WHERE tg_host = ?
		`,
		newTarget.Host,
	)
}
func (targetTx *addTargetTx) prepareGroupTags(tx *sqlx.Tx) {
	newTarget := targetTx.target
	buildGroupTagsForTarget(
		tx, newTarget.Id, newTarget.GroupTags,
	)
}

func buildGroupTagsForTarget(tx *sqlx.Tx, targetId int32, groupTags []string) {
	owlDb.BuildGroupTags(
		tx, groupTags,
		func(tx *sqlx.Tx, groupTag string) {
			tx.MustExec(
				`
				INSERT INTO nqm_target_group_tag(tgt_tg_id, tgt_gt_id)
				VALUES(
					?,
					(
						SELECT gt_id
						FROM owl_group_tag
						WHERE gt_name = ?
					)
				)
				`,
				targetId,
				groupTag,
			)
		},
	)
}

type updateTargetTx struct {
	updatedTarget *nqmModel.TargetForAdding
	oldTarget *nqmModel.TargetForAdding
}
func (agentTx *updateTargetTx) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	agentTx.loadNameTagId(tx)

	updatedTarget, oldTarget := agentTx.updatedTarget, agentTx.oldTarget
	tx.MustExec(
		`
		UPDATE nqm_target
		SET tg_name = ?,
			tg_comment = ?,
			tg_status = ?,
			tg_isp_id = ?,
			tg_pv_id = ?,
			tg_ct_id = ?,
			tg_nt_id = ?
		WHERE tg_id = ?
		`,
		updatedTarget.Name,
		sql.NullString{ updatedTarget.Comment, updatedTarget.Comment != "" },
		updatedTarget.Status,
		updatedTarget.IspId,
		updatedTarget.ProvinceId,
		updatedTarget.CityId,
		updatedTarget.NameTagId,
		oldTarget.Id,
	)

	agentTx.updateGroupTags(tx)
	return commonDb.TxCommit
}
func (agentTx *updateTargetTx) loadNameTagId(tx *sqlx.Tx) {
	updatedTarget, oldTarget := agentTx.updatedTarget, agentTx.oldTarget

	if updatedTarget.NameTagValue == oldTarget.NameTagValue {
		return
	}

	updatedTarget.NameTagId = owlDb.BuildAndGetNameTagId(
		tx, updatedTarget.NameTagValue,
	)
}
func (agentTx *updateTargetTx) updateGroupTags(tx *sqlx.Tx) {
	updatedTarget, oldTarget := agentTx.updatedTarget, agentTx.oldTarget
	if updatedTarget.AreGroupTagsSame(oldTarget) {
		return
	}

	tx.MustExec(
		`
		DELETE FROM nqm_target_group_tag
		WHERE tgt_tg_id = ?
		`,
		oldTarget.Id,
	)

	buildGroupTagsForTarget(
		tx, oldTarget.Id, updatedTarget.GroupTags,
	)
}
