package db

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"

	"github.com/Cepave/open-falcon-backend/common/types"
)

type Bytes16 types.Bytes16

// Supports hex string or bytes
func (b *Bytes16) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	var srcBytes []byte

	switch srcValue := src.(type) {
	case []byte:
		srcBytes = srcValue
	case string:
		decodedHex, err := hex.DecodeString(srcValue)
		if err != nil {
			return err
		}

		srcBytes = decodedHex
	default:
		return fmt.Errorf("Needs type of \"[]byte\" or hex string. Got: %T", src)
	}

	if len(srcBytes) != 16 {
		return fmt.Errorf("Needs len of \"[]byte\" to be 16. Got: %d", len(srcBytes))
	}

	copy(b[:], srcBytes)

	return nil
}
func (b Bytes16) Value() (driver.Value, error) {
	if b.IsNil() {
		return driver.Null{}, nil
	}

	return b[:], nil
}
func (b Bytes16) IsNil() bool {
	for _, v := range b {
		if v != byte(0) {
			return false
		}
	}

	return true
}
