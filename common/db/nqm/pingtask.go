package nqm

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"

	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	gormExt "github.com/Cepave/open-falcon-backend/common/gorm"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	tb "github.com/Cepave/open-falcon-backend/common/textbuilder"
	sqlb "github.com/Cepave/open-falcon-backend/common/textbuilder/sql"
)

type addAgentPingtaskTx struct {
	agentPingtask *nqmModel.AgentPingtask
	err           error
}

func (agentPingtaskTx *addAgentPingtaskTx) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	tx.MustExec(
		`
		INSERT INTO nqm_agent_ping_task(apt_ag_id,apt_pt_id)
		VALUES
		(?,?)
		ON DUPLICATE KEY UPDATE
		apt_ag_id=VALUES(apt_ag_id),
		apt_pt_id=VALUES(apt_pt_id)
		`,
		agentPingtaskTx.agentPingtask.AgentID,
		agentPingtaskTx.agentPingtask.PingtaskID,
	)
	if agentPingtaskTx.err != nil {
		return commonDb.TxRollback
	}

	return commonDb.TxCommit
}

func AssignPingtaskToAgentForAgent(aID int32, pID int32) (*nqmModel.Agent, error) {
	txProcessor := &addAgentPingtaskTx{
		agentPingtask: &nqmModel.AgentPingtask{AgentID: aID, PingtaskID: pID},
	}

	DbFacade.NewSqlxDbCtrl().InTx(txProcessor)
	// :~)

	if txProcessor.err != nil {
		return nil, txProcessor.err
	}

	return GetAgentById(aID), nil
}

type deleteAgentPingtaskTx struct {
	agentPingtask *nqmModel.AgentPingtask
	err           error
}

func (agentPingtaskTx *deleteAgentPingtaskTx) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	tx.MustExec(
		`
		DELETE from nqm_agent_ping_task WHERE
		apt_ag_id=? AND apt_pt_id=? LIMIT 1;
		`,
		agentPingtaskTx.agentPingtask.AgentID,
		agentPingtaskTx.agentPingtask.PingtaskID,
	)

	if agentPingtaskTx.err != nil {
		return commonDb.TxRollback
	}

	return commonDb.TxCommit
}

func RemovePingtaskFromAgentForAgent(aID int32, pID int32) (*nqmModel.Agent, error) {
	txProcessor := &deleteAgentPingtaskTx{
		agentPingtask: &nqmModel.AgentPingtask{AgentID: aID, PingtaskID: pID},
	}

	DbFacade.NewSqlxDbCtrl().InTx(txProcessor)
	// :~)

	if txProcessor.err != nil {
		return nil, txProcessor.err
	}

	return GetAgentById(aID), nil
}

func AssignPingtaskToAgentForPingtask(aID int32, pID int32) (*nqmModel.PingtaskView, error) {
	txProcessor := &addAgentPingtaskTx{
		agentPingtask: &nqmModel.AgentPingtask{AgentID: aID, PingtaskID: pID},
	}

	DbFacade.NewSqlxDbCtrl().InTx(txProcessor)
	// :~)

	if txProcessor.err != nil {
		return nil, txProcessor.err
	}

	return GetPingtaskById(pID), nil
}

func RemovePingtaskFromAgentForPingtask(aID int32, pID int32) (*nqmModel.PingtaskView, error) {
	txProcessor := &deleteAgentPingtaskTx{
		agentPingtask: &nqmModel.AgentPingtask{AgentID: aID, PingtaskID: pID},
	}

	DbFacade.NewSqlxDbCtrl().InTx(txProcessor)
	// :~)

	if txProcessor.err != nil {
		return nil, txProcessor.err
	}

	return GetPingtaskById(pID), nil
}

var orderByDialectForPingtasks = commonModel.NewSqlOrderByDialect(
	map[string]string{
		"id":                    "pt_id",
		"period":                "pt_period",
		"name":                  "pt_name",
		"enable":                "pt_enable",
		"comment":               "pt_comment",
		"num_of_enabled_agents": "pt_num_of_enabled_agents",
	},
)

