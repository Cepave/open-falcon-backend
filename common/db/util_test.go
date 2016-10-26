package db

import (
	"database/sql"
	. "gopkg.in/check.v1"
)

type DbUtilTestSuite struct{}

var _ = Suite(&DbUtilTestSuite{})

// Tests the convertion for like clause(SQL) for binary string
func (suite *DbUtilTestSuite) TestIpV4ToBytesForLike(c *C) {
	testCases := []struct {
		sampleIp string
		expectedResult []byte
	} {
		{
			"10.20.30.40",
			[]byte {
				0x0A, 0x14, 0x1E, 0x28, 0x25,
			},
		},
		{
			"10.20",
			[]byte {
				0x0A, 0x14, 0x25,
			},
		},
		{
			"10.37.12",
			[]byte {
				0x0A, 0x5C, 0x25, 0x0C, 0x25,
			},
		},
	}

	for _, testCase := range testCases {
		testedIp, err := IpV4ToBytesForLike(testCase.sampleIp)

		c.Assert(err, IsNil)
		c.Assert(testedIp, DeepEquals, testCase.expectedResult)
	}
}

// Tests the convertion for array of string values
func (suite *DbUtilTestSuite) TestGroupedStringToStringArray(c *C) {
	testCases := []struct {
		sampleSqlString sql.NullString
		expectedResult []string
	} {
		{ // Normal data
			sql.NullString{
				"v1,v2,v3,v4", true,
			},
			[]string { "v1", "v2", "v3", "v4" },
		},
		{ // Null data
			sql.NullString{
				"<NULL_STRING>", false,
			},
			nil,
		},
	}

	for _, testCase := range testCases {
		testedResult := GroupedStringToStringArray(testCase.sampleSqlString, ",")
		c.Assert(testedResult, DeepEquals, testCase.expectedResult)
	}
}

// Tests the convertion for array of int values
func (suite *DbUtilTestSuite) TestGroupedStringToIntArray(c *C) {
	testCases := []struct {
		sampleSqlString sql.NullString
		expectedResult []int64
	} {
		{ // Normal data
			sql.NullString{
				"34|-67|81|-91", true,
			},
			[]int64 { 34, -67, 81, -91 },
		},
		{ // Null data
			sql.NullString{
				"<NULL_STRING>", false,
			},
			nil,
		},
	}

	for _, testCase := range testCases {
		testedResult := GroupedStringToIntArray(testCase.sampleSqlString, "|")
		c.Assert(testedResult, DeepEquals, testCase.expectedResult)
	}
}

// Tests the convertion for array of int(unsigned) values
func (suite *DbUtilTestSuite) TestGroupedStringToUintArray(c *C) {
	testCases := []struct {
		sampleSqlString sql.NullString
		expectedResult []uint64
	} {
		{ // Normal data
			sql.NullString{
				"22@167@33@42", true,
			},
			[]uint64 { 22, 167, 33, 42 },
		},
		{ // Null data
			sql.NullString{
				"<NULL_STRING>", false,
			},
			nil,
		},
	}

	for _, testCase := range testCases {
		testedResult := GroupedStringToUintArray(testCase.sampleSqlString, "@")
		c.Assert(testedResult, DeepEquals, testCase.expectedResult)
	}
}
