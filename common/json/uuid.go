package json

import (
	"encoding/json"

	uuid "github.com/satori/go.uuid"
)

type Uuid uuid.UUID

// For uuid.Nil, the JSON value would be "null"
func (u Uuid) MarshalJSON() ([]byte, error) {
	nativeUuid := uuid.UUID(u)

	if nativeUuid == uuid.Nil {
		return nullValue, nil
	}

	return json.Marshal(nativeUuid.String())
}

// For empty string or "null" value, the value of uuid would be uuid.Nil
func (u *Uuid) UnmarshalJSON(v []byte) error {
	var stringValue string

	err := json.Unmarshal(v, &stringValue)
	if err != nil {
		return err
	}

	if stringValue == "" {
		*u = Uuid(uuid.Nil)
		return nil
	}

	parsedUuid, err := uuid.FromString(stringValue)
	if err != nil {
		return err
	}

	*u = Uuid(parsedUuid)

	return nil
}
