package nqm

import (
	"database/sql"
	"fmt"
	"net"
	"reflect"

	"github.com/Cepave/open-falcon-backend/common/utils"
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"

	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	sqlxExt "github.com/Cepave/open-falcon-backend/common/db/sqlx"
	gormExt "github.com/Cepave/open-falcon-backend/common/gorm"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	tb "github.com/Cepave/open-falcon-backend/common/textbuilder"
	sqlb "github.com/Cepave/open-falcon-backend/common/textbuilder/sql"
)

type ErrDuplicatedNqmAgent struct {
	ConnectionId string
}

func (err ErrDuplicatedNqmAgent) Error() string {
	return fmt.Sprintf("Duplicated NQM agent. Connection Id: [%s]", err.ConnectionId)
}

// Add and retrieve detail data of agent
//
// Errors:
// 		ErrDuplicatedNqmAgent - The agent is existing with the same connection id
//		ErrNotInSameHierarchy - The city is not belonging to the province
func AddAgent(addedAgent *nqmModel.AgentForAdding) (*nqmModel.Agent, error) {
	/**
	 * Checks the hierarchy over administrative region
	 */
	err := owlDb.CheckHierarchyForCity(addedAgent.ProvinceId, addedAgent.CityId)
	if err != nil {
		return nil, err
	}
	// :~)

	/**
	 * Executes the insertion of agent and its related data
	 */
	txProcessor := &addAgentTx{
		agent: addedAgent,
	}

	DbFacade.NewSqlxDbCtrl().InTx(txProcessor)
	// :~)

	if txProcessor.err != nil {
		return nil, txProcessor.err
	}

	return GetAgentById(addedAgent.Id), nil
}

func UpdateAgent(oldAgent *nqmModel.Agent, updatedAgent *nqmModel.AgentForAdding) (*nqmModel.Agent, error) {
	/**
	 * Checks the hierarchy over administrative region
	 */
	err := owlDb.CheckHierarchyForCity(updatedAgent.ProvinceId, updatedAgent.CityId)
	if err != nil {
		return nil, err
	}
	// :~)

	txProcessor := &updateAgentTx{
		updatedAgent: updatedAgent,
		oldAgent:     oldAgent.ToAgentForAdding(),
	}

	DbFacade.NewSqlxDbCtrl().InTx(txProcessor)

	return GetAgentById(oldAgent.Id), nil
}

func GetAgentById(agentId int32) *nqmModel.Agent {
	var selectAgent = DbFacade.GormDb.Model(&nqmModel.Agent{}).
		Select(`
			ag_id, ag_name, ag_connection_id, ag_hostname, ag_ip_address, ag_status, ag_comment, ag_last_heartbeat,
			COUNT(DISTINCT pt.pt_id) AS ag_num_of_enabled_pingtasks,
			isp_id, isp_name, pv_id, pv_name, ct_id, ct_name, nt_id, nt_value,
			GROUP_CONCAT(gt.gt_id ORDER BY gt_name ASC SEPARATOR ',') AS gt_ids,
			GROUP_CONCAT(gt.gt_name ORDER BY gt_name ASC SEPARATOR '\0') AS gt_names
		`).
		Joins(`
			LEFT JOIN
			nqm_agent_ping_task AS apt
			ON ag_id = apt.apt_ag_id
			LEFT JOIN
			nqm_ping_task AS pt
			ON apt.apt_pt_id = pt.pt_id AND pt.pt_enable=true
			INNER JOIN
			owl_isp AS isp
			ON ag_isp_id = isp.isp_id
			INNER JOIN
			owl_province AS pv
			ON ag_pv_id = pv.pv_id
			INNER JOIN
			owl_city AS ct
			ON ag_ct_id = ct.ct_id
			INNER JOIN
			owl_name_tag AS nt
			ON ag_nt_id = nt.nt_id
			LEFT OUTER JOIN
			nqm_agent_group_tag AS agt
			ON ag_id = agt.agt_ag_id
			LEFT OUTER JOIN
			owl_group_tag AS gt
			ON agt.agt_gt_id = gt.gt_id
		`).
		Where("ag_id = ?", agentId).
		Group(`
			ag_id, ag_name, ag_connection_id, ag_hostname, ag_ip_address, ag_status, ag_comment, ag_last_heartbeat,
			isp_id, isp_name, pv_id, pv_name, ct_id, ct_name, nt_id, nt_value
		`)

	var loadedAgent = &nqmModel.Agent{}
	selectAgent = selectAgent.Find(loadedAgent)

	if selectAgent.Error == gorm.ErrRecordNotFound {
		return nil
	}
	gormExt.ToDefaultGormDbExt(selectAgent).PanicIfError()

	loadedAgent.AfterLoad()
	return loadedAgent
}

