package db

import (
	"database/sql/driver"

	"github.com/satori/go.uuid"
)

// Uses the []byte for database/sql/driver.Value
type DbUuid uuid.UUID

// This function supports []byte or string value
func (u *DbUuid) Scan(src interface{}) error {
	srcUuid := (*uuid.UUID)(u)
	return srcUuid.Scan(src)
}

// This function would supply []byte as value of database driver
func (u DbUuid) Value() (driver.Value, error) {
	srcUuid := uuid.UUID(u)
	if srcUuid == uuid.Nil {
		return driver.Null{}, nil
	}

	return srcUuid.Bytes(), nil
}

func (u DbUuid) IsNil() bool {
	return uuid.UUID(u) == uuid.Nil
}

func (u DbUuid) ToUuid() uuid.UUID {
	return uuid.UUID(u)
}
