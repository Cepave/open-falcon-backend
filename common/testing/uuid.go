package testing

import (
	"encoding/hex"
	"github.com/satori/go.uuid"
	check "gopkg.in/check.v1"
	"strings"
)

// Parses string of UUID, supporting both of HEX with hypen("-") and HEX string
func ParseUuid(c *check.C, uuidString string) uuid.UUID {
	var uuidValue uuid.UUID
	var err error

	if strings.Contains(uuidString, "-") {
		uuidValue, err = uuid.FromString(uuidString)
		c.Assert(err, check.IsNil)
	} else {
		bytesOfUuid, err := hex.DecodeString(uuidString)
		c.Assert(err, check.IsNil)

		uuidValue, err = uuid.FromBytes(bytesOfUuid)
		c.Assert(err, check.IsNil)
	}

	return uuidValue
}
