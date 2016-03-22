package event

import (
	"fmt"
	"time"

	"github.com/Cepave/fe/model/uic"
	"github.com/astaxie/beego/orm"
)

func GetEvent(startTime int64, endTime int64, priority int, status string) (result []Event, err error) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	flag := false
	query_tmp := ""
	if startTime != 0 && endTime != 0 {
		flag = true
		query_tmp = fmt.Sprintf(" %v timestamp >= %d and  update_at >= %d", query_tmp, startTime, endTime)
	}
	if priority != -1 {
		if flag {
			query_tmp = fmt.Sprintf("%v and priority = %d", query_tmp, priority)
		} else {
			flag = true
			query_tmp = fmt.Sprintf("%v priority = %d", query_tmp, priority)
		}
	}
	if status != "ALL" {
		if flag {
			query_tmp = fmt.Sprintf("%v and status = '%s'", query_tmp, status)
		} else {
			flag = true
			query_tmp = fmt.Sprintf("%v status = '%s'", query_tmp, status)
		}
	}
	if query_tmp != "" {
		_, err = q.Raw(fmt.Sprintf("select * from event where %v", query_tmp)).QueryRows(&result)
	} else {
		_, err = q.Raw("select * from event").QueryRows(&result)
	}

	if len(result) == 0 {
		result = []Event{}
	}

	return
}

func CloseEvent(username string, id string) (err error) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	userid := uic.ReadUserIdByName(username)
	_, err = q.Raw("Update event SET user_modified = ?, closed_at = ?, status = ? WHERE id = ?", userid, time.Now(), "SOLVED", id).Exec()
	return
}

func CountNumOfTlp() (c int, err error) {
	var h []Tpl
	q := orm.NewOrm()
	q.Using("falcon_portal")
	_, err = q.Raw("select * from `tpl`").QueryRows(&h)
	c = len(h)
	return
}
