package model

import (
	"database/sql"
)

type BossHost struct {
	Hostname string         `db:"hostname"`
	Ip       string         `db:"ip"`
	Activate sql.NullInt64  `db:"activate"`
	Platform sql.NullString `db:"platform"`
}
