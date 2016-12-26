package owl

import (
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	"github.com/jmoiron/sqlx"
	sqlxExt "github.com/Cepave/open-falcon-backend/common/db/sqlx"
)

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
