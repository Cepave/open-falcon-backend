package cmdb

import (
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	sqlxExt "github.com/Cepave/open-falcon-backend/common/db/sqlx"
	cmdbModel "github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	"github.com/jmoiron/sqlx"
)

type hostTuple struct {
	Hostname       string
	Ip             string
	Maintain_begin uint32
	Maintain_end   uint32
}

type syncHostTx struct {
	hosts []hostTuple
}

type syncHostGroupTx struct {
	groups []cmdbModel.SyncHostGroup
}

type syncRelTx struct {
	relations map[string][]string
}

type relHostTuple struct {
	Id       int
	Hostname string
}

func hostTuples2map(tuples []relHostTuple) map[string]int {
	result := make(map[string]int)
	for _, t := range tuples {
		result[t.Hostname] = t.Id
	}
	return result
}

type relGroupTuple struct {
	Id       int
	Grp_name string
}

func groupTuples2map(tuples []relGroupTuple) map[string]int {
	result := make(map[string]int)
	for _, t := range tuples {
		result[t.Grp_name] = t.Id
	}
	return result
}

func (syncTx *syncRelTx) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	// prepare id, name mapping data
	h := []relHostTuple{}
	tx.Select(&h, "id, hostname from host")
	//hd := hostTuples2map(h)
	g := []relGroupTuple{}
	tx.Select(&g, "id, grp_name from grp")
	//gd := groupTuples2map(g)

	/*transform string 2 string relation to id 2 id relation
	for _, r := syncTx.relations {

	}*/
	return commonDb.TxCommit
}

func (syncTx *syncHostGroupTx) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	tx.MustExec(`DROP TEMPORARY TABLE IF EXISTS tempgrp`)
	tx.MustExec(`
		CREATE TEMPORARY TABLE tempgrp (
		grp_name varchar(255) NOT NULL DEFAULT '',
		create_user varchar(64) NOT NULL DEFAULT '',
		come_from tinyint(4) NOT NULL DEFAULT '0',
		UNIQUE KEY idx_host_grp_grp_name (grp_name)
		) ENGINE=MEMORY DEFAULT CHARSET=utf8
    `)
	/*
	 *  multiple insertion with prepared statement
	 *  come_from = 1
	 *  create_user = ?
	 */
	txExt := sqlxExt.ToTxExt(tx)
	namedStmt := txExt.PrepareNamed(
		`INSERT INTO tempgrp (grp_name, create_user, come_from)
		 VALUES (:name, :creator, 1)`)
	for _, s := range syncTx.groups {
		namedStmt.MustExec(s)
	}
	// :~)

	/*
	 *  update host table from temp_host
	 */
	tx.MustExec(
		`
		UPDATE grp INNER JOIN tempgrp
		ON grp.grp_name = tempgrp.grp_name
		SET grp.create_user = tempgrp.create_user,
    		grp.come_from   = tempgrp.come_from
		`)
	/*
	 * insert host table from temp_host
	 */
	tx.MustExec(
		`
		INSERT INTO grp(grp_name, create_user, come_from)
		SELECT tempgrp.grp_name, tempgrp.create_user, tempgrp.come_from
		FROM tempgrp LEFT JOIN grp
		ON tempgrp.grp_name = grp.grp_name
		WHERE grp.grp_name IS NULL
		`)
	tx.MustExec(`DROP TEMPORARY TABLE tempgrp`)
	return commonDb.TxCommit
}

func api2tuple(hosts []cmdbModel.SyncHost) []hostTuple {
	var begin uint32
	var end uint32
	dbData := []hostTuple{}
	for _, h := range hosts {
		if h.Activate == 1 {
			begin = uint32(0)
			end = uint32(0)
		} else {
			begin = uint32(946684800) //  Sat, 01 Jan 2000 00:00:00 GMT
			end = uint32(4292329420)  // Thu, 07 Jan 2106 17:43:40 GMT
		}
		dbData = append(dbData, hostTuple{
			Hostname:       h.Name,
			Ip:             h.IP,
			Maintain_begin: begin,
			Maintain_end:   end,
		})
	}
	return dbData
}

func (syncTx *syncHostTx) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	tx.MustExec(`DROP TEMPORARY TABLE IF EXISTS temp_host`)
	tx.MustExec(`
		CREATE TEMPORARY TABLE temp_host (
     	hostname varchar(255) NOT NULL DEFAULT '',
     	ip varchar(16) NOT NULL DEFAULT '',
     	maintain_begin int(10) unsigned NOT NULL DEFAULT '0',
     	maintain_end int(10) unsigned NOT NULL DEFAULT '0',
     	UNIQUE KEY idx_host_hostname (hostname)
    	) ENGINE=MEMORY AUTO_INCREMENT=7 DEFAULT CHARSET=utf8
     `)
	/*
	 *  multiple insertion with prepared statement
	 */
	txExt := sqlxExt.ToTxExt(tx)
	namedStmt := txExt.PrepareNamed(
		`INSERT INTO temp_host (hostname, ip, maintain_begin, maintain_end)
		 VALUES (:hostname, :ip, :maintain_begin, :maintain_end)`)
	for _, s := range syncTx.hosts {
		namedStmt.MustExec(s)
	}
	// :~)

	/*
	 *  update host table from temp_host
	 */
	tx.MustExec(
		`
		UPDATE host INNER JOIN temp_host
		ON host.hostname = temp_host.hostname
		SET host.ip = temp_host.ip,
		    host.maintain_begin = temp_host.maintain_begin,
			host.maintain_end   = temp_host.maintain_end
		`)
	/*
	 * insert host table from temp_host
	 */
	tx.MustExec(
		`
		INSERT INTO host(hostname, ip, maintain_begin, maintain_end)
		SELECT temp_host.hostname, temp_host.ip,
		       temp_host.maintain_begin, temp_host.maintain_end
		FROM temp_host LEFT JOIN host
		ON temp_host.hostname = host.hostname
		WHERE host.hostname IS NULL
		`)
	tx.MustExec(`DROP TEMPORARY TABLE temp_host`)
	return commonDb.TxCommit
}

// Start the Synchronization of CMDB data.
func StartSync(syncData *cmdbModel.SyncForAdding) {
	// sync Hosts
	txProcessorHost := &syncHostTx{
		hosts: api2tuple(syncData.Hosts),
	}

	DbFacade.NewSqlxDbCtrl().InTx(txProcessorHost)

	// sync HostGroups

	txProcessorGroup := &syncHostGroupTx{
		groups: syncData.Hostgroups,
	}

	DbFacade.NewSqlxDbCtrl().InTx(txProcessorGroup)

	// sync Relations

	txProcessorRel := &syncRelTx{
		relations: syncData.Relations,
	}

	DbFacade.NewSqlxDbCtrl().InTx(txProcessorRel)
}
