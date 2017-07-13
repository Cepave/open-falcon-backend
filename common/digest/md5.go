package digest

import (
	"crypto/md5"
)

type Md5DigestValue [md5.Size]byte

// This function would call Digestor.GetDigest() and
// "glue" every value of digestor together.
func SumAllToMd5(rest ...Digestor) Md5DigestValue {
	digest := SumAll(
		Md5SumFunc, rest...,
	)

	var finalValue Md5DigestValue
	copy(finalValue[:], digest)

	return finalValue
}

type StringMd5Digestor string

func (d StringMd5Digestor) GetDigest() []byte {
	return Md5SumFunc([]byte(d))
}

type BytesMd5Digestor []byte

func (d BytesMd5Digestor) GetDigest() []byte {
	return Md5SumFunc(d)
}

func Md5SumFunc(v []byte) []byte {
	md5Digest := md5.Sum(v)
	return md5Digest[:]
}
