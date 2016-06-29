package uic

import (
	"time"
)

type User struct {
	Id      int64     `json:"id"`
	Name    string    `json:"name"`
	Cnname  string    `json:"cnname"`
	Passwd  string    `json:"-"`
	Email   string    `json:"email"`
	Phone   string    `json:"phone"`
	IM      string    `json:"im" orm:"column(im)"`
	QQ      string    `json:"qq" orm:"column(qq)"`
	Role    int       `json:"role"`
	Created time.Time `json:"-" orm:"-"`
}

type Team struct {
	Id      int64     `json:"id"`
	Name    string    `json:"name"`
	Resume  string    `json:"resume"`
	Creator int64     `json:"creator"`
	Created time.Time `json:"-" orm:"-"`
}

type RelTeamUser struct {
	Id  int64
	Tid int64
	Uid int64
}

type Session struct {
	Id      int64
	Uid     int64
	Sig     string
	Expired int
}
