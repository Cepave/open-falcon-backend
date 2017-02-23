package reflect

import (
	"bytes"
	"encoding/binary"
	"hash/fnv"
	"reflect"
)

var byteOrder = binary.BigEndian

// Used to generate digest information of a type.
//
// The hash value of reflect.Type could be used to generate key of map on reflect.Type,
// which has better performance on look-up operation on map.
//
// algorithms for generating byte array of source:
//
// 1. UTF-8 Type.PkgPath()
// 2. UTF-8 Type.String()
//
// The algorithms of hashing: hash/fnv.New64()
func DigestType(t reflect.Type) uint64 {
	var fnvHash = fnv.New64()

	buffer := bytes.NewBuffer(make([]byte, 0, 4))

	/**
	 * Writes common properteis of any type
	 */
	if _, err := buffer.WriteString(t.PkgPath()); err != nil {
		panic(err)
	}
	if _, err := buffer.WriteString(t.String()); err != nil {
		panic(err)
	}
	// :~)

	if _, err := fnvHash.Write(buffer.Bytes()); err != nil {
		panic(err)
	}

	return fnvHash.Sum64()
}