func GetSimpleAgent1ById(agentId int32) *nqmModel.SimpleAgent1 {
	var result nqmModel.SimpleAgent1

	if !DbFacade.SqlxDbCtrl.GetOrNoRow(
		&result,
		`
		SELECT ag_id, ag_name, ag_hostname, ag_ip_address,
			ag_isp_id, ag_pv_id, ag_ct_id, ag_nt_id
		FROM nqm_agent
		WHERE ag_id = ?
		`,
		agentId,
	) {
		return nil
	}

	return &result
}

func LoadSimpleAgent1sByFilter(filter *nqmModel.AgentFilter) []*nqmModel.SimpleAgent1 {
	var result []*nqmModel.SimpleAgent1

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
	sqlArgs = utils.AppendToAny(sqlArgs, filter.Hostname)
	sqlArgs = utils.AppendToAny(sqlArgs, filter.ConnectionId)

	sqlArgs = utils.MakeAbstractArray(sqlArgs).
		MapTo(
			utils.TypedFuncToMapper(func(v string) string {
				return v + "%"
			}),
			utils.TypeOfString,
		).GetArrayOfAny()

	// Process ip address
	sqlArgs = append(
		sqlArgs,
		utils.MakeAbstractArray(filter.IpAddress).
			MapTo(
				utils.TypedFuncToMapper(func(v string) []byte {
					bytes, err := commonDb.IpV4ToBytesForLike(v)
					if err != nil {
						panic(err)
					}

					return bytes
				}),
				utils.ATypeOfByte,
			).GetArrayOfAny()...,
	)
	// :~)

	DbFacade.SqlxDbCtrl.Select(
		&result,
		`
		SELECT ag_id, ag_name, ag_hostname, ag_ip_address,
			ag_isp_id, ag_pv_id, ag_ct_id, ag_nt_id
		FROM nqm_agent
		`+
			sqlb.Where(
				sqlb.And(
					buildRepeatOr("ag_name LIKE ?", filter.Name),
					buildRepeatOr("ag_hostname LIKE ?", filter.Hostname),
					buildRepeatOr("ag_connection_id LIKE ?", filter.ConnectionId),
					buildRepeatOr("ag_ip_address LIKE ?", filter.IpAddress),
				),
			).String(),
		sqlArgs...,
	)

	return result
}

type agentInCity struct {
	AgentId        int32  `db:"ag_id"`
	AgentName      string `db:"ag_name"`
	AgentHostname  string `db:"ag_hostname"`
	AgentIpAddress net.IP `db:"ag_ip_address"`
	CityId         int16  `db:"ct_id"`
	CityName       string `db:"ct_name"`
	CityPostCode   string `db:"ct_post_code"`
}

func (a *agentInCity) getCity2() *owlModel.City2 {
	return &owlModel.City2{
		Id:       a.CityId,
		Name:     a.CityName,
		PostCode: a.CityPostCode,
	}
}
func (a *agentInCity) getSimpleAgent1() *nqmModel.SimpleAgent1 {
	return &nqmModel.SimpleAgent1{
		Id:        a.AgentId,
		Name:      a.AgentName,
		Hostname:  a.AgentName,
		IpAddress: a.AgentIpAddress,
	}
}

var typeOfSimpleAgent1 = reflect.TypeOf(&nqmModel.SimpleAgent1{})

