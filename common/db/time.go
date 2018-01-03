package db

import (
	"database/sql/driver"
	"time"
)

// As of database driver, this type supported nullable-value on time.Time
//
// If the time.Time.IsZero() gives true value, this driver would gives null value.
//
// You could use "time.Time{}" to initialize the value, shich as null value in database.
type DbTime time.Time

// Supports hex string or bytes
func (t *DbTime) Scan(src interface{}) error {
	if src == nil {
		*t = DbTime(time.Time{})
		return nil
	}

	*t = DbTime(src.(time.Time))

	return nil
}
func (t DbTime) Value() (driver.Value, error) {
	if t.IsNil() {
		return driver.Null{}, nil
	}

	return t.ToTime(), nil
}
func (t DbTime) ToTime() time.Time {
	return time.Time(t)
}
func (t DbTime) IsNil() bool {
	return t.ToTime().IsZero()
}
