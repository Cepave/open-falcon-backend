package uic

import (
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
)

type RelTeamUser struct {
	ID  int64
	Tid int64
	Uid int64
}

func (this RelTeamUser) TableName() string {
	return "rel_team_user"
}

func (this RelTeamUser) Me() {
	db := config.Con()
	db.Uic.Where("id = 1")
}