func LoadEffectiveAgentsInProvince(provinceId int16) []*nqmModel.SimpleAgent1InCity {
	var rawResult []*agentInCity

	/**
	 * Query data from database
	 */
	DbFacade.SqlxDbCtrl.Select(
		&rawResult,
		`
		SELECT DISTINCT ag.ag_id, ag.ag_name, ag.ag_hostname, ag.ag_ip_address,
			ct.ct_id, ct.ct_name, ct.ct_post_code
		FROM nqm_agent AS ag
			/**
			 * Only lists the agents(enabled) has ping task(enabled)
			 */
			INNER JOIN
			nqm_agent_ping_task AS apt
			ON ag.ag_id = apt.apt_ag_id
				AND ag.ag_status = TRUE
			INNER JOIN
			nqm_ping_task AS pt
			ON pt.pt_id = apt.apt_pt_id
				AND pt.pt_enable = TRUE
			# //:~)
			INNER JOIN
			owl_city AS ct
			ON ag.ag_ct_id = ct.ct_id
				AND ct.ct_pv_id = ?
		`,
		provinceId,
	)
	// :~)

	/**
	 * Grouping data
	 */
	grouping := utils.NewGroupingProcessor(typeOfSimpleAgent1)
	for _, agent := range rawResult {
		grouping.Put(agent.getCity2(), agent.getSimpleAgent1())
	}
	// :~)

	keys := grouping.Keys()
	result := make([]*nqmModel.SimpleAgent1InCity, len(keys))

	/**
	 * Builds the result list of agents grouped by city
	 */
	i := 0
	for _, key := range keys {
		result[i] = &nqmModel.SimpleAgent1InCity{
			City:   grouping.KeyObject(key).(*owlModel.City2),
			Agents: grouping.Children(key).([]*nqmModel.SimpleAgent1),
		}

		i++
	}
	// :~)

	return result
}

// Lists the agents by query condition
func ListAgents(query *nqmModel.AgentQuery, paging commonModel.Paging) ([]*nqmModel.Agent, *commonModel.Paging) {
	var result []*nqmModel.Agent

	var funcTxLoader gormExt.TxCallbackFunc = func(txGormDb *gorm.DB) commonDb.TxFinale {
		/**
		 * Retrieves the page of data
		 */
		var selectAgent = txGormDb.Model(&nqmModel.Agent{}).
			Select(`SQL_CALC_FOUND_ROWS
				ag_id, ag_name, ag_connection_id, ag_hostname, ag_ip_address, ag_status, ag_comment, ag_last_heartbeat,
				isp_id, isp_name, pv_id, pv_name, ct_id, ct_name, nt_id, nt_value,
				COUNT(gt.gt_id) AS gt_number,
				GROUP_CONCAT(gt.gt_id ORDER BY gt_name ASC SEPARATOR ',') AS gt_ids,
				GROUP_CONCAT(gt.gt_name ORDER BY gt_name ASC SEPARATOR '\0') AS gt_names
			`).
			Joins(`
				INNER JOIN
				owl_isp AS isp
				ON ag_isp_id = isp.isp_id
				INNER JOIN
				owl_province AS pv
				ON ag_pv_id = pv.pv_id
				INNER JOIN
				owl_city AS ct
				ON ag_ct_id = ct.ct_id
				INNER JOIN
				owl_name_tag AS nt
				ON ag_nt_id = nt.nt_id
				LEFT OUTER JOIN
				nqm_agent_group_tag AS agt
				ON ag_id = agt.agt_ag_id
				LEFT OUTER JOIN
				owl_group_tag AS gt
				ON agt.agt_gt_id = gt.gt_id
			`).
			Limit(paging.Size).
			Group(`
				ag_id, ag_name, ag_connection_id, ag_hostname, ag_ip_address, ag_status, ag_comment, ag_last_heartbeat,
				isp_id, isp_name, pv_id, pv_name, ct_id, ct_name, nt_id, nt_value
			`).
			Order(buildSortingClauseOfAgents(&paging)).
			Offset(paging.GetOffset())

		if query.Name != "" {
			selectAgent = selectAgent.Where("ag_name LIKE ?", query.Name+"%")
		}
		if query.ConnectionId != "" {
			selectAgent = selectAgent.Where("ag_connection_id LIKE ?", query.ConnectionId+"%")
		}
		if query.Hostname != "" {
			selectAgent = selectAgent.Where("ag_hostname LIKE ?", query.Hostname+"%")
		}
		if query.HasIspId {
			selectAgent = selectAgent.Where("ag_isp_id = ?", query.IspId)
		}
		if query.HasStatusCondition {
			selectAgent = selectAgent.Where("ag_status = ?", query.Status)
		}
		if query.IpAddress != "" {
			selectAgent = selectAgent.Where("ag_ip_address LIKE ?", query.GetIpForLikeCondition())
		}
		// :~)

		gormExt.ToDefaultGormDbExt(selectAgent.Find(&result)).PanicIfError()

		return commonDb.TxCommit
	}

	gormExt.ToDefaultGormDbExt(DbFacade.GormDb).SelectWithFoundRows(
		funcTxLoader, &paging,
	)

	/**
	 * Loads group tags
	 */
	for _, agent := range result {
		agent.AfterLoad()
	}
	// :~)

	return result, &paging
}

