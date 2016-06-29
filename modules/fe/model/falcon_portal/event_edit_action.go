package falconPortal

import (
	"fmt"
	"time"

	"github.com/Cepave/fe/model/uic"
	"github.com/astaxie/beego/orm"
)

//will deprecated
func CloseEvent(username string, colsed_note string, id string) (err error) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	userid := uic.ReadUserIdByName(username)
	_, err = q.Raw("Update event_cases SET user_modified = ?, closed_at = ?, status = ?, closed_note = ? WHERE id = ?", userid, time.Now(), "SOLVED", colsed_note, id).Exec()
	return
}

func UpdateCaseStatus(eventcaseid string, processNoteID int, processStatus string) (err error) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	_, err = q.Raw("Update event_cases SET process_note = ?, process_status = ? WHERE id = ?", processNoteID, processStatus, eventcaseid).Exec()
	return
}

func AddNote(username string, processNote string, eventcaseid string, processStatus string, bossId string) (err error) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	userid := uic.ReadUserIdByName(username)
	sqlbase := fmt.Sprintf("SET event_caseId = '%s' , user_id = %d", eventcaseid, userid)
	if processNote != "" {
		sqlbase = fmt.Sprintf("%s , note = '%s'", sqlbase, processNote)
	}
	//for set boss case id
	if bossId != "" {
		sqlbase = fmt.Sprintf("%s, case_id = '%s'", sqlbase, bossId)
	}
	if processStatus != "" {
		sqlbase = fmt.Sprintf("%s, status = '%s'", sqlbase, processStatus)
	}
	var processNoteID int
	q.Raw(fmt.Sprintf("Insert INTO event_note %s, timestamp = ? ;", sqlbase), time.Now()).Exec()
	err = q.Raw("SELECT LAST_INSERT_ID()").QueryRow(&processNoteID)
	if processNoteID != 0 && (processStatus == "resolved" || processStatus == "in progress" || processStatus == "ignored") && err == nil {
		err = UpdateCaseStatus(eventcaseid, processNoteID, processStatus)
	}
	return
}
