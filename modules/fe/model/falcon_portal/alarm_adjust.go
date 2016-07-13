package falconPortal

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego/orm"
)

var (
	removedStatus  = "REMOVED"
	modifiedStatus = "UNKNOWN"
	processStatus  = "ignored"
)

func WhenStrategyUpdated(strategyId int) (err error, affectedRows int) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	//set all related event cases status with "UNKNOWN"
	_, err = q.Raw("UPDATE event_cases SET status = ? WHERE strategy_id = ?", modifiedStatus, strategyId).Exec()
	if err != nil {
		return
	}
	affectedAlerms := []string{}
	_, err = q.Raw("SELECT id FROM event_cases WHERE strategy_id = ? ", strategyId).QueryRows(&affectedAlerms)
	affectedRows = len(affectedAlerms)
	// AddNote("root", fmt.Sprintf("it's cause of strategyId: %d has been modified by user", strategyId), strings.Join(affectedAlerms, ","), processStatus, "")
	return
}

func WhenStrategyDeleted(strategyId int) (err error, affectedRows int) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	//set all related event cases status with "REMOVED"
	_, err = q.Raw("UPDATE event_cases SET status = ? WHERE strategy_id = ?", removedStatus, strategyId).Exec()
	if err != nil {
		return
	}
	affectedAlerms := []string{}
	_, err = q.Raw("SELECT id FROM event_cases WHERE strategy_id = ? ", strategyId).QueryRows(&affectedAlerms)
	affectedRows = len(affectedAlerms)
	// AddNote("root", fmt.Sprintf("it's cause of strategyId: %d has been deleted by user", strategyId), strings.Join(affectedAlerms, ","), processStatus, "")
	return
}

func WhenTempleteDeleted(templateId int) (err error, affectedRows int) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	//set all related event cases status with "REMOVED"
	_, err = q.Raw("UPDATE event_cases SET status = ? WHERE template_id = ?", removedStatus, templateId).Exec()
	if err != nil {
		return
	}
	affectedAlerms := []string{}
	_, err = q.Raw(fmt.Sprintf("SELECT id FROM event_cases WHERE template_id = %d ", templateId)).QueryRows(&affectedAlerms)
	affectedRows = len(affectedAlerms)
	// AddNote("root", fmt.Sprintf("it's cause of templateId: %d has been deleted by user", templateId), strings.Join(affectedAlerms, ","), processStatus, "")
	return
}

func WhenTempleteUnbind(templateId int, hostgroupId int) (err error, affectedRows int) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	//set all related event cases status with "REMOVED"
	hostnameList := []string{}
	//get hostname list by hostgroupId
	_, err = q.Raw("SELECT hostname FROM host INNER JOIN grp_host ON host.id = grp_host.host_id WHERE grp_id = ?", hostgroupId).QueryRows(&hostnameList)
	if err != nil {
		log.Debugf("WhenTempleteUnbind : %s", err.Error())
		return
	}
	if len(hostnameList) > 0 {
		hosts := fmt.Sprintf("('%s')", strings.Join(hostnameList, "','"))
		filterCond := fmt.Sprintf(" template_id = %d AND endpoint in %s", templateId, hosts)
		//set all related event cases status with "REMOVED"
		_, err = q.Raw(fmt.Sprintf("UPDATE event_cases SET status = '%s' WHERE %s ", removedStatus, filterCond)).Exec()
		if err != nil {
			return
		}
		affectedAlerms := []string{}
		_, err = q.Raw(fmt.Sprintf("SELECT id FROM event_cases WHERE %s ", filterCond)).QueryRows(&affectedAlerms)
		affectedRows = len(affectedAlerms)
		// AddNote("root", fmt.Sprintf("it's cause of templateId: %d has been unbind form hostgroupId: %d", templateId, hostgroupId), strings.Join(affectedAlerms, ","), processStatus, "")
		return
	}
	return
}

func WhenEndpointUnbind(hostId int, hostgroupId int) (err error, affectedRows int) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	tplids := []int{}
	//get bind templete list by hostgroupid
	_, err = q.Raw("SELECT tpl_id FROM grp_tpl WHERE grp_id = ?", hostgroupId).QueryRows(&tplids)
	if err != nil {
		return
	}
	log.Printf("sss %v", tplids)
	//get hostname by host id
	hostname := ""
	err = q.Raw("SELECT hostname FROM host WHERE id = ?", hostId).QueryRow(&hostname)
	if len(tplids) > 0 && hostname != "" {
		//set all related event cases status with "REMOVED"
		_, err = q.Raw("UPDATE event_cases SET status = ? WHERE template_id IN ? AND hostname = ?", removedStatus, tplids, hostname).Exec()
		if err != nil {
			return
		}
		affectedAlerms := []string{}
		_, err = q.Raw("SELECT id FROM event_cases WHERE template_id IN ? AND hostname = ?", tplids, hostname).QueryRows(affectedAlerms)
		affectedRows = len(affectedAlerms)
		// AddNote("root", fmt.Sprintf("it's cause of endpoint: %s has been unbind form hostgroupId: %d", hostname, hostgroupId), strings.Join(affectedAlerms, ","), processStatus, "")
		return
	}
	return
}
