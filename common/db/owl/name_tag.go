package owl

import (
	"fmt"
	"github.com/jmoiron/sqlx"

	"github.com/Cepave/open-falcon-backend/common/db"
	sqlxExt "github.com/Cepave/open-falcon-backend/common/db/sqlx"
	"github.com/Cepave/open-falcon-backend/common/model"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	t "github.com/Cepave/open-falcon-backend/common/textbuilder"
	tsql "github.com/Cepave/open-falcon-backend/common/textbuilder/sql"
	"github.com/Cepave/open-falcon-backend/common/utils"
)

var orderByDialectForNameTag = model.NewSqlOrderByDialect(
	map[string]string{
		"value": "nt_value",
	},
)

// This function skipes the -1 name tag
func ListNameTags(value string, p *model.Paging) []*owlModel.NameTag {
	var result = make([]*owlModel.NameTag, 0)

	if len(p.OrderBy) == 0 {
		p.OrderBy = append(p.OrderBy, &model.OrderByEntity{"value", utils.Ascending})
	}

	var sqlParams = make([]interface{}, 0)
	if value != "" {
		sqlParams = append(sqlParams, value+"%")
	}

	txFunc := sqlxExt.TxCallbackFunc(func(tx *sqlx.Tx) db.TxFinale {
		sql := fmt.Sprintf(
			`
			SELECT SQL_CALC_FOUND_ROWS nt_id, nt_value
			FROM owl_name_tag
			%s
			%s
			`,
			tsql.Where(
				tsql.And(
					t.Dsl.S("nt_id >= 1"),
					t.Dsl.S("nt_value LIKE ?").
						Post().Viable(value != ""),
				),
			),
			model.GetOrderByAndLimit(p, orderByDialectForNameTag),
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

func BuildAndGetNameTagId(tx *sqlx.Tx, valueOfNameTag string) int16 {
	if valueOfNameTag == "" {
		return -1
	}

	tx.MustExec(
		`
		INSERT INTO owl_name_tag(nt_value)
		SELECT ?
		FROM DUAL
		WHERE NOT EXISTS (
			SELECT *
			FROM owl_name_tag
			WHERE nt_value = ?
		)
		`,
		valueOfNameTag, valueOfNameTag,
	)

	var nameTagId int16
	sqlxExt.ToTxExt(tx).Get(
		&nameTagId,
		`
		SELECT nt_id FROM owl_name_tag
		WHERE nt_value = ?
		`,
		valueOfNameTag,
	)

	return nameTagId
}

func GetNameTagById(id int16) *owlModel.NameTag {
	nameTag := &owlModel.NameTag{}

	if !DbFacade.SqlxDbCtrl.GetOrNoRow(
		nameTag,
		`
		SELECT nt_id, nt_value
		FROM owl_name_tag
		WHERE nt_id = ?
		`,
		id,
	) {
		return nil
	}

	return nameTag
}
