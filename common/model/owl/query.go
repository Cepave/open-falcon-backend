package owl

import (
	"fmt"
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

func (q *Query) String() string {
	return fmt.Sprintf(
		"Named Id: [%s]. Uuid: [%s]. Md5 Content: [%x]",
		q.NamedId, uuid.UUID(q.Uuid).String(),
		q.Md5Content,
	)
}

func (q *Query) ToJson() map[string]interface{} {
	return map[string]interface{} {
		"query_id": uuid.UUID(q.Uuid).String(),
	}
}
