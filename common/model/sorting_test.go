package model

import (
	. "gopkg.in/check.v1"
)

type TestSortingSuite struct{}

var _ = Suite(&TestSortingSuite{})

// Tests the omitting of syntax for SQL
func (suite *TestSortingSuite) TestNewSqlOrderByDialect(c *C) {
	testCases := []*struct {
		sampleEntities []*OrderByEntity
		expectedSyntax string
		hasError       bool
	}{
		{ // Only one soring entity
			[]*OrderByEntity{
				{"name", Descending},
			},
			"tb_name DESC",
			false,
		},
		{ // Multiple sorting entities
			[]*OrderByEntity{
				{"name", DefaultDirection},
				{"age", Ascending},
				{"address", Descending},
			},
			"tb_name, tb_age ASC, tb_address DESC",
			false,
		},
		{ // Empty
			[]*OrderByEntity{}, "", false,
		},
		{ // Error case(no mapping of property)
			[]*OrderByEntity{
				{"name2", DefaultDirection},
			},
			"", true,
		},
		{ // Error case(no mapping of direction)
			[]*OrderByEntity{
				{"name", 99},
			},
			"", true,
		},
	}

	testedDialect := NewSqlOrderByDialect(
		map[string]string{
			"name":    "tb_name",
			"age":     "tb_age",
			"address": "tb_address",
		},
	)

	for _, testCase := range testCases {
		testedResult, err := testedDialect.ToQuerySyntax(testCase.sampleEntities)
		c.Logf("Omit query syntax: \"%s\"", testedResult)
		if testCase.hasError {
			c.Assert(err, NotNil)
			continue
		}

		c.Assert(err, IsNil)
		c.Assert(testedResult, Equals, testCase.expectedSyntax)
	}
}
