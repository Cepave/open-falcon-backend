package json

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

// This type of time would be serialized to UNIX time for MarshalJSON()
type JsonTime time.Time

// MarshalJSON does the serialization of UNIX timestamp
func (t JsonTime) MarshalJSON() ([]byte, error) {
	if time.Time(t).IsZero() {
		return ([]byte)("null"), nil
	}
	return ([]byte)(t.String()), nil
}

// UnmarshalJSON does the deserialization of UNIX timestamp
func (t *JsonTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" {
		*t = JsonTime(time.Time{})
		return nil
	}
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return fmt.Errorf("Cannot parse timestamp string: %s\n", data)
	}
	*t = JsonTime(time.Unix(i, 0))
	return nil
}

func (t JsonTime) Value() (driver.Value, error) {
	return time.Time(t).Unix(), nil
}

func (t JsonTime) String() string {
	return fmt.Sprintf("%d", time.Time(t).Unix())
}