func buildSortingClauseOfPingtasks(paging *commonModel.Paging) string {
	querySyntax, err := orderByDialectForPingtasks.ToQuerySyntax(paging.OrderBy)
	gormExt.DefaultGormErrorConverter.PanicIfError(err)

	return querySyntax
}

// Lists the pingtasks according to the query parameters
func ListPingtasks(query *nqmModel.PingtaskQuery, paging commonModel.Paging) ([]*nqmModel.PingtaskView, *commonModel.Paging) {
	var result []*nqmModel.PingtaskView

	var funcTxLoader gormExt.TxCallbackFunc = func(txGormDb *gorm.DB) commonDb.TxFinale {
		sqlStr := `SELECT SQL_CALC_FOUND_ROWS
			pt_id, pt_period, pt_name, pt_enable, pt_comment,
			COUNT(DISTINCT ag.ag_id) AS pt_num_of_enabled_agents,
			GROUP_CONCAT(DISTINCT isp.isp_id ORDER BY isp_id ASC SEPARATOR ',') AS pt_isp_filter_ids,
			GROUP_CONCAT(DISTINCT isp.isp_name ORDER BY isp_id ASC SEPARATOR '\0') AS pt_isp_filter_names,
			GROUP_CONCAT(DISTINCT pv.pv_id ORDER BY pv_id ASC SEPARATOR ',') AS pt_province_filter_ids,
			GROUP_CONCAT(DISTINCT pv.pv_name ORDER BY pv_id ASC SEPARATOR '\0') AS pt_province_filter_names,
			GROUP_CONCAT(DISTINCT ct.ct_id ORDER BY ct_id ASC SEPARATOR ',') AS pt_city_filter_ids,
			GROUP_CONCAT(DISTINCT ct.ct_pv_id ORDER BY ct_id ASC SEPARATOR ',') AS pt_city_filter_pv_ids,
			GROUP_CONCAT(DISTINCT ct.ct_name ORDER BY ct_id ASC SEPARATOR '\0') AS pt_city_filter_names,
			GROUP_CONCAT(DISTINCT nt.nt_id ORDER BY nt_id ASC SEPARATOR ',') AS pt_name_tag_filter_ids,
			GROUP_CONCAT(DISTINCT nt.nt_value ORDER BY nt_id ASC SEPARATOR '\0') AS pt_name_tag_filter_values,
			GROUP_CONCAT(DISTINCT gt.gt_id ORDER BY gt_id ASC SEPARATOR ',') AS pt_group_tag_filter_ids,
			GROUP_CONCAT(DISTINCT gt.gt_name ORDER BY gt_id ASC SEPARATOR '\0') AS pt_group_tag_filter_names
			From nqm_ping_task
			LEFT JOIN
			nqm_agent_ping_task AS apt
			ON pt_id = apt.apt_pt_id
			LEFT JOIN
			nqm_agent AS ag
			ON apt.apt_ag_id = ag.ag_id AND ag.ag_status=true
			LEFT JOIN
			nqm_pt_target_filter_isp AS tfisp
			ON pt_id = tfisp.tfisp_pt_id
			LEFT JOIN
			owl_isp AS isp
			ON tfisp.tfisp_isp_id = isp.isp_id
			LEFT JOIN
			nqm_pt_target_filter_province AS tfpv
			ON pt_id = tfpv.tfpv_pt_id
			LEFT JOIN
			owl_province AS pv
			ON tfpv.tfpv_pv_id = pv.pv_id
			LEFT JOIN
			nqm_pt_target_filter_city AS tfct
			ON pt_id = tfct.tfct_pt_id
			LEFT JOIN
			owl_city AS ct
			ON tfct.tfct_ct_id = ct.ct_id
			LEFT JOIN
			nqm_pt_target_filter_name_tag AS tfnt
			ON pt_id = tfnt.tfnt_pt_id
			LEFT JOIN
			owl_name_tag AS nt
			ON tfnt.tfnt_nt_id = nt.nt_id
			LEFT JOIN
			nqm_pt_target_filter_group_tag AS tfgt
			ON pt_id = tfgt.tfgt_pt_id
			LEFT JOIN
			owl_group_tag AS gt
			ON tfgt.tfgt_gt_id = gt.gt_id
			%s
			GROUP BY pt_id, pt_period, pt_name, pt_enable, pt_comment
			ORDER BY %s
			Limit %d, %d
		`

		var conditions []tb.TextGetter
		var sqlParams []interface{}

		if query.Period != "" {
			conditions = append(conditions, tb.Dsl.S("pt_period = ?"))
			sqlParams = append(sqlParams, query.Period)
		}
		if query.Name != "" {
			conditions = append(conditions, tb.Dsl.S("pt_name LIKE ?"))
			sqlParams = append(sqlParams, query.Name+"%")
		}
		if ena, err := strconv.ParseBool(query.Enable); query.Enable != "" && err == nil {
			conditions = append(conditions, tb.Dsl.S("pt_enable = ?"))
			sqlParams = append(sqlParams, ena)
		}
		if query.Comment != "" {
			conditions = append(conditions, tb.Dsl.S("pt_comment LIKE ?"))
			sqlParams = append(sqlParams, query.Comment+"%")
		}
		if query.NumOfEnabledAgents != "" {
			conditions = append(conditions, tb.Dsl.S("pt_num_of_enabled_agents = ?"))
			sqlParams = append(sqlParams, query.NumOfEnabledAgents)
		}

		whereClause := sqlb.Where(sqlb.And(conditions...))

		sqlStr = fmt.Sprintf(sqlStr, whereClause.String(), buildSortingClauseOfPingtasks(&paging), paging.GetOffset(), paging.Size)

		var selectPingtask = txGormDb.Model(&nqmModel.PingtaskView{}).Raw(sqlStr, sqlParams...)
		gormExt.ToDefaultGormDbExt(selectPingtask.Find(&result)).PanicIfError()
		return commonDb.TxCommit
	}

	gormExt.ToDefaultGormDbExt(DbFacade.GormDb).SelectWithFoundRows(
		funcTxLoader, &paging,
	)

	/**
	 * Loads group tags
	 */
	for _, pingtask := range result {
		pingtask.AfterLoad()
	}
	// :~)

	return result, &paging
}

