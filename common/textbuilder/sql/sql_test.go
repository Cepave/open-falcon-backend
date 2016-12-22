package sql

import (
	//ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	. "gopkg.in/check.v1"
)

type TestSqlSuite struct{}

var _ = Suite(&TestSqlSuite{})

// Tests the sql builder
//func (suite *TestSqlSuite) TestComplexSql(c *C) {
	//testedGetter := SqlWhere(
		//SqlAnd(
			//TextGetterByString("tb_name = 'OK'"),
			//TextGetterByString("tb_id = 33"),
			//SqlParenthesis(
				//SqlOr(
					//TextGetterByString("tb_v1 = 'OK'"),
					//TextGetterByString("tb_v2 = 33"),
				//),
			//),
			//&Surrounding {
				//TextGetterByString("tc_g1 IN ("),
				//TextGetterByString(")"),
				//SqlPlaceholders(4),
			//},
		//),
	//)

	//testedResult := testedGetter.GetText()

	//c.Assert(testedResult, ocheck.StringContains, "tb_name")
	//c.Assert(testedResult, ocheck.StringContains, "tb_v1")
	//c.Assert(testedResult, ocheck.StringContains, "tc_g1 IN")
	//c.Assert(testedResult, ocheck.StringContains, "AND")
	//c.Assert(testedResult, ocheck.StringContains, "OR")
	//c.Assert(testedResult, ocheck.StringContains, "WHERE")
	//c.Assert(testedResult, ocheck.StringContains, "?")
//}
