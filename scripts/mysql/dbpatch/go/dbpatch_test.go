package main

import (
	"testing"

	"github.com/Cepave/open-falcon-backend/scripts/mysql/dbpatch/go/changelog"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type DefaultSuite struct{}

var _ = Suite(&DefaultSuite{})

/**
 * Tests the checking of configuration
 */
func (s *DefaultSuite) TestCheckRunPatchConfig(c *C) {
	var testCases = []struct {
		driverName     string
		dsn            string
		changeLog      string
		patchFileBase  string
		expectedResult bool
	}{
		{"mysql", "msyql://localhost/", "changeLog.json", "patch-files", true}, // Passed check
		{"", "msyql://localhost/", "changeLog.json", "patch-files", false},     // Failed check
		{"mysql", "", "changeLog.json", "patch-files", false},                  // Failed check
	}

	for _, testCase := range testCases {
		var testedResult = checkRunPatchConfig(&changelog.ChangeLogConfig{
			DriverName:    testCase.driverName,
			Dsn:           testCase.dsn,
			ChangeLog:     testCase.changeLog,
			PatchFileBase: testCase.patchFileBase,
		})

		c.Assert(testedResult, Equals, testCase.expectedResult)
	}
}