func GetPingtaskById(id int32) *nqmModel.PingtaskView {
	var selectPingtask = DbFacade.GormDb.Model(&nqmModel.PingtaskView{}).
		Select(`
			pt_id, pt_period, pt_name, pt_enable, pt_comment,
			COUNT(DISTINCT ag.ag_id) AS pt_num_of_enabled_agents,
			GROUP_CONCAT(DISTINCT isp.isp_id ORDER BY isp_id ASC SEPARATOR ',') AS pt_isp_filter_ids,
			GROUP_CONCAT(DISTINCT isp.isp_name ORDER BY isp_id ASC SEPARATOR '\0') AS pt_isp_filter_names,
			GROUP_CONCAT(DISTINCT pv.pv_id ORDER BY pv_id ASC SEPARATOR ',') AS pt_province_filter_ids,
			GROUP_CONCAT(DISTINCT pv.pv_name ORDER BY pv_id ASC SEPARATOR '\0') AS pt_province_filter_names,
			GROUP_CONCAT(DISTINCT ct.ct_id ORDER BY ct_id ASC SEPARATOR ',') AS pt_city_filter_ids,
			GROUP_CONCAT(DISTINCT ct.ct_name ORDER BY ct_id ASC SEPARATOR '\0') AS pt_city_filter_names,
			GROUP_CONCAT(DISTINCT nt.nt_id ORDER BY nt_id ASC SEPARATOR ',') AS pt_name_tag_filter_ids,
			GROUP_CONCAT(DISTINCT nt.nt_value ORDER BY nt_id ASC SEPARATOR '\0') AS pt_name_tag_filter_values,
			GROUP_CONCAT(DISTINCT gt.gt_id ORDER BY gt_id ASC SEPARATOR ',') AS pt_group_tag_filter_ids,
			GROUP_CONCAT(DISTINCT gt.gt_name ORDER BY gt_id ASC SEPARATOR '\0') AS pt_group_tag_filter_names
		`).
		Joins(`
			LEFT JOIN
			nqm_agent_ping_task AS apt
			ON pt_id = apt.apt_pt_id
			LEFT JOIN
			nqm_agent AS ag
			ON apt.apt_ag_id = ag.ag_id AND ag.ag_status=true
			LEFT JOIN
			nqm_pt_target_filter_isp AS tfisp
			ON pt_id = tfisp.tfisp_pt_id
			LEFT JOIN
			owl_isp AS isp
			ON tfisp.tfisp_isp_id = isp.isp_id
			LEFT JOIN
			nqm_pt_target_filter_province AS tfpv
			ON pt_id = tfpv.tfpv_pt_id
			LEFT JOIN
			owl_province AS pv
			ON tfpv.tfpv_pv_id = pv.pv_id
			LEFT JOIN
			nqm_pt_target_filter_city AS tfct
			ON pt_id = tfct.tfct_pt_id
			LEFT JOIN
			owl_city AS ct
			ON tfct.tfct_ct_id = ct.ct_id
			LEFT JOIN
			nqm_pt_target_filter_name_tag AS tfnt
			ON pt_id = tfnt.tfnt_pt_id
			LEFT JOIN
			owl_name_tag AS nt
			ON tfnt.tfnt_nt_id = nt.nt_id
			LEFT JOIN
			nqm_pt_target_filter_group_tag AS tfgt
			ON pt_id = tfgt.tfgt_pt_id
			LEFT JOIN
			owl_group_tag AS gt
			ON tfgt.tfgt_gt_id = gt.gt_id
		`).
		Where("pt_id = ?", id).
		Group(`
			pt_id, pt_period, pt_name, pt_enable, pt_comment
		`)

	var loadedPingtask = &nqmModel.PingtaskView{}
	selectPingtask = selectPingtask.Find(loadedPingtask)

	if selectPingtask.Error == gorm.ErrRecordNotFound {
		return nil
	}
	gormExt.ToDefaultGormDbExt(selectPingtask).PanicIfError()

	loadedPingtask.AfterLoad()
	return loadedPingtask
}

