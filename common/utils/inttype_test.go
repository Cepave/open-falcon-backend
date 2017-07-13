package utils

import (
	. "gopkg.in/check.v1"
)

type TestIntTypeSuite struct{}

var _ = Suite(&TestIntTypeSuite{})

// Tests the conversion from array of uint64 to uint32
func (suite *TestIntTypeSuite) TestUintTo32(c *C) {
	testCases := []*struct {
		source         []uint64
		expectedResult []uint32
	}{
		{
			[]uint64{0, 77, 4294967295},
			[]uint32{0, 77, 4294967295},
		},
		{nil, nil},
	}

	for _, testCase := range testCases {
		c.Assert(
			UintTo32(testCase.source), DeepEquals, testCase.expectedResult,
		)
	}
}

// Tests the conversion from array of uint64 to uint16
func (suite *TestIntTypeSuite) TestUintTo16(c *C) {
	testCases := []*struct {
		source         []uint64
		expectedResult []uint16
	}{
		{
			[]uint64{0, 308, 65535},
			[]uint16{0, 308, 65535},
		},
		{nil, nil},
	}

	for _, testCase := range testCases {
		c.Assert(
			UintTo16(testCase.source), DeepEquals, testCase.expectedResult,
		)
	}
}

// Tests the conversion from array of uint64 to uint8
func (suite *TestIntTypeSuite) TestUintTo8(c *C) {
	testCases := []*struct {
		source         []uint64
		expectedResult []uint8
	}{
		{
			[]uint64{0, 13, 255},
			[]uint8{0, 13, 255},
		},
		{nil, nil},
	}

	for _, testCase := range testCases {
		c.Assert(
			UintTo8(testCase.source), DeepEquals, testCase.expectedResult,
		)
	}
}

// Tests the conversion from array of int64 to int32
func (suite *TestIntTypeSuite) TestIntTo32(c *C) {
	testCases := []*struct {
		source         []int64
		expectedResult []int32
	}{
		{
			[]int64{-2147483648, 1377, 2147483647},
			[]int32{-2147483648, 1377, 2147483647},
		},
		{nil, nil},
	}

	for _, testCase := range testCases {
		c.Assert(
			IntTo32(testCase.source), DeepEquals, testCase.expectedResult,
		)
	}
}

// Tests the conversion from array of int64 to int16
func (suite *TestIntTypeSuite) TestIntTo16(c *C) {
	testCases := []*struct {
		source         []int64
		expectedResult []int16
	}{
		{
			[]int64{-32768, 508, 32767},
			[]int16{-32768, 508, 32767},
		},
		{nil, nil},
	}

	for _, testCase := range testCases {
		c.Assert(
			IntTo16(testCase.source), DeepEquals, testCase.expectedResult,
		)
	}
}

// Tests the conversion from array of int64 to int8
func (suite *TestIntTypeSuite) TestIntTo8(c *C) {
	testCases := []*struct {
		source         []int64
		expectedResult []int8
	}{
		{
			[]int64{-128, 55, 127},
			[]int8{-128, 55, 127},
		},
		{nil, nil},
	}

	for _, testCase := range testCases {
		c.Assert(
			IntTo8(testCase.source), DeepEquals, testCase.expectedResult,
		)
	}
}

// Tests the sorting and unique for int64
func (suite *TestIntTypeSuite) TestSortAndUniqueInt64(c *C) {
	testCases := []*struct {
		source         []int64
		expectedResult []int64
	}{
		{
			[]int64{30, 30, -10, 10, -7, 22},
			[]int64{-10, -7, 10, 22, 30},
		},
		{nil, nil},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := SortAndUniqueInt64(testCase.source)

		c.Assert(testedResult, DeepEquals, testCase.expectedResult, comment)
	}
}

// Tests the sorting and unique for uint64
func (suite *TestIntTypeSuite) TestSortAndUniqueUint64(c *C) {
	testCases := []*struct {
		source         []uint64
		expectedResult []uint64
	}{
		{
			[]uint64{30, 30, 10, 10, 7, 22},
			[]uint64{7, 10, 22, 30},
		},
		{nil, nil},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := SortAndUniqueUint64(testCase.source)

		c.Assert(testedResult, DeepEquals, testCase.expectedResult, comment)
	}
}
