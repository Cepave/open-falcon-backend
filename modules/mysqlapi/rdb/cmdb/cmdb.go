package cmdb

import (
	"github.com/jmoiron/sqlx"

	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	sqlxExt "github.com/Cepave/open-falcon-backend/common/db/sqlx"
	"github.com/Cepave/open-falcon-backend/common/utils"

	cmdbModel "github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
)

const (
	MAINTAIN_PERIOD_BEGIN = 946684800  // Sat, 01 Jan 2000 00:00:00 GMT
	MAINTAIN_PERIOD_END   = 4292329420 // Thu, 07 Jan 2106 17:43:40 GMT
)

type hostTuple struct {
	Hostname      string `db:"hostname"`
	Ip            string `db:"ip"`
	MaintainBegin uint32 `db:"maintain_begin"`
	MaintainEnd   uint32 `db:"maintain_end"`
}

type syncHostTx struct {
	hosts []*hostTuple
}

type syncHostGroupTx struct {
	groups []*cmdbModel.SyncHostGroup
}

type syncHostgroupContaining struct {
	relations map[string][]string
}

func (syncTx *syncHostgroupContaining) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	// transform string string relation into temporary memory table
	tx.MustExec(`DROP TEMPORARY TABLE IF EXISTS tmp_hostgroup_host`)

	tx.MustExec(`
		CREATE TEMPORARY TABLE tmp_hostgroup_host (
			hh_group_name varchar(255) NOT NULL,
			hh_hostname varchar(255) NOT NULL,
			PRIMARY KEY(hh_group_name, hh_hostname)
		) DEFAULT CHARSET=utf8
    `)
	txExt := sqlxExt.ToTxExt(tx)

	batchStmt := txExt.Preparex(
		`
		INSERT INTO tmp_hostgroup_host(hh_group_name, hh_hostname)
		VALUES (?, ?), (?, ?), (?, ?), (?, ?),
			(?, ?), (?, ?), (?, ?), (?, ?),
			(?, ?), (?, ?), (?, ?), (?, ?)
		`,
	)
	restStmt := txExt.Preparex(
		`
		INSERT INTO tmp_hostgroup_host(hh_group_name, hh_hostname)
		VALUES (?, ?)
		`,
	)
	for groupName, hosts := range syncTx.relations {
		utils.MakeAbstractArray(hosts).BatchProcess(
			12,
			func(batch interface{}) {
				batchStmt.Exec(
					utils.FlattenToSlice(
						batch,
						func(v interface{}) []interface{} {
							return []interface{}{
								groupName, v,
							}
						},
					)...,
				)
			},
			func(rest interface{}) {
				for _, hostName := range rest.([]string) {
					restStmt.MustExec(
						groupName, hostName,
					)
				}
			},
		)
	}

	/**
	 * Removes hosts which are not appeared in BOSS
	 */
	tx.MustExec(`
		DELETE gh
		FROM grp_host AS gh
			INNER JOIN
			grp gp
			ON gh.grp_id = gp.id
				AND gp.come_from = 1
			INNER JOIN
			host AS hs
			ON hs.id = gh.host_id
			LEFT OUTER JOIN
			tmp_hostgroup_host AS hh
			ON hs.hostname = hh.hh_hostname
				AND gp.grp_name = hh.hh_group_name
		WHERE hh.hh_hostname IS NULL AND
			hh.hh_group_name IS NULL
	`)
	// :~)
	/**
	 * Adds hosts which are appeared in BOSS
	 */
	tx.MustExec(`
		INSERT INTO grp_host (grp_id, host_id)
		SELECT gp.id, hs.id
		FROM tmp_hostgroup_host AS hh
			INNER JOIN
			grp gp
			ON hh.hh_group_name = gp.grp_name
			INNER JOIN
			host AS hs
			ON hh.hh_hostname = hs.hostname
			LEFT OUTER JOIN
			grp_host AS gh
			ON hs.id = gh.host_id
				AND gp.id = gh.grp_id
		WHERE gh.host_id IS NULL
			AND gh.grp_id IS NULL
	`)
	// :~)

	tx.MustExec(`DROP TEMPORARY TABLE tmp_hostgroup_host`)
	return commonDb.TxCommit
}