func AddAndGetPingtask(pm *nqmModel.PingtaskModify) *nqmModel.PingtaskView {
	txProcessor := &addPingtaskTx{
		pingtask: pm,
	}

	DbFacade.NewSqlxDbCtrl().InTx(txProcessor)
	// :~)

	return GetPingtaskById(txProcessor.pingtaskID)
}

type addPingtaskTx struct {
	pingtask   *nqmModel.PingtaskModify
	pingtaskID int32
}

func (p *addPingtaskTx) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	r := tx.MustExec(
		`
		INSERT INTO nqm_ping_task(pt_period, pt_name, pt_enable, pt_comment)
		VALUES (?, ?, ?, ?)
		`,
		p.pingtask.Period,
		p.pingtask.Name,
		p.pingtask.Enable,
		p.pingtask.Comment,
	)
	pingTaskId := int32(commonDb.ToResultExt(r).LastInsertId())

	pingTaskFilterModifiers["isp"].buildData(tx, pingTaskId, p.pingtask.Filter.IspIds)
	pingTaskFilterModifiers["province"].buildData(tx, pingTaskId, p.pingtask.Filter.ProvinceIds)
	pingTaskFilterModifiers["city"].buildData(tx, pingTaskId, p.pingtask.Filter.CityIds)
	pingTaskFilterModifiers["name_tag"].buildData(tx, pingTaskId, p.pingtask.Filter.NameTagIds)
	pingTaskFilterModifiers["group_tag"].buildData(tx, pingTaskId, p.pingtask.Filter.GroupTagIds)

	p.pingtaskID = pingTaskId
	return commonDb.TxCommit
}

