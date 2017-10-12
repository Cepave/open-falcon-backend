package types

import (
	"github.com/juju/errors"
	"github.com/satori/go.uuid"
)

type Uuid uuid.UUID

func (self *Uuid) MustFromString(v string) {
	newUuid, err := uuid.FromString(v)
	if err != nil {
		err = errors.Annotatef(err, "Parse UUID string has error. String: [%s]", v)
		panic(errors.Details(err))
	}

	copy(self[:], newUuid[:])
}

func (self *Uuid) MustFromBytes(v []byte) {
	newUuid, err := uuid.FromBytes(v)
	if err != nil {
		err = errors.Annotatef(err, "UUID from bytes has error. Bytes: [%#x]", v)
		panic(errors.Details(err))
	}

	copy(self[:], newUuid[:])
}
