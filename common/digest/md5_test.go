package digest

import (
	"encoding/hex"
	. "gopkg.in/check.v1"
)

type TestMd5Suite struct{}

var _ = Suite(&TestMd5Suite{})

// Tests the digesting of MD5 from multiple Digestor
func (suite *TestMd5Suite) TestSumAllToMd5(c *C) {
	testCases := []*struct {
		sourceData  []StringMd5Digestor
		expectedMd5 string
	}{
		{
			[]StringMd5Digestor{"", ""},
			"5873dd45edd01f09c1ef2e7819369e8e",
		},
		{
			[]StringMd5Digestor{"", "", ""},
			"2bc3f09a6fbadbe687825030c6e7b9a4",
		},
		{
			[]StringMd5Digestor{"this-1", "that-2"},
			"c9a528931f6a9001d2fcbb11d0e600d3",
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		expectedResult, err := hex.DecodeString(testCase.expectedMd5)
		c.Assert(err, IsNil)

		var digestors = make([]Digestor, len(testCase.sourceData))
		for i, d := range testCase.sourceData {
			digestors[i] = d
		}

		testedResult := SumAllToMd5(digestors...)
		c.Logf("MD5: [%s]", hex.EncodeToString(testedResult[:]))

		c.Assert(testedResult[:], DeepEquals, expectedResult, comment)
	}
}
