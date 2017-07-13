package owl

import (
	"fmt"
	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/common/utils"

	"github.com/Cepave/open-falcon-backend/common/db"
	sqlxExt "github.com/Cepave/open-falcon-backend/common/db/sqlx"
	"github.com/jmoiron/sqlx"

	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	t "github.com/Cepave/open-falcon-backend/common/textbuilder"
	tsql "github.com/Cepave/open-falcon-backend/common/textbuilder/sql"
)

type ProcessGroupTagFunc func(tx *sqlx.Tx, nameOfGroupTag string)

var orderByDialectForGroupTag = model.NewSqlOrderByDialect(
	map[string]string{
		"name": "gt_name",
	},
)

func ListGroupTags(name string, p *model.Paging) []*owlModel.GroupTag {
	var result = make([]*owlModel.GroupTag, 0)

	if len(p.OrderBy) == 0 {
		p.OrderBy = append(p.OrderBy, &model.OrderByEntity{"name", utils.Ascending})
	}

	var sqlParams = make([]interface{}, 0)
	if name != "" {
		sqlParams = append(sqlParams, name+"%")
	}

	txFunc := sqlxExt.TxCallbackFunc(func(tx *sqlx.Tx) db.TxFinale {
		sql := fmt.Sprintf(
			`
			SELECT SQL_CALC_FOUND_ROWS gt_id, gt_name
			FROM owl_group_tag
			%s
			%s
			`,
			tsql.Where(
				t.Dsl.S("gt_name LIKE ?").Post().Viable(name != ""),
			),
			model.GetOrderByAndLimit(p, orderByDialectForGroupTag),
		)

		sqlxExt.ToTxExt(tx).Select(
			&result, sql, sqlParams...,
		)

		return db.TxCommit
	})

	DbFacade.SqlxDbCtrl.SelectWithFoundRows(
		txFunc, p,
	)

	return result
}

func BuildGroupTags(tx *sqlx.Tx, groupTags []string, handleFunc ProcessGroupTagFunc) {
	for _, groupTag := range groupTags {
		tx.MustExec(
			`
			INSERT INTO owl_group_tag(gt_name)
			SELECT ?
			FROM DUAL
			WHERE NOT EXISTS (
				SELECT *
				FROM owl_group_tag
				WHERE gt_name = ?
			)
			`,
			groupTag,
			groupTag,
		)

		handleFunc(tx, groupTag)
	}
}

func GetGroupTagById(id int32) *owlModel.GroupTag {
	groupTag := &owlModel.GroupTag{}

	if !DbFacade.SqlxDbCtrl.GetOrNoRow(
		groupTag,
		`
		SELECT gt_id, gt_name
		FROM owl_group_tag
		WHERE gt_id = ?
		`,
		id,
	) {
		return nil
	}

	return groupTag
}
