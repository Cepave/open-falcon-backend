package host

import (
	"errors"
	"fmt"
	"strconv"

	h "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/helper"
	f "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/falcon_portal"
	u "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func GetHostBindToWhichHostGroup(c *gin.Context) {
	HostIdTmp := c.Params.ByName("host_id")
	if HostIdTmp == "" {
		h.JSONR(c, badstatus, "host id is missing")
		return
	}
	hostID, err := strconv.Atoi(HostIdTmp)
	if err != nil {
		log.Debugf("HostId: %v", HostIdTmp)
		h.JSONR(c, badstatus, err)
		return
	}
	grpHostMap := []f.GrpHost{}
	db.Falcon.Select("grp_id").Where("host_id = ?", hostID).Find(&grpHostMap)
	grpIds := []int64{}
	for _, g := range grpHostMap {
		grpIds = append(grpIds, g.GrpID)
	}
	hostgroups := []f.HostGroup{}
	if len(grpIds) != 0 {
		grpIdsStr, _ := u.ArrInt64ToString(grpIds)
		db.Falcon.Where(fmt.Sprintf("id in (%s)", grpIdsStr)).Find(&hostgroups)
	}
	h.JSONR(c, hostgroups)
	return
}

func GetHostGroupWithTemplate(c *gin.Context) {
	grpIDtmp := c.Params.ByName("host_group")
	if grpIDtmp == "" {
		h.JSONR(c, badstatus, "grp id is missing")
		return
	}
	grpID, err := strconv.Atoi(grpIDtmp)
	if err != nil {
		log.Debugf("grpIDtmp: %v", grpIDtmp)
		h.JSONR(c, badstatus, err)
		return
	}
	hostgroup := f.HostGroup{ID: int64(grpID)}
	if dt := db.Falcon.Find(&hostgroup); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	hosts := []f.Host{}
	grpHosts := []f.GrpHost{}
	if dt := db.Falcon.Where("grp_id = ?", grpID).Find(&grpHosts); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	for _, grph := range grpHosts {
		var host f.Host
		db.Falcon.Find(&host, grph.HostID)
		if host.ID != 0 {
			hosts = append(hosts, host)
		}
	}
	h.JSONR(c, map[string]interface{}{
		"hostgroup": hostgroup,
		"hosts":     hosts,
	})
	return
}

func GetGrpsRelatedHost(c *gin.Context) {
	hostIDtmp := c.Params.ByName("host_id")
	if hostIDtmp == "" {
		h.JSONR(c, badstatus, "host id is missing")
		return
	}
	hostID, err := strconv.Atoi(hostIDtmp)
	if err != nil {
		log.Debugf("host id: %v", hostIDtmp)
		h.JSONR(c, badstatus, err)
		return
	}

	host := f.Host{ID: int64(hostID)}
	if dt := db.Falcon.Find(&host); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	grps := host.RelatedGrp()
	h.JSONR(c, grps)
	return
}

func GetTplsRelatedHost(c *gin.Context) {
	hostIDtmp := c.Params.ByName("host_id")
	if hostIDtmp == "" {
		h.JSONR(c, badstatus, "host id is missing")
		return
	}
	hostID, err := strconv.Atoi(hostIDtmp)
	if err != nil {
		log.Debugf("host id: %v", hostIDtmp)
		h.JSONR(c, badstatus, err)
		return
	}
	host := f.Host{ID: int64(hostID)}
	if dt := db.Falcon.Find(&host); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	tpls := host.RelatedTpl()
	h.JSONR(c, tpls)
	return
}

type APIHostsSetToMaintainInputs struct {
	Hosts     []string `json:"hosts" form:"hosts" binding:"required"`
	StartTime int32    `json:"start_time" form:"start_time" binding:"gt=-1"`
	EndTime   int32    `json:"end_time" form:"end_time" binding:"gt=-1"`
}

func (mine APIHostsSetToMaintainInputs) Check() (err error) {
	if mine.StartTime == 0 && mine.EndTime == 0 {
		return
	}
	if mine.StartTime >= mine.EndTime {
		err = errors.New("start_time can not greater than end_time")
		return
	}
	return
}

func HostsSetToMaintain(c *gin.Context) {
	inputs := APIHostsSetToMaintainInputs{EndTime: -1, StartTime: -1}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	if err := inputs.Check(); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	fhosts := []f.Host{}
	tx := db.Falcon.Begin()
	dtmp := tx.Model(&fhosts).Where("hostname IN (?)", inputs.Hosts).UpdateColumn(map[string]int32{
		"maintain_begin": inputs.StartTime,
		"maintain_end":   inputs.EndTime,
	})
	if dtmp.Error != nil {
		tx.Rollback()
		h.JSONR(c, badstatus, dtmp.Error)
		return
	}
	tx.Commit()
	db.Falcon.Model(&fhosts).Where("hostname IN (?)", inputs.Hosts).Scan(&fhosts)
	h.JSONR(c, fhosts)
	return
}
