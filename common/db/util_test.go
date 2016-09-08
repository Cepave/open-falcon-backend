package db

import (
	"testing"
	"database/sql"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type DbUtilTestSuite struct{}

var _ = Suite(&DbUtilTestSuite{})

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
