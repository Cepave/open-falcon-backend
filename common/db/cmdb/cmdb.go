package cmdb

import (
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	sqlxExt "github.com/Cepave/open-falcon-backend/common/db/sqlx"
	cmdbModel "github.com/Cepave/open-falcon-backend/common/model/cmdb"
	"github.com/jmoiron/sqlx"
)

type syncHostTx struct {
	hosts []cmdbModel.SyncHost
	err   error
}

type hostTuple struct {
	Hostname       string
	Ip             string
	Maintain_begin uint32
	Maintain_end   uint32
}

type syncHostGroupTx struct {
	groups []cmdbModel.SyncHostGroup
	err    error
}

func (syncTx *syncHostGroupTx) InTx(tx *sqlx.Tx) commonDb.TxFinale {
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
		id int(11) NOT NULL AUTO_INCREMENT,
     	hostname varchar(255) NOT NULL DEFAULT '',
     	ip varchar(16) NOT NULL DEFAULT '',
     	agent_version varchar(16) NOT NULL DEFAULT '',
     	plugin_version varchar(128) NOT NULL DEFAULT '',
     	maintain_begin int(10) unsigned NOT NULL DEFAULT '0',
     	maintain_end int(10) unsigned NOT NULL DEFAULT '0',
     	update_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
     	resource_object_id int(10) DEFAULT NULL,
     	PRIMARY KEY (id),
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
	hostsData := api2tuple(syncTx.hosts)
	for _, s := range hostsData {
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
func StartSync(syncData *cmdbModel.SyncForAdding) (*cmdbModel.SyncItem, error) {
	/*
	 * Checks if this sync can be carried on
	 * Required parameters should be name, function, and timeout
	 */
	/*
		item, err := chemin.checkDoableSync("cmdbSync")
		if err != nil {
			return nil, err
		}
	*/

	// sync Hosts

	txProcessorHost := &syncHostTx{
		hosts: syncData.Hosts,
	}

	DbFacade.NewSqlxDbCtrl().InTx(txProcessorHost)

	if txProcessorHost.err != nil {
		return nil, txProcessorHost.err
	}

	// sync HostGroups

	txProcessorGroup := &syncHostGroupTx{
		groups: syncData.Hostgroups,
	}

	DbFacade.NewSqlxDbCtrl().InTx(txProcessorGroup)

	if txProcessorGroup.err != nil {
		return nil, txProcessorGroup.err
	}
	// wait to fix
	var item *cmdbModel.SyncItem
	item = nil
	return item, nil
}
