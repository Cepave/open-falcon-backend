package falconPortal

import (
	"fmt"
	"time"

	"github.com/Cepave/fe/model/uic"
	"github.com/astaxie/beego/orm"
)

func CloseEvent(username string, colsed_note string, id string) (err error) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	userid := uic.ReadUserIdByName(username)
	_, err = q.Raw("Update event_cases SET user_modified = ?, closed_at = ?, status = ?, closed_note = ? WHERE id = ?", userid, time.Now(), "SOLVED", colsed_note, id).Exec()
	return
}

func UpdateProcess(processNoteID int, processStatus string, id string) (err error) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	_, err = q.Raw("Update event_cases SET process_note = ?, process_status = ? WHERE id = ?", processNoteID, processStatus, id).Exec()
	return
}

func AddNote(username string, note string, id string, status string, caseId string) (err error) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	userid := uic.ReadUserIdByName(username)
	sqlbase := fmt.Sprintf("SET event_caseId = %s , user_id = %d, timestamp = %v", id, userid, time.Now())
	if note != "" {
		sqlbase = fmt.Sprintf("%s , note = '%s'", sqlbase, note)
	}
	if status != "" {
		sqlbase = fmt.Sprintf("%s, status = '%s'", sqlbase, status)
	}
	//for set boss case id
	if caseId != "" {
		sqlbase = fmt.Sprintf("%s, case_id = '%s'", sqlbase, caseId)
	}
	var eventNote EventNote
	q.Raw(fmt.Sprintf("Insert INTO event_note %s", sqlbase)).QueryRow(&eventNote)

	if eventNote.Status == "Sovled" || eventNote.Status == "In Process" {
		err = UpdateProcess(eventNote.Id, eventNote.Status, id)
	}
	return
}
