package test

import (
	"testing"

	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/utils"
	log "github.com/Sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
)

func TestHash(t *testing.T) {
	viper.Set("salt", "falcontestsalt")
	log.SetLevel(log.DebugLevel)
	Convey("Test Hash method", t, func() {
		val := utils.HashIt("test2")
		So(val, ShouldEqual, "2eef51f314bf9572c49fc8d19913474a")
	})
}
