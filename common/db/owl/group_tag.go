package owl

import (
	"github.com/jmoiron/sqlx"
)

type ProcessGroupTagFunc func(tx *sqlx.Tx, nameOfGroupTag string)

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
