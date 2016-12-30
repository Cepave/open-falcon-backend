package nqm

import (
	"fmt"
	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	ojson "github.com/Cepave/open-falcon-backend/common/json"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	t "github.com/Cepave/open-falcon-backend/common/testing"
	. "gopkg.in/check.v1"
	"time"
)

type TestQueryDetailSuite struct{}

var _ = Suite(&TestQueryDetailSuite{})

// Tests the marshalling of JSON for time filter
func (suite *TestQueryDetailSuite) TestMarshalJSONOfTimeFilter(c *C) {
	sampleRelative := &TimeFilter{
		timeRangeType: TimeRangeRelative,
		ToNow: &TimeWithUnit {
			Unit: TimeUnitHour,
			Value: 6,
		},
	}
	relativeStartValue, relativeEndValue := sampleRelative.GetNetTimeRange()

	testCases := []*struct {
		timeFilter *TimeFilterDetail
		expectedJson string
	} {
		{
			&TimeFilterDetail{
				StartTime: t.ParseTimeToJsonTime(c, "2014-07-23T10:00:00Z"),
				EndTime: t.ParseTimeToJsonTime(c, "2014-07-23T15:00:00Z"),
				timeRangeType: TimeRangeAbsolute,
			},
			`{ "start_time": 1406109600, "end_time": 1406127600 }`,
		},
		{
			(*TimeFilterDetail)(sampleRelative),
			fmt.Sprintf(
				`{ "to_now": { "unit": "h", "value": 6 }, "start_time": %d, "end_time": %d }`,
				relativeStartValue.Unix(),
				relativeEndValue.Unix(),
			),
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		c.Logf("JSON: %s", ojson.MarshalJSON(testCase.timeFilter))
		c.Assert(testCase.timeFilter, ocheck.JsonEquals, testCase.expectedJson, comment)
	}
}

// Tests the marshalling of detail of query
func (suite *TestQueryDetailSuite) TestCompoundQueryDetail(c *C) {
	sampleQuery := &CompoundQueryDetail {
		Time: &TimeFilterDetail{
			timeRangeType: TimeRangeAbsolute,
			StartTime: ojson.JsonTime(time.Unix(897060500, 0)),
			EndTime: ojson.JsonTime(time.Unix(897064500, 0)),
		},
		Metrics: "$max >= 150",
		Agent: &AgentOfQueryDetail {
			Name: []string { "ag-name-1", "ag-name-2" },
			Hostname: []string{},
			IpAddress: []string{ "12.10.1", "12.10.2" },
			ConnectionId: []string{},

			Isps: []*owlModel.Isp {
				{ Id: 3, Name: "ag-isp-1" },
				{ Id: 4, Name: "ag-isp-2" },
			},
			Provinces: []*owlModel.Province{
				{ Id: 4, Name: "安東省" },
				{ Id: 5, Name: "天山省" },
			},
			Cities: []*owlModel.City2{
				{ Id: 11, Name: "吉林市" },
				{ Id: 12, Name: "香山市" },
			},
			NameTags: []*owlModel.NameTag{},
			GroupTags: []*owlModel.GroupTag{},
		},
		Target:&TargetOfQueryDetail {
			Name: []string { "tg-name-1", "tg-name-2" },
			Host: []string{ "wine-1", "wine-2" },

			Isps: []*owlModel.Isp {},
			Provinces: []*owlModel.Province{},
			Cities: []*owlModel.City2{},
			NameTags: []*owlModel.NameTag{
				{ Id: 2861, Value: "nt-3" },
				{ Id: 2862, Value: "nt-4" },
			},
			GroupTags: []*owlModel.GroupTag{
				{ Id: 30071, Name: "gt-3" },
				{ Id: 30072, Name: "gt-4" },
			},
		},
		Output: &OutputDetail {
			Agent: []string { AgentGroupingName },
			Target: []string { TargetGroupingName },
			Metrics: []string { MetricAvg, MetricNumTarget },
		},
	}

	c.Logf("Source JSON: %s", ojson.MarshalJSON(sampleQuery))

	testedJson := ojson.UnmarshalToJsonExt(sampleQuery)

	c.Assert(testedJson.GetPath("time", "start_time"), ocheck.JsonEquals, "897060500", Commentf("\"time.start_time\" is not as expected"))
	c.Assert(testedJson.GetPath("time", "end_time"), ocheck.JsonEquals, "897064500", Commentf("\"time.start_time\" is not as expected"))
	c.Assert(testedJson.GetPath("metrics"), ocheck.JsonEquals, "\"$max >= 150\"", Commentf("\"metrics\" is not as expected"))

	testedAgent := testedJson.Get("agent")
	c.Assert(testedAgent.GetPath("name"), ocheck.JsonEquals, `[ "ag-name-1", "ag-name-2" ]`, Commentf("\"agent.name\" is not as expected"))
	c.Assert(testedAgent.GetPath("hostname"), ocheck.JsonEquals, `[]`, Commentf("\"agent.hostname\" is not as expected"))
	c.Assert(testedAgent.GetPath("ip_address"), ocheck.JsonEquals, `[ "12.10.1", "12.10.2" ]`, Commentf("\"agent.ip_address\" is not as expected"))
	c.Assert(testedAgent.GetPath("connection_id"), ocheck.JsonEquals, `[]`, Commentf("\"agent.connection_id\" is not as expected"))

	c.Assert(testedAgent.GetPath("isps").GetIndex(1), ocheck.JsonEquals, `{ "id": 4, "name": "ag-isp-2", "acronym": "" }`, Commentf("\"agent.isps[1]\" is not as expected"))
	c.Assert(testedAgent.GetPath("provinces").GetIndex(1), ocheck.JsonEquals, `{ "id": 5, "name": "天山省" }`, Commentf("\"agent.provinces[1]\" is not as expected"))
	c.Assert(testedAgent.GetPath("cities").GetIndex(1), ocheck.JsonEquals, `{ "id": 12, "name": "香山市", "post_code": "" }`, Commentf("\"agent.cities[1]\" is not as expected"))
	c.Assert(testedAgent.GetPath("name_tags"), ocheck.JsonEquals, `[]`, Commentf("\"agent.name_tags\" is not as expected"))
	c.Assert(testedAgent.GetPath("group_tags"), ocheck.JsonEquals, `[]`, Commentf("\"agent.group_tags\" is not as expected"))

	testedTarget := testedJson.Get("target")
	c.Assert(testedTarget.GetPath("name"), ocheck.JsonEquals, `[ "tg-name-1", "tg-name-2" ]`, Commentf("\"target.name\" is not as expected"))
	c.Assert(testedTarget.GetPath("host"), ocheck.JsonEquals, `[ "wine-1", "wine-2" ]`, Commentf("\"target.host\" is not as expected"))

	c.Assert(testedTarget.GetPath("isps"), ocheck.JsonEquals, `[]`, Commentf("\"target.isps[1]\" is not as expected"))
	c.Assert(testedTarget.GetPath("provinces"), ocheck.JsonEquals, `[]`, Commentf("\"target.provinces[1]\" is not as expected"))
	c.Assert(testedTarget.GetPath("cities"), ocheck.JsonEquals, `[]`, Commentf("\"target.cities[1]\" is not as expected"))
	c.Assert(testedTarget.GetPath("name_tags").GetIndex(1), ocheck.JsonEquals, `{ "id": 2862, "value": "nt-4" }`, Commentf("\"target.name_tags\" is not as expected"))
	c.Assert(testedTarget.GetPath("group_tags").GetIndex(1), ocheck.JsonEquals, `{ "id": 30072, "name": "gt-4" }`, Commentf("\"target.group_tags\" is not as expected"))

	testedOutput := testedJson.Get("output")
	c.Assert(testedOutput.Get("agent"), ocheck.JsonEquals, `[ "name" ]`, Commentf("\"output.agent\" is not as expected"))
	c.Assert(testedOutput.Get("target"), ocheck.JsonEquals, `[ "name" ]`, Commentf("\"output.target\" is not as expected"))
	c.Assert(testedOutput.Get("metrics"), ocheck.JsonEquals, `[ "avg", "num_target" ]`, Commentf("\"output.metrics\" is not as expected"))
}
