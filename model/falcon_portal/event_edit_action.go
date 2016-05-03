package falconPortal

import (
	"github.com/Cepave/fe/model/uic"
	"github.com/astaxie/beego/orm"
	"time"
)

func CloseEvent(username string, colsed_note string, id string) (err error) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	userid := uic.ReadUserIdByName(username)
	_, err = q.Raw("Update event_cases SET user_modified = ?, closed_at = ?, status = ?, closed_note = ? WHERE id = ?", userid, time.Now(), "SOLVED", colsed_note, id).Exec()
	return
}
