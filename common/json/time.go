package json

import (
	"fmt"
	"time"
)

// This type of time would be serialized to UNIX time for MarshalJSON()
type JsonTime time.Time

func (t JsonTime) MarshalJSON() ([]byte, error) {
	jsonResult := "null"

	timeValue := time.Time(t)
	if !timeValue.IsZero() {
		jsonResult = fmt.Sprintf("%d", timeValue.Unix())
	}

	return ([]byte)(jsonResult), nil
}
