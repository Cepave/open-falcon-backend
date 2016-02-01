package main

import (
	"testing"
	"github.com/Cepave/scripts/dbpatch/go/changelog"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type DefaultSuite struct{}

var _ = Suite(&DefaultSuite{})

/**
 * Tests the checking of configuration
 */
func (s *DefaultSuite) TestCheckRunPatchConfig(c *C) {
	var testCases = [][]interface{} {
		{ "mysql", "msyql://localhost/", "changeLog.json", "patch-files", true }, // Passed check
		{ "", "msyql://localhost/", "changeLog.json", "patch-files", false }, // Failed check
		{ "mysql", "", "changeLog.json", "patch-files", false }, // Failed check
	}

	for _, testCase := range testCases {
		var testedResult = checkRunPatchConfig(&changelog.ChangeLogConfig{
			DriverName: testCase[0].(string),
			Dsn: testCase[1].(string),
			ChangeLog: testCase[2].(string),
			PatchFileBase: testCase[3].(string),
		})

		c.Assert(testedResult, Equals, testCase[4].(bool))
	}
}
