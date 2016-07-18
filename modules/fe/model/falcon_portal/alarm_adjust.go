package falconPortal

import (
	"errors"
	"fmt"
	"strings"

	"time"

	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego/orm"
)

const (
	removedStatus  = "REMOVED"
	modifiedStatus = "UNKNOWN"
)

func UpdateCloseNote(eventCaseID []string, closedNote string) error {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	layoutf := "2006-01-02 15:04:05"
	now := time.Now().Format(layoutf)
	var err error
	for _, cid := range eventCaseID {
		_, err = q.Raw("UPDATE event_cases SET closed_note = ?, closed_at = ? WHERE id = ?", closedNote, now, cid).Exec()
		if err != nil {
			return err
		}
	}
	return err
}

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
	UpdateCloseNote(affectedAlerms, fmt.Sprintf("Because of strategyId: %d has been modified by user", strategyId))
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
	UpdateCloseNote(affectedAlerms, fmt.Sprintf("Because of strategyId: %d has been deleted by user", strategyId))
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
	UpdateCloseNote(affectedAlerms, fmt.Sprintf("Because of templateId: %d has been deleted by user", templateId))
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
		UpdateCloseNote(affectedAlerms, fmt.Sprintf("Because of templateId: %d has been unbind form hostgroupId: %d", templateId, hostgroupId))
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
	//get hostname by host id
	hostname := ""
	err = q.Raw("SELECT hostname FROM host WHERE id = ?", hostId).QueryRow(&hostname)
	if err != nil {
		log.Debug(err.Error())
	}
	log.Debug("host: %s", hostname)
	if len(tplids) > 0 && hostname != "" {
		log.Debug("will disable hostname: %s", hostname)
		templetIds := "("
		for ind, mid := range tplids {
			if ind == 0 {
				templetIds = fmt.Sprintf("(%s", strconv.Itoa(mid))
			} else {
				templetIds = fmt.Sprintf("%s,%s", templetIds, strconv.Itoa(mid))
			}
		}
		templetIds = fmt.Sprintf("%s)", templetIds)
		//set all related event cases status with "REMOVED"
		_, err = q.Raw(fmt.Sprintf("UPDATE event_cases SET status = '%s' WHERE template_id IN %s AND endpoint = '%s'", removedStatus, templetIds, hostname)).Exec()
		if err != nil {
			log.Debug(err.Error())
			return
		}
		affectedAlerms := []string{}
		_, err = q.Raw(fmt.Sprintf("SELECT id FROM event_cases WHERE template_id IN %s AND endpoint = '%s'", templetIds, hostname)).QueryRows(&affectedAlerms)
		if err != nil {
			log.Debug(err.Error())
			return
		}
		affectedRows = len(affectedAlerms)
		UpdateCloseNote(affectedAlerms, fmt.Sprintf("Because of endpoint: %s has been unbind form hostgroupId: %d", hostname, hostgroupId))
		return
	}
	log.Debug("addected alerm: %v", affectedRows)
	return
}

func WhenEndpointOnMaintain(hostId int, maintainBegin int64, maintainEnd int64) (err error, affectedRows int) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	var host Host
	if hostId != 0 {
		err = q.Raw("SELECT * FROM host WHERE id = ?", hostId).QueryRow(&host)
		log.Debugf("host: %v", host)
		if err != nil {
			return
		}
		if maintainBegin <= time.Now().Unix() && maintainEnd >= time.Now().Unix() {
			_, err = q.Raw("UPDATE event_cases SET status = ? WHERE endpoint = ?", modifiedStatus, host.Hostname).Exec()
		} else {
			err = errors.New("MaintainBegin / MaintainEnd is wrong please check it!")
			log.Debugf("%v , %v, %v", time.Unix(maintainBegin, 10), time.Unix(maintainEnd, 10), err.Error())
			return
		}
	} else {
		return
	}
	affectedAlerms := []string{}
	_, err = q.Raw("SELECT id FROM event_cases WHERE endpoint = ? ", host.Hostname).QueryRows(&affectedAlerms)
	if err != nil {
		log.Debug(err.Error())
	}
	affectedRows = len(affectedAlerms)
	UpdateCloseNote(affectedAlerms, fmt.Sprintf("Because of endpoint: %s is under maintenance", host.Hostname))
	return
}
