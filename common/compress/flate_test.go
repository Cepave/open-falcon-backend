package compress

import (
	"compress/flate"
	. "gopkg.in/check.v1"
)

type TestFlateSuite struct{}

var _ = Suite(&TestFlateSuite{})

// Tests the compression(and decompress of data
func (suite *TestFlateSuite) TestMustCompressString(c *C) {
	testCases := []*struct {
		sampleString string
	}{
		{"{}"},
		{
			`面對二○一七年全球經濟變局，林全研判，川普不可能揚棄自由貿易的方向，不論TPP替代方案為何，台灣必須準備好。他也第一次披露，行政院即將新訂針對外籍專業人士來台的特別法；及未來三年，台灣將推千億鐵道交通建設；個人與企業所得稅要併案稅改等新建設。`,
		},
		{
			`{"filters":{"time":{"start_time":2908001,"end_time":2909001,"to_now":{"unit":"","value":0}},"agent":{"name":["CB1","KC2"],"hostname":["GA3","ZC0"],"ip_address":["10.9","11.56.71.89"],"connection_i
d":["AB@13","AC@13"],"isp_ids":[11,12],"province_ids":[5,8,9],"city_ids":[31,34],"name_tag_ids":[10,19],"group_tag_ids":[45,51]},"target":{"name":["CB1","KC2"],"host":["GA3","ZC0"],"isp_ids":[13,17],"province_id
s":[24,39,81],"city_ids":[14,23],"name_tag_ids":[39,46],"group_tag_ids":[61,63]},"metrics":"$max \u003e 100 or $min \u003c 30"},"grouping":{"agent":["name","province"],"target":["isp"]},"output":{"metrics":["min
","loss"]}}`,
		},
	}

	testedCompressor := &FlateCompressor{flate.DefaultCompression}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResultBytes := testedCompressor.MustCompressString(testCase.sampleString)
		c.Logf("[Case %d] Source size:[%d] Compressed Size: [%d]", i+1, len(testCase.sampleString), len(testedResultBytes))

		testedResultString := testedCompressor.MustDecompressToString(testedResultBytes)

		c.Assert(testedResultString, DeepEquals, testCase.sampleString, comment)
	}
}
