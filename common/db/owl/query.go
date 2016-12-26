package owl

import (
	"database/sql"
	"github.com/satori/go.uuid"
	model "github.com/Cepave/open-falcon-backend/common/model/owl"
	"github.com/Cepave/open-falcon-backend/common/db"
	sqlxExt "github.com/Cepave/open-falcon-backend/common/db/sqlx"
	"github.com/jmoiron/sqlx"
	"time"
)

// Adds a query or refresh an existing one
func AddOrRefreshQuery(query *model.Query, accessTime time.Time) {
	DbFacade.SqlxDbCtrl.InTx(&addOrRefreshQueryTx{ query, accessTime })
}

// Updates access time of a query or adding new one
func UpdateAccessTimeOrAddNewOne(query *model.Query, accessTime time.Time) {
	result := db.ToResultExt(DbFacade.SqlxDbCtrl.NamedExec(
		`
		UPDATE owl_query
		SET qr_time_access = :access_time
		WHERE qr_uuid = :uuid
		`,
		map[string]interface{} {
			"access_time": accessTime,
			"uuid": query.Uuid,
		},
	))

	if result.RowsAffected() == 0 {
		AddOrRefreshQuery(query, accessTime)
	}
}

func LoadQueryByUuidAndUpdateAccessTime(name string, uuid uuid.UUID, accessTime time.Time) *model.Query {
	result := db.ToResultExt(DbFacade.SqlxDbCtrl.NamedExec(
		`
		UPDATE owl_query
		SET qr_time_access = :access_time
		WHERE qr_uuid = :uuid
		`,
		map[string]interface{} {
			"access_time": accessTime,
			"uuid": db.DbUuid(uuid),
		},
	))

	if result.RowsAffected() == 0 {
		return nil
	}

	var queryData model.Query
	err := DbFacade.SqlxDb.QueryRowx(
		`
		SELECT qr_uuid, qr_named_id, qr_content, qr_md5_content
		FROM owl_query
		WHERE qr_uuid = ?
		`,
		db.DbUuid(uuid),
	).StructScan(&queryData)

	if err == sql.ErrNoRows {
		return nil
	}

	db.PanicIfError(err)

	return &queryData
}

type addOrRefreshQueryTx struct {
	query *model.Query
	accessTime time.Time
}

func (queryTx *addOrRefreshQueryTx) InTx(tx *sqlx.Tx) db.TxFinale {
	txExt := sqlxExt.ToTxExt(tx)

	queryTx.performAddOrUpdate(txExt)
	queryTx.loadUuid(txExt)

	return db.TxCommit
}
func (queryTx *addOrRefreshQueryTx) performAddOrUpdate(txExt *sqlxExt.TxExt) {
	queryObject := queryTx.query
	if queryObject.Uuid.IsNil() {
		queryObject.NewUuid()
	}

	txExt.NamedExec(
		`
		INSERT INTO owl_query(
			qr_uuid, qr_named_id, qr_content, qr_md5_content,
			qr_time_creation, qr_time_access
		)
		VALUES(
			:uuid, :named_id, :content, :md5_content,
			:access_time, :access_time
		)
		ON DUPLICATE KEY UPDATE
			qr_time_access = :access_time
		`,
		map[string]interface{} {
			"uuid": queryObject.Uuid,
			"named_id": queryObject.NamedId,
			"content": queryObject.Content,
			"md5_content": queryObject.Md5Content,
			"access_time": queryTx.accessTime,
		},
	)
}
func (queryTx *addOrRefreshQueryTx) loadUuid(txExt *sqlxExt.TxExt) {
	queryObject := queryTx.query

	txExt.QueryRowxAndScan(
		`
		SELECT qr_uuid
		FROM owl_query
		WHERE qr_named_id = ?
			AND qr_md5_content = ?
		`,
		[]interface{} {
			queryObject.NamedId,
			queryObject.Md5Content,
		},
		&queryObject.Uuid,
	)
}
