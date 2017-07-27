package digest

import (
	"encoding/hex"
	. "gopkg.in/check.v1"
	"time"
)

type TestTimeSuite struct{}

var _ = Suite(&TestTimeSuite{})

// Tests the getting of digest for time object
func (suite *TestTimeSuite) TestGetDigest(c *C) {
	testCases := []*struct {
		sampleTime     time.Time
		expectedDigest string
	}{
		{time.Unix(323432, 0), "000000000004ef68"},
		{time.Time{}, ""},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		hexValue := hex.EncodeToString(
			DigestableTime(testCase.sampleTime).GetDigest(),
		)
		c.Logf("Time: [%v]. Digest Value: [%s]", testCase.sampleTime, hexValue)

		c.Assert(hexValue, Equals, testCase.expectedDigest, comment)
	}
}
