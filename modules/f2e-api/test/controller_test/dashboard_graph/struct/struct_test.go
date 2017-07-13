package test

import (
	"testing"

	dg "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/controller/dashboard_graph"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
)

func TestDashboardGraphCustomStruct(t *testing.T) {
	viper.Set("services", map[string]interface{}{
		"test":  "test123",
		"test2": "test456",
	})
	viper.Set("enable_services", true)
	viper.AddConfigPath(".")
	viper.AddConfigPath("../../../")
	viper.SetConfigName("cfg_test")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error(err.Error())
	}
	gin.SetMode(gin.TestMode)
	log.SetLevel(log.DebugLevel)
	config.InitDB(viper.GetBool("db.db_debug"))
	Convey("test struct a", t, func() {
		tests1 := dg.APIGraphCreateReqDataWithNewScreenInputs{
			ScreenName: "CPU",
			Title:      "aaa2",
			Endpoints:  []string{"e1", "e2"},
			Counters:   []string{"m1", "m2"},
			TimeSpan:   3600,
			GraphType:  "h",
			Method:     "sum",
			TimeRange:  "3h",
			SortBy:     "a-z",
		}
		err := tests1.Check()
		So(err.Error(), ShouldContainSubstring, "already existing")
		tests1.ScreenName = "aaa1"
		tests1.GraphType = "xx"
		err = tests1.Check()
		So(err.Error(), ShouldContainSubstring, "value of graph_type only accept")
		tests1.GraphType = "h"
		tests1.SortBy = "xx"
		err = tests1.Check()
		So(err.Error(), ShouldContainSubstring, "sort_by only accept 'a-z' or 'z-a'")
		tests1.SortBy = "z-a"
		tests1.TimeRange = "3x"
		err = tests1.Check()
		So(err.Error(), ShouldContainSubstring, "time_range only accept")
		tests1.TimeRange = "323213,213"
		err = tests1.Check()
		So(err.Error(), ShouldContainSubstring, "time_range only accept")
		tests1.TimeRange = "3d"
		err = tests1.Check()
		So(err, ShouldEqual, nil)
		tests1.TimeRange = "1496293526,1496379923"
		err = tests1.Check()
		So(err, ShouldEqual, nil)
		tests1.YScale = "1u"
		err = tests1.Check()
		So(err.Error(), ShouldContainSubstring, "not vaild")
		tests1.YScale = "10000.1"
		err = tests1.Check()
		So(err, ShouldEqual, nil)
		tests1.YScale = "10000"
		err = tests1.Check()
		So(err, ShouldEqual, nil)
		tests1.YScale = "10000,10000"
		err = tests1.Check()
		tests1.YScale = "10000,1k"
		err = tests1.Check()
		So(err, ShouldEqual, nil)
		tests1.YScale = "1k"
		err = tests1.Check()
		So(err, ShouldEqual, nil)
		tests1.YScale = "1k,2k"
		err = tests1.Check()
		So(err, ShouldEqual, nil)
		tests1.YScale = "1k,20"
		err = tests1.Check()
		So(err, ShouldEqual, nil)
		tests1.SampleMethod = "notvail"
		err = tests1.Check()
		So(err.Error(), ShouldContainSubstring, "not vaild, only accept")
		tests1.SampleMethod = "MAX"
		err = tests1.Check()
		So(err, ShouldEqual, nil)
	})
}
