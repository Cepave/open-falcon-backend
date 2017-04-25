package json

import (
	sjson "github.com/bitly/go-simplejson"
	//ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	. "gopkg.in/check.v1"
)

type TestJsonExtSuite struct{}

var _ = Suite(&TestJsonExtSuite{})

// Tests the getting of path
func (suite *TestJsonExtSuite) TestGetPathExt(c *C) {
	sampleJson := sjson.New()
	sampleJson.Set("v1", 20)

	testedExt := ToJsonExt(sampleJson)

	c.Assert(testedExt.GetExt("v1"), NotNil)
	c.Assert(testedExt.GetExt("v2").IsNil(), Equals, true)
	c.Assert(testedExt.GetPathExt("v1", "ck").IsNil(), Equals, true)
}
