//
// Control Flag
//
// The flag "it.web.enable"(default is false) controls whether or not to skip tests.
//
// If you like to control whether or not to skip tests in Ginkgo,
// the GinkgoHttpIt variable provides enhancement for out-of-box utility over Ginkgo framework.
//
package http

import (
	"flag"

	gk "github.com/onsi/ginkgo"
)

var itFlag = flag.Bool("it.web.enable", false, "Need -it.web.enable=true to execute it")

// Ginkgo functions for testing.
//
// NeedItWeb(func())
//
//	Ginkgo wrapper for skipping tests if flag "it.web.enable" doesn't appear in command line.
//
//	var needWebFlag = tHttp.GinkgoHttpIt.NeedItWeb
//
// 	var _ = Describe("Some web test", needWebFlag(func() {
// 		It("Something should be ....", func() {
// 			/* your test... */
// 		})
// 	}))
var GinkgoHttpIt = &struct {
	NeedItWeb func(func()) func()
}{
	NeedItWeb: func(src func()) func() {
		return func() {
			gk.BeforeEach(func() {
				if *itFlag == false {
					gk.Skip("-it.web.enable is false. Skip testing")
				}
			})

			src()
		}
	},
}
