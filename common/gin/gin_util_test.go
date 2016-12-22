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
		orderBy string
		expectedSize int32
		expectedPos int32
		expectedOrderByEntity int
	} {
		{ "", "", "", 50, 1, 0 },
		{ "20", "4", "p1#asc:p2#desc", 20, 4, 2 },
	}

	defaultPaging := model.NewUndefinedPaging()
	defaultPaging.Size = 50
	defaultPaging.Position = 1

	for _, testCase := range testCases {
		req, _ := http.NewRequest("GET", "http://127.0.0.1/fake", nil)

		req.Header.Add("page-size", testCase.pageSize)
		req.Header.Add("page-pos", testCase.pagePos)
		req.Header.Add("order-by", testCase.orderBy)

		context := &gin.Context{
			Request: req,
		}

		testedPaging := PagingByHeader(context, defaultPaging)
		c.Logf("Paging: %s", testedPaging)

		c.Assert(testedPaging.Size, Equals, testCase.expectedSize)
		c.Assert(testedPaging.Position, Equals, testCase.expectedPos)
		c.Assert(testedPaging.OrderBy, HasLen, testCase.expectedOrderByEntity)
	}
}

func (suite *TestGinUtilSuite) TestParseOrderBy(c *C) {
	testCases := []struct {
		sampleValue string
		expectedOrderBy []*model.OrderByEntity
		hasError bool
	} {
		{ // normal case
			"p1#asc:p1_1#a:p1_2#ascending:p2#desc:p2_1#d:p2_2#descending:p3",
			[]*model.OrderByEntity {
				&model.OrderByEntity{ "p1", model.Ascending },
				&model.OrderByEntity{ "p1_1", model.Ascending },
				&model.OrderByEntity{ "p1_2", model.Ascending },
				&model.OrderByEntity{ "p2", model.Descending },
				&model.OrderByEntity{ "p2_1", model.Descending },
				&model.OrderByEntity{ "p2_2", model.Descending },
				&model.OrderByEntity{ "p3", model.DefaultDirection },
			},
			false,
		},
		{ // UPPER CASE for direction
			"t1#ASC:t2#DESC:t3",
			[]*model.OrderByEntity {
				&model.OrderByEntity{ "t1", model.Ascending },
				&model.OrderByEntity{ "t2", model.Descending },
				&model.OrderByEntity{ "t3", model.DefaultDirection },
			},
			false,
		},
		{ // No direction
			"c_1:c_2",
			[]*model.OrderByEntity {
				&model.OrderByEntity{ "c_1", model.DefaultDirection },
				&model.OrderByEntity{ "c_2", model.DefaultDirection },
			},
			false,
		},
		{
			"#asc",
			[]*model.OrderByEntity {},
			true,
		},
		{
			"abc#asc:",
			[]*model.OrderByEntity {},
			true,
		},
	}

	for _, testCase := range testCases {
		testedResult, err := parseOrderBy(testCase.sampleValue)

		if !testCase.hasError {
			c.Assert(err, IsNil)
		} else {
			c.Assert(err, NotNil)
		}
		c.Assert(testedResult, HasLen, len(testCase.expectedOrderBy))

		for i, orderByEntity := range testedResult {
			c.Assert(testCase.expectedOrderBy[i], DeepEquals, orderByEntity)
		}
	}
}
