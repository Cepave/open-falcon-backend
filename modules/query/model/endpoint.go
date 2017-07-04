package model

import (
	"time"

	"fmt"

	"github.com/Cepave/open-falcon-backend/modules/query/database"
	"github.com/Cepave/open-falcon-backend/modules/query/g"
	log "github.com/sirupsen/logrus"
)

type Endpoint struct {
	Id       int64     `json:"id"`
	Endpoint string    `json:"endpoint"`
	Ts       int64     `json:"ts"`
	TCreate  time.Time `json:"-"`
	TModify  time.Time `json:"-"`
	Ipv4     string    `json:"-"`
}

type EndpointCounter struct {
	Id         int64  `json:"id"`
	EndpointId int    `orm:"endpoint_id";json:"endpoint_id"`
	Counter    string `orm:"counter";json:"counter"`
	Step       int
	Type       string
	Ts         int64     `json:"ts"`
	TCreate    time.Time `json:"-"`
	TModify    time.Time `json:"-"`
}

func EndpointQuery(query string) (endpointList []string) {
	db := database.DBConn()
	gconf := g.Config()
	var enps []Endpoint
	var sqlStr string
	if query == "" {
		sqlStr = "SELECT * FROM graph.endpoint "
	} else {
		sqlStr = fmt.Sprintf("SELECT * FROM graph.endpoint WHERE endpoint regexp '%s'", query)
	}
	if g.Config().GraphDB.Limit != -1 {
		sqlStr = fmt.Sprintf("%s limit %v", sqlStr, gconf.GraphDB.Limit)
	}
	log.Debugf("endpoint query: %s", sqlStr)
	db.Raw(sqlStr).Scan(&enps)
	if len(enps) != 0 {
		for _, host := range enps {
			endpointList = append(endpointList, host.Endpoint)
		}
	}
	return
}

func EndpointIdQuery(endpoints []string) (endpointList []int64) {
	db := database.DBConn()
	percooke := ""
	for indx, e := range endpoints {
		if indx == 0 {
			percooke = e
		} else {
			percooke = fmt.Sprintf("%s\",\"%s", percooke, e)
		}
	}
	percooke = fmt.Sprintf("(\"%s\")", percooke)
	sqlstr := fmt.Sprintf("select id from graph.endpoint where endpoint in %s", percooke)
	log.Debugf("find EndpointIdQuery sql: %s", sqlstr)
	endId := []Endpoint{}
	db.Raw(sqlstr).Scan(&endId)
	for _, s := range endId {
		endpointList = append(endpointList, s.Id)
	}
	log.Debugf("EndpointIdQuery result: %v", endpointList)
	return
}

func FindMatchedCounters(endpointList []int64, counter string) (result []string) {
	db := database.DBConn()
	percooke := ""
	for indx, e := range endpointList {
		if indx == 0 {
			percooke = fmt.Sprintf("%v", e)
		} else {
			percooke = fmt.Sprintf("%s,%d", percooke, e)
		}
	}
	percooke = fmt.Sprintf("(%s)", percooke)
	enpc := []EndpointCounter{}
	sqlstr := fmt.Sprintf("SELECT distinct(counter) FROM graph.endpoint_counter WHERE endpoint_id in %s and counter like '%s'", percooke, counter)
	log.Debugf("find matched counters sql: %s", sqlstr)
	db.Raw(sqlstr).Scan(&enpc)
	for _, c := range enpc {
		result = append(result, c.Counter)
	}
	return
}
