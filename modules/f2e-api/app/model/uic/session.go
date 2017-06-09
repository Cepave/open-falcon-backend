package uic

import (
	"strings"
	"time"

	con "github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
	uuid "github.com/satori/go.uuid"
)

type Session struct {
	ID      int64
	Uid     int64
	Sig     string
	Expired int64
}

func (this Session) ExpiredUnixGen() int64 {
	// will expired 3 day
	return time.Now().Unix() + 3600*24*3
}

func (this Session) UUIDGen() string {
	sig := uuid.NewV1().String()
	sig = strings.Replace(sig, "-", "", -1)
	return sig
}

func (this Session) DeleteSessionByUid() (err error) {
	db := con.Con()
	if dt := db.Uic.Delete(this, "uid = ?", this.Uid); dt.Error != nil {
		err = dt.Error
	}
	return
}

// will find vaild session
// and delete other expired session
func (this Session) FindVaildSession() Session {
	db := con.Con()
	db.Uic.Model(&this).Order("id", true).Where("uid = ? AND expired > ?", this.Uid, time.Now().Unix()).Limit(1).Scan(&this)
	if this.ID != 0 {
		db.Uic.Delete(this, "id < ? AND uid = ?", this.ID, this.Uid)
	} else {
		// this mean no any vaild session, create one!
		this.Expired = this.ExpiredUnixGen()
		this.Sig = this.UUIDGen()
		db.Uic.Model(&this).Save(&this)
	}
	return this
}

func (this Session) TableName() string {
	return "session"
}
