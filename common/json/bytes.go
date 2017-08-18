package json

import (
	"encoding/base64"
	"encoding/json"

	"github.com/Cepave/open-falcon-backend/common/types"
	"github.com/juju/errors"
)

type Bytes16 types.Bytes16

// For empty bytes, the value of base64 would be "AAAAAAAAAAAAAAAAAAAAAA=="
func (b Bytes16) MarshalJSON() ([]byte, error) {
	return json.Marshal(base64.StdEncoding.EncodeToString(b[:]))
}

// For empty string or "null" of JSON, the value of byte array[16] would be [16]byte{}
func (b *Bytes16) UnmarshalJSON(v []byte) error {
	var jsonString string

	err := json.Unmarshal(v, &jsonString)
	if err != nil {
		return err
	}

	if jsonString == "" {
		for i, _ := range b {
			b[i] = 0
		}
		return nil
	}

	decoded, err := base64.StdEncoding.DecodeString(jsonString)
	if err != nil {
		return errors.Annotate(err, "Decode base64 string has error")
	}

	if len(decoded) != 16 {
		return errors.Errorf("The length of byte array is not 16: %d", len(decoded))
	}

	copy(b[:], decoded)
	return nil
}

func (b *Bytes16) IsZero() bool {
	for _, v := range b {
		if v != 0 {
			return false
		}
	}

	return true
}

type VarBytes types.VarBytes

// For nil slice and empty slice, the value of JSON would be "null"
func (b VarBytes) MarshalJSON() ([]byte, error) {
	if len(b) == 0 {
		return nullValue, nil
	}

	return json.Marshal(base64.StdEncoding.EncodeToString(b[:]))
}

// For empty string or "null" of JSON, the value of byte slice would be []byte(nil)
func (b *VarBytes) UnmarshalJSON(v []byte) error {
	var jsonString string

	err := json.Unmarshal(v, &jsonString)
	if err != nil {
		return errors.Annotate(err, "JSON unmarshaling has error")
	}

	if jsonString == "" {
		*b = nil
		return nil
	}

	decoded, err := base64.StdEncoding.DecodeString(jsonString)
	if err != nil {
		return errors.Annotate(err, "Decode base64 string has error")
	}

	*b = decoded
	return nil
}
