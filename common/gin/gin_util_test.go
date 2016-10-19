package gin

import (
	model "github.com/Cepave/open-falcon-backend/common/model"
	"net/http"
	"gopkg.in/gin-gonic/gin.v1"
	. "gopkg.in/check.v1"
)

type TestGinUtilSuite struct{}

var _ = Suite(&TestGinUtilSuite{})

// Tests the paging parameters by header
func (suite *TestGinUtilSuite) TestPagingByHeader(c *C) {
	testCases := []struct {
		pageSize string
		pagePos string
		expectedSize int32
		expectedPos int32
	} {
		{ "", "", 50, 1 },
		{ "20", "4", 20, 4 },
	}

	defaultPaging := model.NewUndefinedPaging()
	defaultPaging.Size = 50
	defaultPaging.Position = 1

	for _, testCase := range testCases {
		req, _ := http.NewRequest("GET", "http://127.0.0.1/fake", nil)

		req.Header.Add("page-size", testCase.pageSize)
		req.Header.Add("page-pos", testCase.pagePos)

		context := &gin.Context{
			Request: req,
		}

		testedPaging := PagingByHeader(context, defaultPaging)
		c.Logf("Paging: %s", testedPaging)

		c.Assert(testedPaging.Size, Equals, testCase.expectedSize)
		c.Assert(testedPaging.Position, Equals, testCase.expectedPos)
	}
}
