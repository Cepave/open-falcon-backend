package json

import (
	"encoding/json"
)

// When marshalling, the empty string would be "null" in JSON format
type JsonString string
func (s JsonString) MarshalJSON() ([]byte, error) {
	if s == "" {
		return []byte("null"), nil
	}

	return json.Marshal(string(s))
}

// When marshalling, the string would be the raw form of JSON format
type RawJsonForm string
func (s RawJsonForm) MarshalJSON() ([]byte, error) {
	return []byte(s), nil
}
