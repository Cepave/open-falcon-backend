package nqm

import (
	"fmt"
	model "github.com/Cepave/open-falcon-backend/modules/query/model/nqm"
	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	qtest "github.com/Cepave/open-falcon-backend/modules/query/test"
	"github.com/satori/go.uuid"
	. "gopkg.in/check.v1"
)

type TestCompoundReportSuite struct{}

var _ = Suite(&TestCompoundReportSuite{})

// Tests the convertion of query to detail information
func (suite *TestCompoundReportSuite) TestToQueryDetail(c *C) {
	/**
	 * Sets-up sample query
	 */
	sampleQuery := model.NewCompoundQuery()

	sampleQuery.Filters.Metrics = "$max >= $min"
	sampleQuery.Grouping.Agent = []string{ model.AgentGroupingHostname, model.GroupingIsp }
	sampleQuery.Grouping.Target = []string{ model.TargetGroupingHost, model.GroupingProvince, model.GroupingIsp }
	sampleQuery.Output.Metrics = []string{ model.MetricAvg, model.MetricMax, model.MetricLoss }

	agentFilter := sampleQuery.Filters.Agent

	agentFilter.Name = []string { "Cool1", "Cool2" }
	agentFilter.Hostname = []string { "gc1.com", "gc2.com" }
	agentFilter.IpAddress = []string { "123.71.1", "123.71.2" }
	agentFilter.ConnectionId = []string { "c01", "c02" }
	agentFilter.IspIds = []int16 { 2, 3, 4, 1, 5, 6 }
	agentFilter.ProvinceIds = []int16 { 3, 4, 7, 8 }
	agentFilter.CityIds = []int16 { 13, 14 }
	agentFilter.NameTagIds = []int16 { 3375, 3376 }
	agentFilter.GroupTagIds = []int32 { 90801, 90802, 90803 }

	targetFilter := sampleQuery.Filters.Target

	targetFilter.Name = []string { "Zoo-1", "Zoo-2" }
	targetFilter.Host = []string { "kz1.com", "kz2.com" }
	targetFilter.IspIds = []int16 { 5, 7 }
	targetFilter.ProvinceIds = []int16 { 1, 3, 8 }
	targetFilter.CityIds = []int16 { 15, 16 }
	targetFilter.NameTagIds = []int16 { 3375, 3376 }
	targetFilter.GroupTagIds = []int32 { 90801, 90802, 90803 }
	// :~)

	testedDetail := ToQueryDetail(sampleQuery)

	c.Assert(string(testedDetail.Metrics), Equals, sampleQuery.Filters.Metrics)

	/**
	 * Asserts the query detail on agent conditions
	 */
	testedAgentDetail := testedDetail.Agent
	c.Assert(testedAgentDetail.Name, DeepEquals, agentFilter.Name)
	c.Assert(testedAgentDetail.Hostname, DeepEquals, agentFilter.Hostname)
	c.Assert(testedAgentDetail.IpAddress, DeepEquals, agentFilter.IpAddress)
	c.Assert(testedAgentDetail.ConnectionId, DeepEquals, agentFilter.ConnectionId)
	c.Assert(testedAgentDetail.Isps, HasLen, 6)
	c.Assert(testedAgentDetail.Provinces, HasLen, 4)
	c.Assert(testedAgentDetail.Cities, HasLen, 2)
	c.Assert(testedAgentDetail.NameTags, HasLen, 2)
	c.Assert(testedAgentDetail.GroupTags, HasLen, 3)
	// :~)

	/**
	 * Asserts the query detail on target conditions
	 */
	testedTargetDetail := testedDetail.Target
	c.Assert(testedTargetDetail.Name, DeepEquals, targetFilter.Name)
	c.Assert(testedTargetDetail.Host, DeepEquals, targetFilter.Host)
	c.Assert(testedTargetDetail.Isps, HasLen, 2)
	c.Assert(testedTargetDetail.Provinces, HasLen, 3)
	c.Assert(testedTargetDetail.Cities, HasLen, 2)
	c.Assert(testedTargetDetail.NameTags, HasLen, 2)
	c.Assert(testedTargetDetail.GroupTags, HasLen, 3)
	// :~)

	/**
	 * Asserts the output detail
	 */
	testedOutputDetail := testedDetail.Output
	c.Assert(testedOutputDetail.Agent, HasLen, 2)
	c.Assert(testedOutputDetail.Target, HasLen, 3)
	c.Assert(testedOutputDetail.Metrics, HasLen, 3)
	// :~)
}

