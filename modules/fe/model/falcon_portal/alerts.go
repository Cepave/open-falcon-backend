package falconPortal

import (
	"strconv"
	"strings"

	"github.com/Cepave/open-falcon-backend/modules/fe/model/boss"
	log "github.com/Sirupsen/logrus"
	"github.com/emirpasic/gods/sets/hashset"
)

func AlertsConvert(result []EventCases) (resp []AlertsResp, endpointSet *hashset.Set, err error) {
	endpointSet = hashset.New()
	resp = []AlertsResp{}
	const notFound = "not found"
	for _, item := range result {
		recordOne := AlertsResp{}
		priority := strconv.Itoa(item.Priority)
		tlayout := "2006-01-02 15:04"
		sTime := item.Timestamp.Format(tlayout)
		eTime := item.UpdateAt.Format(tlayout)
		recordOne.Hash = item.Id
		recordOne.HostName = item.Endpoint
		recordOne.Metric = item.Metric
		recordOne.Author = item.TplCreator
		recordOne.TemplateID = item.TemplateId
		recordOne.Priority = priority
		recordOne.Severity = getSeverity(priority)
		recordOne.Status = getStatus(item.Status)
		recordOne.StatusRaw = item.Status
		recordOne.Type = strings.Split(item.Metric, ".")[0]
		recordOne.Content = item.Note
		recordOne.TimeStart = sTime
		recordOne.TimeUpdate = eTime
		recordOne.Duration = getDuration(eTime)
		// recordOne.Notes = getNote(item.Id, sTime)
		recordOne.Notes = []map[string]string{}
		recordOne.Events = item.Events
		recordOne.Process = item.ProcessStatus
		recordOne.Function = item.Func
		recordOne.Condition = item.Cond
		recordOne.StepLimit = item.MaxStep
		recordOne.Step = item.CurrentStep
		recordOne.Contact = []boss.Contactor{boss.Contactor{"-", "-", "-"}}
		recordOne.Platform = notFound
		recordOne.IP = notFound
		recordOne.IDC = notFound
		// ///make compatible for overall , if need it, please uncomment below
		// recordOne.ID = item.Id
		// recordOne.Endpoint = item.Endpoint
		// recordOne.Func = item.Func
		// recordOne.Cond = item.Cond
		// recordOne.Note = item.Note
		// recordOne.MaxStep = item.MaxStep
		// recordOne.CurrentStep = item.CurrentStep
		// Status     string `json:"status"` //conflict
		// recordOne.ProcessStatus = item.ProcessStatus
		// recordOne.StartAt = item.Timestamp
		// recordOne.UpdateAt = item.UpdateAt
		resp = append(resp, recordOne)
		endpointSet.Add(item.Endpoint)
	}
	return
}

func GetAlertInfo(resp []AlertsResp, endpointList *hashset.Set, showAll bool) (respComplete []AlertsResp) {
	respComplete = resp
	platformInfo, err := boss.GetPlatformASJSON()
	if err != nil {
		log.Error("query platform failed, please check boss api status")
		return
	}
	// //get ip mapping , platformList, popIds
	ipMapping, platList, popIDs := boss.GenPlatMap(platformInfo, endpointList)
	log.Debugf("ipMapping: %v, platList: %v, popIDs: %v", ipMapping.Values(), platList.Values(), popIDs)
	if ipMapping.Size() == 0 {
		log.Warnf("can not found any platform info that matched alarms, current alarms case got: %d", len(resp))
		return
	}
	//get contact mapping
	contactInfo, err := boss.GetPlatfromContactInfo(platList)
	if err != nil {
		log.Error("query contact info failed, please check boss api status")
	}
	//get idc mapping
	idcMapping, err := boss.IdcMapping(popIDs)
	if err != nil {
		log.Error("query idc info failed, please check boss api status")
	}
	log.Debugf("contactInfo: %v, idcMapping: %v", contactInfo, idcMapping)
	respCompleteTmp := []AlertsResp{}
	for _, item := range resp {
		ipInfotmp, ok := ipMapping.Get(item.HostName)
		if !ok {
			log.Debugf("item.HostName: is missing", item.HostName)
			if showAll {
				respCompleteTmp = append(respCompleteTmp, item)
			} else {
				//if not found any ip mapping of this host on boss will skip it
				continue
			}
		}
		ipInfo := ipInfotmp.(boss.IPInfo)
		item.Platform = ipInfo.Platform
		if ipInfo.IPStatus == "0" {
			if showAll {
				item.IP = ipInfo.IP + "(deactivated)"
			} else {
				// if this host is set to inactive, will skip it
				continue
			}
		} else {
			item.IP = ipInfo.IP
		}
		if contact, ok := contactInfo[item.Platform]; ok {
			log.Debugf("item.Hsotname: got contact -> %s", item.HostName)
			item.Contact = contact
		}
		popid, _ := strconv.Atoi(ipInfo.POPID)
		name, ok := idcMapping[popid]
		if !ok {
			respCompleteTmp = append(respCompleteTmp, item)
			continue
		}
		item.IDC = name
		respCompleteTmp = append(respCompleteTmp, item)
	}
	respComplete = respCompleteTmp
	return
}

func GetAlertsNotes(alerts []AlertsResp) (alertsResp []AlertsResp) {
	alertsResp = []AlertsResp{}
	for _, record := range alerts {
		record.Notes = getNote(record.Hash, record.TimeStart)
		alertsResp = append(alertsResp, record)
	}
	return
}
