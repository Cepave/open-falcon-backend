package cmdb

import (
	"github.com/jmoiron/sqlx"

	cmdbModel "github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"

	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	sqlxExt "github.com/Cepave/open-falcon-backend/common/db/sqlx"
)

const (
	MAINTAIN_PERIOD_BEGIN = 946684800  // Sat, 01 Jan 2000 00:00:00 GMT
	MAINTAIN_PERIOD_END   = 4292329420 // Thu, 07 Jan 2106 17:43:40 GMT
)

type hostTuple struct {
	Hostname       string
	Ip             string
	Maintain_begin uint32
	Maintain_end   uint32
}

type syncHostTx struct {
	hosts []*hostTuple
}

type syncHostGroupTx struct {
	groups []*cmdbModel.SyncHostGroup
}

type syncRelTx struct {
	relations map[string][]string
}

func (syncTx *syncRelTx) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	// delete all the grp relation that come_from = 1
	tx.MustExec(`
		DELETE FROM grp_host WHERE grp_id IN
		(SELECT id FROM grp WHERE come_from = 1)
	`)
	// transform string string relation into temporary memory table
	tx.MustExec(`DROP TEMPORARY TABLE IF EXISTS trel`)

	tx.MustExec(`
		CREATE TEMPORARY TABLE trel (
			grp_name varchar(255) NOT NULL DEFAULT '',
			hostname varchar(255) NOT NULL DEFAULT ''
		) ENGINE=MEMORY DEFAULT CHARSET=utf8
    `)
	txExt := sqlxExt.ToTxExt(tx)
	namedStmt := txExt.PrepareNamed(`
		INSERT INTO trel(grp_name, hostname)
		VALUES (:gname, :hname)
	`)
	for key, val := range syncTx.relations {
		for _, hname := range val {
			m := map[string]interface{}{
				"gname": key,
				"hname": hname,
			}
			namedStmt.MustExec(m)
		}
	}

	// insert into grp_host
	tx.MustExec(`
		INSERT INTO grp_host (grp_id, host_id)
		SELECT grp.id, host.id FROM grp, trel, host
			WHERE grp.grp_name = trel.grp_name
			AND host.hostname = trel.hostname
	`)
	tx.MustExec(`DROP TEMPORARY TABLE trel`)
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

func api2tuple(hosts []*cmdbModel.SyncHost) []*hostTuple {
	var begin uint32
	var end uint32
	dbData := []*hostTuple{}
	for _, h := range hosts {
		if h.Activate == 1 {
			begin = uint32(0)
			end = uint32(0)
		} else {
			begin = MAINTAIN_PERIOD_BEGIN // Sat, 01 Jan 2000 00:00:00 GMT
			end = MAINTAIN_PERIOD_END     // Thu, 07 Jan 2106 17:43:40 GMT
		}
		dbData = append(dbData, &hostTuple{
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
func SyncForHosts(syncData *cmdbModel.SyncForAdding) {
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
