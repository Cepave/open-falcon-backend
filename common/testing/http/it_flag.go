package http

import (
	"flag"

	gk "github.com/onsi/ginkgo"
)

var itFlag = flag.Bool("it.web.enable", false, "Need -it.web.enable=true to execute it")

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
