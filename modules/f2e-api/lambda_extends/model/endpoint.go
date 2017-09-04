package model

import (
	"fmt"

	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/graph"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func EndpointQuery(query string) (endpointList []string) {
	db := config.Con()
	var enps []graph.Endpoint
	var sqlStr string
	if query == "" {
		sqlStr = "SELECT * FROM graph.endpoint "
	} else {
		sqlStr = fmt.Sprintf("SELECT * FROM graph.endpoint WHERE endpoint regexp '%s'", query)
	}
	qlimit := viper.GetInt("lambda_extends.qlimit")
	if !(qlimit <= 0) {
		sqlStr = fmt.Sprintf("%s limit %d", sqlStr, qlimit)
	}
	log.Debugf("endpoint query: %s", sqlStr)
	db.Graph.Raw(sqlStr).Scan(&enps)
	if len(enps) != 0 {
		for _, host := range enps {
			endpointList = append(endpointList, host.Endpoint)
		}
	}
	return
}

func BossEndpointQuery(platformName string) (endpointList []string) {
	db := config.Con()
	sqlStr := fmt.Sprintf("SELECT ge.id as id, ge.endpoint as endpoint FROM boss.hosts as ho JOIN graph.endpoint as ge ON ho.hostname = ge.endpoint WHERE ho.platform = '%s'", platformName)
	var enps []graph.Endpoint
	qlimit := viper.GetInt("lambda_extends.qlimit")
	if !(qlimit <= 0) {
		sqlStr = fmt.Sprintf("%s limit %d", sqlStr, qlimit)
	}
	log.Debugf("boss endpoint query: %s", sqlStr)
	db.Graph.Raw(sqlStr).Scan(&enps)
	if len(enps) != 0 {
		for _, host := range enps {
			endpointList = append(endpointList, host.Endpoint)
		}
	}
	return
}

func EndpointIdQuery(endpoints []string) (endpointList []uint) {
	db := config.Con()
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
	endId := []graph.Endpoint{}
	db.Graph.Raw(sqlstr).Scan(&endId)
	for _, s := range endId {
		endpointList = append(endpointList, s.ID)
	}
	log.Debugf("EndpointIdQuery result: %v", endpointList)
	return
}

func FindMatchedCounters(endpointList []uint, counter string) (result []string) {
	db := config.Con()
	percooke := ""
	for indx, e := range endpointList {
		if indx == 0 {
			percooke = fmt.Sprintf("%v", e)
		} else {
			percooke = fmt.Sprintf("%s,%d", percooke, e)
		}
	}
	percooke = fmt.Sprintf("(%s)", percooke)
	enpc := []graph.EndpointCounter{}
	sqlstr := fmt.Sprintf("SELECT distinct(counter) FROM graph.endpoint_counter WHERE endpoint_id in %s and counter like '%s'", percooke, counter)
	log.Debugf("find matched counters sql: %s", sqlstr)
	db.Graph.Raw(sqlstr).Scan(&enpc)
	for _, c := range enpc {
		result = append(result, c.Counter)
	}
	return
}
