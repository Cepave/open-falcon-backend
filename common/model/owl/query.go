package owl

import (
	"github.com/Cepave/open-falcon-backend/common/db"
	"github.com/satori/go.uuid"
)

// Represents the query object
type Query struct {
	Uuid db.DbUuid `db:"qr_uuid"`

	NamedId string `db:"qr_named_id"`
	Content []byte `db:"qr_content"`
	Md5Content db.Bytes16 `db:"qr_md5_content"`
}

// Generate UUID v4 for this object
func (q *Query) NewUuid() {
	q.Uuid = db.DbUuid(uuid.NewV4())
}
