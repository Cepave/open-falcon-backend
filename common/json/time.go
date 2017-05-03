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
	jsonResult := "null"

	timeValue := time.Time(t)
	if !timeValue.IsZero() {
		jsonResult = fmt.Sprintf("%d", timeValue.Unix())
	}

	return ([]byte)(jsonResult), nil
}

// UnmarshalJSON does the deserialization of UNIX timestamp
func (t *JsonTime) UnmarshalJSON(data []byte) error {
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return fmt.Errorf("Cannot parse timestamp string: %s\n", data)
	}
	*t = JsonTime(time.Unix(i, 0))
	return nil
}

func (t JsonTime) Value() (driver.Value, error) {
	return time.Time(t).Unix(), nil
}
