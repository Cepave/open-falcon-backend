package utils

import (
	. "gopkg.in/check.v1"
)

type TestMapSuite struct{}

var _ = Suite(&TestMapSuite{})

// Tests the conversion of types for a map
func (suite *TestMapSuite) Test(c *C) {
	type s2 string

	sampleAMap := MakeAbstractMap(map[int16]s2{
		1: "Nice",
		2: "Good",
	})

	testedMap := sampleAMap.ToTypeOfTarget(int32(0), "").(map[int32]string)

	c.Logf("Map: %#v", testedMap)
	c.Assert(testedMap, HasLen, 2)
}