func (syncTx *syncHostGroupTx) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	tx.MustExec(`DROP TEMPORARY TABLE IF EXISTS tempgrp`)
	tx.MustExec(`
		CREATE TEMPORARY TABLE tempgrp (
			grp_name varchar(255) NOT NULL,
			create_user varchar(64) NOT NULL,
			come_from tinyint(4) NOT NULL,
			UNIQUE KEY idx_host_grp_grp_name (grp_name)
		) DEFAULT CHARSET=utf8
    `)
	/*
	 *  multiple insertion with prepared statement
	 *  come_from = 1
	 *  create_user = ?
	 */
	txExt := sqlxExt.ToTxExt(tx)

	batchStmt := txExt.Preparex(
		`
		INSERT INTO tempgrp (grp_name, create_user, come_from)
		VALUES (?, ?, 1), (?, ?, 1), (?, ?, 1), (?, ?, 1),
			 (?, ?, 1), (?, ?, 1), (?, ?, 1), (?, ?, 1),
			 (?, ?, 1), (?, ?, 1), (?, ?, 1), (?, ?, 1)
		`,
	)
	restStmt := txExt.PrepareNamed(
		`
		INSERT INTO tempgrp (grp_name, create_user, come_from)
		VALUES (:name, :creator, 1)
		`,
	)
	utils.MakeAbstractArray(syncTx.groups).BatchProcess(
		12,
		func(batch interface{}) {
			batchStmt.MustExec(
				utils.FlattenToSlice(
					batch,
					func(group interface{}) []interface{} {
						groupData := group.(*cmdbModel.SyncHostGroup)
						return []interface{}{
							groupData.Name, groupData.Creator,
						}
					},
				)...,
			)
		},
		func(rest interface{}) {
			for _, group := range rest.([]*cmdbModel.SyncHostGroup) {
				restStmt.MustExec(group)
			}
		},
	)
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
			Hostname:      h.Name,
			Ip:            h.IP,
			MaintainBegin: begin,
			MaintainEnd:   end,
		})
	}
	return dbData
}

func (syncTx *syncHostTx) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	tx.MustExec(`DROP TEMPORARY TABLE IF EXISTS temp_host`)
	tx.MustExec(`
		CREATE TEMPORARY TABLE temp_host (
			hostname varchar(255) NOT NULL PRIMARY KEY,
			ip varchar(16) NOT NULL,
			maintain_begin int(10) unsigned NOT NULL,
			maintain_end int(10) unsigned NOT NULL
    	) DEFAULT CHARSET=utf8
     `)
	/*
	 *  multiple insertion with prepared statement
	 */
	txExt := sqlxExt.ToTxExt(tx)

	batchStmt := txExt.Preparex(
		`
		INSERT INTO temp_host (hostname, ip, maintain_begin, maintain_end)
		VALUES (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?),
			 (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?)
		`,
	)
	restStmt := txExt.PrepareNamed(
		`
		INSERT INTO temp_host (hostname, ip, maintain_begin, maintain_end)
		VALUES (:hostname, :ip, :maintain_begin, :maintain_end)
		`,
	)
	utils.MakeAbstractArray(syncTx.hosts).BatchProcess(
		8,
		func(batch interface{}) {
			batchStmt.Exec(
				utils.FlattenToSlice(
					batch,
					func(host interface{}) []interface{} {
						hostData := host.(*hostTuple)
						return []interface{}{
							hostData.Hostname, hostData.Ip,
							hostData.MaintainBegin, hostData.MaintainEnd,
						}
					},
				)...,
			)
		},
		func(rest interface{}) {
			for _, host := range rest.([]*hostTuple) {
				restStmt.MustExec(host)
			}
		},
	)
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

	txProcessorRel := &syncHostgroupContaining{
		relations: syncData.Relations,
	}

	DbFacade.NewSqlxDbCtrl().InTx(txProcessorRel)
}