var orderByDialectForAgents = commonModel.NewSqlOrderByDialect(
	map[string]string{
		"status":              "ag_status",
		"name":                "ag_name",
		"connection_id":       "ag_connection_id",
		"comment":             "ag_comment",
		"province":            "pv_name",
		"city":                "ct_name",
		"last_heartbeat_time": "ag_last_heartbeat",
		"name_tag":            "nt_value",
	},
)

func init() {
	originFunc := orderByDialectForAgents.FuncEntityToSyntax
	orderByDialectForAgents.FuncEntityToSyntax = func(entity *commonModel.OrderByEntity) (string, error) {
		switch entity.Expr {
		case "group_tag":
			return owlDb.GetSyntaxOfOrderByGroupTags(entity), nil
		}

		return originFunc(entity)
	}
}

func buildSortingClauseOfAgents(paging *commonModel.Paging) string {
	if len(paging.OrderBy) == 0 {
		paging.OrderBy = append(paging.OrderBy, &commonModel.OrderByEntity{"status", commonModel.Descending})
	}

	if len(paging.OrderBy) == 1 {
		switch paging.OrderBy[0].Expr {
		case "province":
			paging.OrderBy = append(paging.OrderBy, &commonModel.OrderByEntity{"city", commonModel.Ascending})
		}
	}

	if paging.OrderBy[len(paging.OrderBy)-1].Expr != "last_heartbeat_time" {
		paging.OrderBy = append(paging.OrderBy, &commonModel.OrderByEntity{"last_heartbeat_time", commonModel.Descending})
	}

	querySyntax, err := orderByDialectForAgents.ToQuerySyntax(paging.OrderBy)
	gormExt.DefaultGormErrorConverter.PanicIfError(err)

	return querySyntax
}

type addAgentTx struct {
	agent *nqmModel.AgentForAdding
	err   error
}

func (agentTx *addAgentTx) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	agentTx.prepareHost(tx)

	agentTx.agent.NameTagId = owlDb.BuildAndGetNameTagId(
		tx, agentTx.agent.NameTagValue,
	)

	agentTx.addAgent(tx)
	if agentTx.err != nil {
		return commonDb.TxRollback
	}

	agentTx.prepareGroupTags(tx)
	return commonDb.TxCommit
}
func (agentTx *addAgentTx) prepareHost(tx *sqlx.Tx) {
	newAgent := agentTx.agent

	tx.MustExec(
		`
		INSERT INTO host(hostname, ip, agent_version, plugin_version)
		SELECT ?, ?, '0', '0'
		FROM DUAL
		WHERE NOT EXISTS (
			SELECT *
			FROM host
			WHERE hostname = ?
		)
		`,
		newAgent.Hostname,
		newAgent.GetIpAddressAsString(),
		newAgent.Hostname,
	)
}
func (agentTx *addAgentTx) addAgent(tx *sqlx.Tx) {
	txExt := sqlxExt.ToTxExt(tx)
	newAgent := agentTx.agent

	addedNqmAgent := txExt.NamedExec(
		`
		INSERT INTO nqm_agent(
			ag_name, ag_connection_id, ag_status,
			ag_hostname, ag_ip_address,
			ag_isp_id, ag_pv_id, ag_ct_id, ag_nt_id,
			ag_comment,
			ag_hs_id
		)
		SELECT :name, :connection_id, :status,
			:hostname, :ip_address,
			:isp_id, :province_id, :city_id, :name_tag_id,
			:comment,
			(
				SELECT id
				FROM host
				WHERE hostname = :hostname
			)
		FROM DUAL
		WHERE NOT EXISTS (
			SELECT *
			FROM nqm_agent
			WHERE ag_connection_id = :connection_id
		)
		`,
		map[string]interface{}{
			"status":        newAgent.Status,
			"hostname":      newAgent.Hostname,
			"ip_address":    newAgent.GetIpAddressAsBytes(),
			"isp_id":        newAgent.IspId,
			"province_id":   newAgent.ProvinceId,
			"city_id":       newAgent.CityId,
			"name_tag_id":   newAgent.NameTagId,
			"connection_id": newAgent.ConnectionId,
			"name": sql.NullString{
				newAgent.Name,
				newAgent.Name != "",
			},
			"comment": sql.NullString{
				newAgent.Comment,
				newAgent.Comment != "",
			},
		},
	)

	/**
	 * Rollback if the NQM agent is existing(duplicated by connection id)
	 */
	if commonDb.ToResultExt(addedNqmAgent).RowsAffected() == 0 {
		agentTx.err = ErrDuplicatedNqmAgent{newAgent.ConnectionId}
		return
	}
	// :~)

	txExt.Get(
		&newAgent.Id,
		`
		SELECT ag_id FROM nqm_agent
		WHERE ag_connection_id = ?
		`,
		newAgent.ConnectionId,
	)
}
func (agentTx *addAgentTx) prepareGroupTags(tx *sqlx.Tx) {
	newAgent := agentTx.agent
	buildGroupTagsForAgent(
		tx, newAgent.Id, newAgent.GroupTags,
	)
}

