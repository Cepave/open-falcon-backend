package conf

import (
	"fmt"
	"strings"
	"testing"

	"os"

	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
)

func TestConfRead(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootpath := strings.Replace(pwd, "lambda_extends/conf", "", 1)
	viper.Set("lambda_extends.root_dir", rootpath)
	ReadConf()
	c := Config()
	testfunname := "avgCompare"
	Convey("ReadFuncation", t, func() {
		Convey("ReadJsonConfTest", func() {
			So(c["avgCompare"].FuncationName, ShouldNotBeEmpty)
			log.Printf("%v", c["top"].Params)
			So(len(c[testfunname].Params), ShouldBeGreaterThan, -1)
			So(c["top"].Params[0], ShouldEqual, "limit:int")
		})
	})

	Convey("JsReadFuncation", t, func() {
		Convey("ReadJsonFileTest", func() {
			contain := jsFileReader("../js/avgCompare.js")
			So(contain, ShouldNotBeEmpty)
		})
		Convey("GetAvaibleFun", func() {
			funcations := GetAvaibleFun()
			So(len(funcations), ShouldBeGreaterThan, 0)
		})
	})

}
