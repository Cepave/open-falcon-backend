package falconPortal

import (
	"fmt"
	"time"

	"github.com/Cepave/open-falcon-backend/modules/fe/model/boss"
	"github.com/astaxie/beego/orm"
	log "github.com/sirupsen/logrus"
)

// //make compatible for overall , if need it, please uncomment below
// type AlarmResp struct {
// 	ID          string `json:"id"`
// 	Endpoint    string `json:"endpoint"`
// 	Func        string `json:"func"`
// 	Cond        string `json:"cond"`
// 	Note        string `json:"note"`
// 	MaxStep     int    `json:"max_step"`
// 	CurrentStep int    `json:"current_step"`
// 	// Status     string `json:"status"` //conflict
// 	ProcessStatus string    `json:"process_status"`
// 	StartAt       time.Time `json:"start_at"`
// 	UpdateAt      time.Time `json:"update_at"`
// }

type PlatformExtend struct {
	Contact  []boss.Contactor `json:"contact"`
	IDC      string           `json:"idc"`
	IP       string           `json:"ip"`
	Platform string           `json:"platform"`
}

type AlertsResp struct {
	// //make compatible for overall , if need it, please uncomment below
	// AlarmResp
	PlatformExtend
	Hash       string `json:"hash"`
	CTmpName   string `json:"c_tmp_name"`
	HostName   string `json:"hostname"`
	Metric     string `json:"metric"`
	Author     string `json:"author"`
	TemplateID int    `json:"templateID"`
	Priority   string `json:"priority"`
	Severity   string `json:"severity"`
	Status     string `json:"status"`
	StatusRaw  string `json:"statusRaw"`
	//metricTyoe
	Type       string              `json:"type"`
	Content    string              `json:"content"`
	TimeStart  int64               `json:"timeStart"`
	TimeUpdate int64               `json:"timeUpdate"`
	Notes      []map[string]string `json:"notes"`
	Events     []*Events           `json:"events"`
	Process    string              `json:"process"`
	Function   string              `json:"function"`
	Condition  string              `json:"condition"`
	StepLimit  int                 `json:"stepLimit"`
	Step       int                 `json:"step"`
	// add by 201707
	Activate     int    `json:"active"`
	AlarmType    string `json:"alarm_type"`
	ExtendedBlob string `json:"extended_blob"`
	InternalData int    `json:"internal_data"`
}

func getSeverity(priority string) string {
	severity := "Lower"
	if priority == "0" {
		severity = "High"
	} else if priority == "1" {
		severity = "Medium"
	} else if priority == "2" || priority == "3" {
		severity = "Low"
	}
	return severity
}

func getStatus(statusRaw string) string {
	status := ""
	if statusRaw == "PROBLEM" {
		status = "Triggered"
	} else if statusRaw == "OK" {
		status = "Recovered"
	}
	return status
}

func getDuration(timeTriggered string) string {
	date, _ := time.Parse("2006-01-02 15:04", timeTriggered)
	now := time.Now().Unix()
	diff := now - date.Unix()
	if diff <= 60 {
		return "just now"
	}
	if diff <= 120 {
		return "1 minute ago"
	}
	if diff <= 3600 {
		return fmt.Sprintf("%d minutes ago", diff/60)
	}
	if diff <= 7200 {
		return "1 hour ago"
	}
	if diff <= 3600*24 {
		return fmt.Sprintf("%d hours ago", diff/3600)
	}
	if diff <= 3600*24*2 {
		return "1 day ago"
	}
	return fmt.Sprintf("%d days ago", diff/3600/24)
}

func getNote(hash string, timestamp int64) []map[string]string {
	o := orm.NewOrm()
	var rows []orm.Params
	queryStr := fmt.Sprintf(`SELECT note.id as id, note.event_caseId as event_caseId, note.note as note, note.case_id as case_id, note.status as status, note.timestamp as timestamp, user.name as name from
	(SELECT * from falcon_portal.event_note WHERE event_caseId = '%s' AND timestamp >= FROM_UNIXTIME(%d))
	 note LEFT JOIN uic.user as user on note.user_id = user.id;`, hash, timestamp)

	num, err := o.Raw(queryStr).Values(&rows)
	notes := []map[string]string{}
	if err != nil {
		log.Error(err.Error())
	} else if num == 0 {
		return notes
	} else {
		for _, row := range rows {
			hash := row["event_caseId"].(string)
			time := row["timestamp"].(string)
			time = time[:len(time)-3]
			user := row["name"].(string)
			note := map[string]string{
				"note":   row["note"].(string),
				"status": row["status"].(string),
				"user":   user,
				"hash":   hash,
				"time":   time,
			}
			notes = append(notes, note)
		}
	}
	return notes
}
