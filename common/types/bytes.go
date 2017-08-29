package types

import (
	"encoding/base64"

	"github.com/juju/errors"
)

type SupportBase64 interface {
	ToBase64() string
	FromBase64(string) error
	MustFromBase64(string)
}

type VarBytes []byte

func (p *VarBytes) ToBase64() string {
	return base64.StdEncoding.EncodeToString(*p)
}
func (p *VarBytes) FromBase64(v string) error {
	sourceBytes, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return errors.Annotatef(err, "Cannot parse string as base64 of bytes. Source string: [%s]", v)
	}

	*p = sourceBytes
	return nil
}
func (p *VarBytes) MustFromBase64(v string) {
	if err := p.FromBase64(v); err != nil {
		panic(errors.Details(err))
	}
}

type Bytes16 [16]byte

func (p *Bytes16) FromVarBytes(srcBytes VarBytes) error {
	if len(srcBytes) != 16 {
		return errors.Errorf("Len of source bytes is not 16. It is [%d].", len(srcBytes))
	}

	copy(p[:], srcBytes[:16])
	return nil
}
func (p *Bytes16) MustFromVarBytes(srcBytes VarBytes) {
	if err := p.FromVarBytes(srcBytes); err != nil {
		panic(errors.Details(err))
	}
}
func (p *Bytes16) ToVarBytes() VarBytes {
	varBytes := make(VarBytes, 16)

	copy(varBytes, p[:])
	return varBytes
}
func (p *Bytes16) ToBase64() string {
	return base64.StdEncoding.EncodeToString(p[:])
}
func (p *Bytes16) FromBase64(v string) error {
	sourceBytes, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return errors.Annotatef(err, "Cannot parse string as base64 of bytes. Source string: [%s]", v)
	}

	if len(sourceBytes) != 16 {
		return errors.Errorf("The length of byte array is not [16]: %d", len(sourceBytes))
	}

	copy(p[:], sourceBytes)
	return nil
}
func (p *Bytes16) MustFromBase64(v string) {
	if err := p.FromBase64(v); err != nil {
		panic(errors.Details(err))
	}
}
