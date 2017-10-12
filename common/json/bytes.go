package json

import (
	"encoding/json"

	"github.com/Cepave/open-falcon-backend/common/types"
	"github.com/juju/errors"
)

type Bytes16 types.Bytes16

// For empty bytes, the value of base64 would be "AAAAAAAAAAAAAAAAAAAAAA=="
func (b Bytes16) MarshalJSON() ([]byte, error) {
	bytes16 := types.Bytes16(b)

	return json.Marshal(bytes16.ToBase64())
}

// For empty string or "null" of JSON, the value of byte array[16] would be [16]byte{}
func (b *Bytes16) UnmarshalJSON(v []byte) error {
	var jsonString string

	err := json.Unmarshal(v, &jsonString)
	if err != nil {
		return err
	}

	if jsonString == "" {
		for i := range b {
			b[i] = 0
		}
		return nil
	}

	if err := (*types.Bytes16)(b).FromBase64(jsonString); err != nil {
		return errors.Annotate(err, "[JSON] Decode base64 string has error")
	}

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
func (b *VarBytes) MarshalJSON() ([]byte, error) {
	if len(*b) == 0 {
		return nullValue, nil
	}

	return json.Marshal((*types.VarBytes)(b).ToBase64())
}

// For empty string or "null" of JSON, the value of byte slice would be []byte(nil)
func (b *VarBytes) UnmarshalJSON(v []byte) error {
	var jsonString string

	if err := json.Unmarshal(v, &jsonString); err != nil {
		return errors.Annotate(err, "JSON unmarshaling has error")
	}

	if jsonString == "" {
		*b = nil
		return nil
	}

	if err := (*types.VarBytes)(b).FromBase64(jsonString); err != nil {
		return errors.Annotate(err, "[JSON] Decode base64 string has error")
	}

	return nil
}