func UpdateAndGetPingtask(id int32, pm *nqmModel.PingtaskModify) *nqmModel.PingtaskView {
	txProcessor := &updatePingtaskTx{
		pingtaskID: id,
		pingtask:   pm,
	}

	DbFacade.NewSqlxDbCtrl().InTx(txProcessor)
	// :~)

	return GetPingtaskById(txProcessor.pingtaskID)
}

var pingTaskFilterModifiers = map[string]*filterModifier{
	"isp":       newFilterModifier("isp", "tfisp", "isp_id"),
	"province":  newFilterModifier("province", "tfpv", "pv_id"),
	"city":      newFilterModifier("city", "tfct", "ct_id"),
	"name_tag":  newFilterModifier("name_tag", "tfnt", "nt_id"),
	"group_tag": newFilterModifier("group_tag", "tfgt", "gt_id"),
}

type updatePingtaskTx struct {
	pingtask   *nqmModel.PingtaskModify
	pingtaskID int32
}

func (u *updatePingtaskTx) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	tx.MustExec(
		`
		UPDATE nqm_ping_task SET
			pt_period = ?,
			pt_name = ?,
			pt_enable = ?,
			pt_comment = ?
		WHERE pt_id = ?
		`,
		u.pingtask.Period,
		u.pingtask.Name,
		u.pingtask.Enable,
		u.pingtask.Comment,
		u.pingtaskID,
	)

	pingTaskId := u.pingtaskID

	pingTaskFilterModifiers["isp"].setData(tx, pingTaskId, u.pingtask.Filter.IspIds)
	pingTaskFilterModifiers["province"].setData(tx, pingTaskId, u.pingtask.Filter.ProvinceIds)
	pingTaskFilterModifiers["city"].setData(tx, pingTaskId, u.pingtask.Filter.CityIds)
	pingTaskFilterModifiers["name_tag"].setData(tx, pingTaskId, u.pingtask.Filter.NameTagIds)
	pingTaskFilterModifiers["group_tag"].setData(tx, pingTaskId, u.pingtask.Filter.GroupTagIds)

	return commonDb.TxCommit
}

func newFilterModifier(
	tableNameSuffix string,
	columnPrefix string,
	propertySuffix string,
) *filterModifier {
	return &filterModifier{
		deleteSql: fmt.Sprintf(
			`
			DELETE FROM nqm_pt_target_filter_%s WHERE %s_pt_id = ?
			`,
			tableNameSuffix, columnPrefix,
		),
		insertSqlTmpl: fmt.Sprintf(
			`
			INSERT INTO nqm_pt_target_filter_%s(%s_pt_id, %s_%s)
			VALUES %%s
			`,
			tableNameSuffix,
			columnPrefix,
			columnPrefix, propertySuffix,
		),
	}
}

type filterModifier struct {
	deleteSql     string
	insertSqlTmpl string
}

func (m *filterModifier) setData(tx *sqlx.Tx, pingTaskId int32, data interface{}) {
	buildDataSql, params := m.buildSqlOfNewData(pingTaskId, data)

	tx.MustExec(m.deleteSql, pingTaskId)
	if buildDataSql != "" {
		tx.MustExec(buildDataSql, params...)
	}
}
func (m *filterModifier) buildData(tx *sqlx.Tx, pingTaskId int32, data interface{}) {
	buildDataSql, params := m.buildSqlOfNewData(pingTaskId, data)

	if buildDataSql != "" {
		tx.MustExec(buildDataSql, params...)
	}
}

func (m *filterModifier) buildSqlOfNewData(pingTaskId int32, data interface{}) (string, []interface{}) {
	values := reflect.ValueOf(data)
	if values.Len() == 0 {
		return "", nil
	}

	newValuesSql := make([]string, values.Len())
	sqlParams := make([]interface{}, values.Len()*2)

	p := 0
	for i := 0; i < values.Len(); i++ {
		newValuesSql[i] = "(?, ?)"

		sqlParams[p] = pingTaskId
		p++
		sqlParams[p] = values.Index(i).Interface()
		p++
	}

	finalSql := fmt.Sprintf(
		m.insertSqlTmpl,
		strings.Join(newValuesSql, ", "),
	)

	return finalSql, sqlParams
}
