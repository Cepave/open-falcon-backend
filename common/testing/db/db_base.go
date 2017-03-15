package db

import (
	"flag"

	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	f "github.com/Cepave/open-falcon-backend/common/db/facade"
	check "gopkg.in/check.v1"
)

// The base environment for RDB testing

var dsnMysql = flag.String("dsn_mysql", "", "DSN of MySql")

var DeleteNqmAgentPingtaskSQL = `DELETE FROM nqm_agent_ping_task WHERE apt_ag_id >= 24021 AND apt_ag_id <= 24025`

var InsertPingtaskSQL = `
INSERT INTO nqm_ping_task(pt_id, pt_name, pt_period)
VALUES(10119, 'test1-pingtask_name', 40),
(10120, 'test2-pingtask_name', 3)
`
var DeletePingtaskSQL = `DELETE FROM nqm_ping_task WHERE pt_id >= 10119`

var InsertHostSQL = `
INSERT INTO host(id, hostname, agent_version, plugin_version)
VALUES(36091, 'ct-agent-1', '', ''),
	(36092, 'ct-agent-2', '', ''),
	(36093, 'ct-agent-3', '', ''),
	(36094, 'ct-agent-4', '', '')
`
var DeleteHostSQL = `DELETE FROM host WHERE id >= 36091 AND id <= 36095`

var InsertNqmAgentSQL = `
INSERT INTO nqm_agent(
	ag_id, ag_hs_id, ag_name, ag_connection_id, ag_hostname, ag_ip_address,
	ag_pv_id, ag_ct_id
)
VALUES(24021, 36091, 'ct-255-1', 'ct-255-1@201.3.116.1', 'ct-1', x'C9037401', 7, 255),
	(24022, 36092, 'ct-255-2', 'ct-255-2@201.3.116.2', 'ct-2', x'C9037402', 7, 255),
	(24023, 36093, 'ct-255-3', 'ct-255-3@201.4.23.3', 'ct-3', x'C9037403', 7, 255),
	(24024, 36094, 'ct-263-1', 'ct-63-1@201.77.23.3', 'ct-4', x'C9022403', 7, 263)
`
var DeleteNqmAgentSQL = `DELETE FROM nqm_agent WHERE ag_id >= 24021 AND ag_id <= 24025`

var InitNqmAgentAndPingtaskSQL = []string{InsertPingtaskSQL, InsertHostSQL, InsertNqmAgentSQL}
var CleanNqmAgentAndPingtaskSQL = []string{DeleteNqmAgentPingtaskSQL, DeleteNqmAgentSQL, DeleteHostSQL, DeletePingtaskSQL}

// This callback is used to setup a viable database configuration while testing
type ViableDbConfigFunc func(config *commonDb.DbConfig)

// This function is used to:
//
// 1) Check whether or not the configuration o "dsn_mysql" has been supplied
// 2) If it does, supply the data of configuration to callback function
func SetupByViableDbConfig(c *check.C, configFunc ViableDbConfigFunc) bool {
	config := GetDbConfig(c)

	if config != nil {
		configFunc(config)
	}

	return config != nil
}

func GetDbConfig(c *check.C) *commonDb.DbConfig {
	if *dsnMysql == "" {
		c.Skip("Skip database testing. Needs \"-dsn_mysql=<MySQL DSN>\"")
		return nil
	}

	return &commonDb.DbConfig{
		Dsn:     *dsnMysql,
		MaxIdle: 2,
	}
}
func InitDbFacade(c *check.C) *f.DbFacade {
	var dbFacade = &f.DbFacade{}
	dbConfig := GetDbConfig(c)

	if dbConfig == nil {
		return nil
	}

	err := dbFacade.Open(dbConfig)
	c.Assert(err, check.IsNil)

	return dbFacade
}
func ReleaseDbFacade(c *check.C, dbFacade *f.DbFacade) {
	if dbFacade != nil {
		dbFacade.Release()
	}
}

// Checks whether or not skipping testing by viable arguments
func HasDbEnvForMysqlOrSkip(c *check.C) bool {
	var hasMySqlDsn = *dsnMysql != ""

	if !hasMySqlDsn {
		c.Skip("Skip Mysql Test: -dsn_mysql=<dsn>")
	}

	return hasMySqlDsn
}