type updateAgentTx struct {
	updatedAgent *nqmModel.AgentForAdding
	oldAgent     *nqmModel.AgentForAdding
}

func (agentTx *updateAgentTx) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	agentTx.loadNameTagId(tx)

	updatedAgent, oldAgent := agentTx.updatedAgent, agentTx.oldAgent
	tx.MustExec(
		`
		UPDATE nqm_agent
		SET ag_name = ?,
			ag_comment = ?,
			ag_status = ?,
			ag_isp_id = ?,
			ag_pv_id = ?,
			ag_ct_id = ?,
			ag_nt_id = ?
		WHERE ag_id = ?
		`,
		sql.NullString{updatedAgent.Name, updatedAgent.Name != ""},
		sql.NullString{updatedAgent.Comment, updatedAgent.Comment != ""},
		updatedAgent.Status,
		updatedAgent.IspId,
		updatedAgent.ProvinceId,
		updatedAgent.CityId,
		updatedAgent.NameTagId,
		oldAgent.Id,
	)

	agentTx.updateGroupTags(tx)
	return commonDb.TxCommit
}
func (agentTx *updateAgentTx) loadNameTagId(tx *sqlx.Tx) {
	updatedAgent, oldAgent := agentTx.updatedAgent, agentTx.oldAgent

	if updatedAgent.NameTagValue == oldAgent.NameTagValue {
		return
	}

	updatedAgent.NameTagId = owlDb.BuildAndGetNameTagId(
		tx, updatedAgent.NameTagValue,
	)
}
func (agentTx *updateAgentTx) updateGroupTags(tx *sqlx.Tx) {
	updatedAgent, oldAgent := agentTx.updatedAgent, agentTx.oldAgent
	if updatedAgent.AreGroupTagsSame(oldAgent) {
		return
	}

	tx.MustExec(
		`
		DELETE FROM nqm_agent_group_tag
		WHERE agt_ag_id = ?
		`,
		oldAgent.Id,
	)

	buildGroupTagsForAgent(
		tx, oldAgent.Id, updatedAgent.GroupTags,
	)
}

func buildGroupTagsForAgent(tx *sqlx.Tx, agentId int32, groupTags []string) {
	owlDb.BuildGroupTags(
		tx, groupTags,
		func(tx *sqlx.Tx, groupTag string) {
			tx.MustExec(
				`
				INSERT INTO nqm_agent_group_tag(agt_ag_id, agt_gt_id)
				VALUES(
					?,
					(
						SELECT gt_id
						FROM owl_group_tag
						WHERE gt_name = ?
					)
				)
				`,
				agentId,
				groupTag,
			)
		},
	)
}
