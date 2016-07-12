package uic

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego/orm"
	"github.com/toolkits/cache"
	"github.com/toolkits/logger"
	"time"
)

func SelectSessionBySig(sig string) *Session {

	obj := Session{Sig: sig}
	if sig == "" {
		obj.Uid = -1
		return &obj
	}

	err := orm.NewOrm().Read(&obj, "Sig")
	if err != nil {
		if err != orm.ErrNoRows {
			logger.Errorln(err)
		}
		obj.Uid = -1
	}
	return &obj
}

func ReadSessionBySig(sig string) *Session {
	if sig == "" {
		return nil
	}

	key := fmt.Sprintf("session:obj:%s", sig)
	var obj Session
	if err := cache.Get(key, &obj); err != nil {
		objPtr := SelectSessionBySig(sig)
		if objPtr != nil {
			go cache.Set(key, objPtr, time.Hour)
		}
		return objPtr
	}

	return &obj
}

func (this *Session) Save() (int64, error) {
	return orm.NewOrm().Insert(this)
}

func SaveSessionAttrs(uid int64, sig string, expired int) (int64, error) {
	s := &Session{Uid: uid, Sig: sig, Expired: expired}
	return s.Save()
}

func ReadSessionByName(name string) (sig string, expired int) {
	uid := ReadUserIdByName(name)
	sig, expired = ReadSessionByUid(uid)
	return
}

func ReadSessionByUid(uid int64) (sig string, expired int) {
	var ss []Session
	Sessions().Filter("Uid", uid).All(&ss, "Id", "Sig", "expired")
	if len(ss) != 0 {
		sig = ss[0].Sig
		expired = ss[0].Expired
	}
	return
}

func RemoveSessionByUid(uid int64) (num int64, err error) {
	var ss []Session
	Sessions().Filter("Uid", uid).All(&ss, "Id", "Sig")
	if ss == nil || len(ss) == 0 {
		return
	}

	for _, s := range ss {
		num, err = DeleteSessionById(s.Id)
		if err == nil && num > 0 {
			deletekey := fmt.Sprintf("session:obj:%s", s.Sig)
			log.Printf("%v", deletekey)
			cache.Delete(deletekey)
		} else {
			return
		}
	}
	return
}

func ReadSessionById(id int64) (s *Session, err error) {
	var ss []Session
	Sessions().Filter("Id", id).All(&ss, "Id", "Sig", "expired")
	if ss == nil || len(ss) == 0 {
		err = errors.New("not found this session token")
	} else {
		s = &ss[0]
	}
	return
}

func DeleteSessionById(id int64) (int64, error) {
	_, err := ReadSessionById(id)
	if err != nil {
		return 0, err
	}
	r, err := orm.NewOrm().Raw("DELETE FROM `session` WHERE `id` = ?", id).Exec()
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}

func Sessions() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(Session))
}
