package owl

import (
	"github.com/Cepave/open-falcon-backend/common/db"
	"github.com/satori/go.uuid"
)

// Represents the query object
type Query struct {
	Uuid db.DbUuid

	NamedId string
	Content []byte
	Md5Content [16]byte
}

// Generate UUID v4 for this object
func (q *Query) NewUuid() {
	q.Uuid = db.DbUuid(uuid.NewV4())
}
