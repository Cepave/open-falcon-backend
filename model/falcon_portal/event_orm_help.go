package falconPortal

import (
	"fmt"
	"log"

	"github.com/Cepave/fe/model/uic"
	"github.com/astaxie/beego/orm"
)

func getUserRole(username string) (int64, bool) {
	user := uic.ReadUserByName(username)
	// Role of root is 2. Role of admin assigned by root is 1.
	if user.Role == 2 || user.Role == 1 {
		return user.Id, true
	} else {
		return user.Id, false
	}
}

//get teamid from userid
func GetTeamIdsFromUser(uid int64) ([]int64, string, error) {
	tids := ""
	tidsI := []int64{}
	tidsI, err := uic.Tids(uid)
	if err != nil {
		log.Println(err.Error())
	} else if len(tidsI) != 0 {
		for _, v := range tidsI {
			if tids == "" {
				tids = string(v)
			} else {
				tids = fmt.Sprintf("%s,%d", tids, v)
			}
		}
	}
	return tidsI, tids, err
}

//get teamNames from teamids
func GetTeamNameFromTeamIds(teamIds []int64) ([]int64, string) {
	teamNames := ""
	teamNameI := []int64{}
	for _, v := range teamIds {
		iteam := uic.SelectTeamById(v)
		teamNameI = append(teamNameI, iteam.Id)
		if teamNames == "" {
			teamNames = fmt.Sprintf("\"%s\"", iteam.Name)
		} else {
			teamNames = fmt.Sprintf("%s,\"%s\"", teamNames, iteam.Name)
		}
	}
	return teamNameI, teamNames
}

//get actionIds from teamNames
func GetActionIdsFromTeamNames(teamNames string) ([]int, string, error) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	var actions []Action
	aids := ""
	aidI := []int{}
	_, err := q.Raw(fmt.Sprintf("select * from falcon_portal.action where uic IN(%s)", teamNames)).QueryRows(&actions)
	for _, v := range actions {
		aidI = append(aidI, v.Id)
		if aids == "" {
			aids = fmt.Sprintf("%d", v.Id)
		} else {
			aids = fmt.Sprintf("%s,%d", aids, v.Id)
		}
	}
	return aidI, aids, err
}

func GetTplIdFromActionId(aids string) ([]int, string, error) {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	var tpls []Tpl
	log.Printf("select * from falcon_portal.tpl where action_id IN(%s)", aids)
	_, err := q.Raw(fmt.Sprintf("select * from falcon_portal.tpl where action_id IN(%s)", aids)).QueryRows(&tpls)
	tplIds := ""
	tplIdI := []int{}
	for _, v := range tpls {
		tplIdI = append(tplIdI, v.Id)
		if tplIds == "" {
			tplIds = fmt.Sprintf("%d", v.Id)
		} else {
			tplIds = fmt.Sprintf("%s,%d", tplIds, v.Id)
		}
	}
	return tplIdI, tplIds, err
}

func GetCasePermission(username string) (isAdmin bool, tplIds string, err error) {
	var uid int64
	uid, isAdmin = getUserRole(username)
	if isAdmin {
		return
	}
	teamId, _, err := GetTeamIdsFromUser(uid)
	if err != nil || len(teamId) == 0 {
		return
	}
	_, teamNames := GetTeamNameFromTeamIds(teamId)
	if teamNames == "" {
		return
	}
	_, aids, err := GetActionIdsFromTeamNames(teamNames)
	if err != nil || aids == "" {
		return
	}
	_, tplIds, err = GetTplIdFromActionId(aids)
	return
}
