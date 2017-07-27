package owl

import (
	"github.com/Cepave/open-falcon-backend/common/db"
	sqlxExt "github.com/Cepave/open-falcon-backend/common/db/sqlx"
	model "github.com/Cepave/open-falcon-backend/common/model/owl"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
	"time"
)

// Adds a query or refresh an existing one
func AddOrRefreshQuery(query *model.Query, accessTime time.Time) {
	DbFacade.SqlxDbCtrl.InTx(&addOrRefreshQueryTx{query, accessTime})
}

// Updates access time of a query or adding new one
func UpdateAccessTimeOrAddNewOne(query *model.Query, accessTime time.Time) {
	if isQueryExistingByUuid(query.Uuid) {
		updateAccessTimeByUuid(
			query.Uuid, accessTime,
		)
		return
	}

	AddOrRefreshQuery(query, accessTime)
}

// It is possible that the access times of two requests are very closed(less than 1 second)
func LoadQueryByUuidAndUpdateAccessTime(name string, uuid uuid.UUID, accessTime time.Time) *model.Query {
	updateAccessTimeByUuid(
		db.DbUuid(uuid), accessTime,
	)

	var queryData model.Query

	if !DbFacade.SqlxDbCtrl.GetOrNoRow(
		&queryData,
		`
		SELECT qr_uuid, qr_named_id, qr_content, qr_md5_content
		FROM owl_query
		WHERE qr_uuid = ?
		`,
		db.DbUuid(uuid),
	) {
		return nil
	}

	return &queryData
}

type addOrRefreshQueryTx struct {
	query      *model.Query
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
		map[string]interface{}{
			"uuid":        queryObject.Uuid,
			"named_id":    queryObject.NamedId,
			"content":     queryObject.Content,
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
		[]interface{}{
			queryObject.NamedId,
			queryObject.Md5Content,
		},
		&queryObject.Uuid,
	)
}

func updateAccessTimeByUuid(uuid db.DbUuid, accessTime time.Time) {
	DbFacade.SqlxDbCtrl.NamedExec(
		`
		UPDATE owl_query
		SET qr_time_access = :access_time
		WHERE qr_uuid = :uuid
		`,
		map[string]interface{}{
			"access_time": accessTime,
			"uuid":        uuid,
		},
	)
}
func isQueryExistingByUuid(uuid db.DbUuid) bool {
	return DbFacade.SqlxDbCtrl.GetOrNoRow(
		&struct {
			Uuid db.DbUuid `db:"qr_uuid"`
		}{},
		`
		SELECT qr_uuid
		FROM owl_query
		WHERE qr_uuid = ?
		`,
		uuid,
	)
}