// Tests the building of query object
func (suite *TestCompoundReportSuite) TestBuildQuery(c *C) {
	sampleQuery := model.NewCompoundQuery()

	sampleJson := []byte(`
	{
		"filters": {
			"time": {
				"start_time": 1336608000,
				"end_time": 1336622400
			},
			"agent": {
				"name": [ "GD-1", "GD-2" ]
			},
			"target": {
				"host": [ "18.98.7.61", "google.com" ]
			}
		}
	}`)
	c.Assert(sampleQuery.UnmarshalJSON(sampleJson), IsNil)

	testedResult1 := BuildQuery(sampleQuery)

	c.Logf("[T-1] Query object: %s", testedResult1)
	c.Assert(testedResult1.NamedId, Equals, queryNamedId)

	/**
	 * Asserts sample query with same conditions
	 */
	testedResult2 := BuildQuery(sampleQuery)
	c.Logf("[T-2] Query object: %s", testedResult2)
	c.Assert(testedResult1, DeepEquals, testedResult2)
	// :~)
}

// Tests the loading of compound query by UUID
func (suite *TestCompoundReportSuite) TestGetCompoundQueryByUuid(c *C) {
	sampleQuery := model.NewCompoundQuery()
	err := sampleQuery.UnmarshalJSON([]byte(`
	{
		"filters": {
			"time": {
				"start_time": 190807000,
				"end_time": 190827000
			}
		},
		"output": {
			"metrics": [ "max", "min", "avg", "loss" ]
		}
	}
	`))
	c.Assert(err, IsNil)

	/**
	 * Builds query object(persist the query) and load it with generated UUID
	 */
	queryObject := BuildQuery(sampleQuery)
	testedQuery := GetCompoundQueryByUuid(uuid.UUID(queryObject.Uuid))
	// :~)

	c.Assert(testedQuery.Filters.Time, DeepEquals, sampleQuery.Filters.Time)
	c.Assert(testedQuery.Output, DeepEquals, sampleQuery.Output)
}

func (s *TestCompoundReportSuite) SetUpTest(c *C) {
	inTx := owlDb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestCompoundReportSuite.TestToQueryDetail":
		inTx(
			`
			INSERT INTO owl_group_tag(gt_id, gt_name)
			VALUES(90801, 'gt-ab-1'), (90802, 'gt-ab-2'), (90803, 'gt-ab-3')
			`,
			`
			INSERT INTO owl_name_tag(nt_id, nt_value)
			VALUES(3375, 'nt-ab-1'), (3376, 'nt-ab-2')
			`,
		)
	}
}
func (s *TestCompoundReportSuite) TearDownTest(c *C) {
	inTx := owlDb.DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case "TestCompoundReportSuite.TestToQueryDetail":
		inTx(
			`DELETE FROM owl_group_tag WHERE gt_id >= 90801 AND gt_id <= 90803`,
			`DELETE FROM owl_name_tag WHERE nt_id >= 3375 AND nt_id <= 3376`,
		)
	case "TestCompoundReportSuite.TestGetCompoundQueryByUuid":
		inTx(
			fmt.Sprintf(`DELETE FROM owl_query WHERE qr_named_id = '%s'`, queryNamedId),
		)
	case "TestCompoundReportSuite.TestBuildQuery":
		inTx(
			fmt.Sprintf(`DELETE FROM owl_query WHERE qr_named_id = '%s'`, queryNamedId),
		)
	}
}

func (s *TestCompoundReportSuite) SetUpSuite(c *C) {
	qtest.InitDb(c)
	initServices()
}
func (s *TestCompoundReportSuite) TearDownSuite(c *C) {
	qtest.ReleaseDb(c)
}
