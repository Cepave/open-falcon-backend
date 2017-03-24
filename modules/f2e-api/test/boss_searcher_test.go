package test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/helper/filter"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/boss"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBossSeacher(t *testing.T) {
	dat, err := ioutil.ReadFile("./fakeData.json")
	if err != nil {
		log.Fatalln(err.Error())
	}
	var testData []boss.BossHost
	err = json.Unmarshal(dat, &testData)
	if err != nil {
		log.Println(err.Error())
	}

	Convey("search platform", t, func() {
		res := filter.PlatformFilter(testData, "02")
		So(len(res), ShouldEqual, 3)
		res = filter.PlatformFilter(testData, "01")
		So(len(res), ShouldEqual, 13)
	})

	Convey("search isp", t, func() {
		res := filter.IspFilter(testData, "ctt")
		So(len(res), ShouldEqual, 8)
	})

	Convey("search idc", t, func() {
		res := filter.IdcFilter(testData, "北京一区讯通")
		So(len(res), ShouldEqual, 11)
	})

	Convey("search ip", t, func() {
		res := filter.IpFilter(testData, ".86")
		So(len(res), ShouldEqual, 2)
	})

	Convey("search hostname", t, func() {
		res := filter.HostNameFilter(testData, "51-")
		So(len(res), ShouldEqual, 2)
	})

}
