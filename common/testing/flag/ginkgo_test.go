package flag

import (
	"fmt"

	"github.com/spf13/viper"

	gb "github.com/Cepave/open-falcon-backend/common/testing/ginkgo/builder"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BuildSkipFactoryByBool()", func() {
	Context("Skipping", func() {
		failedTest := func() {
			It("Should be skipped", func() {
				Expect(1).To(Equal(2))
			})
		}

		testedSkip := BuildSkipFactoryByBool(true, "Should skipping")

		Context("PrependBeforeEach()", testedSkip.PrependBeforeEach(func() {
			failedTest()
		}))

		Context("BeforeEachSkip()", func() {
			testedSkip.BeforeEachSkip()
			failedTest()
		})

		Context("Skip()", func() {
			It("Should skipped", func() {
				testedSkip.Skip()
				Expect(1).To(Equal(2))
			})
		})
	})

	Context("Not-skipping", func() {
		testedSkip := BuildSkipFactoryByBool(false, "Nothing skipping")

		Context("PrependBeforeEach()", func() {
			touched := 0

			Context("Touch Execution", testedSkip.PrependBeforeEach(func() {
				It("Sample Touch", func() {
					touched = 1
				})
			}))

			It("Touched should be 1", func() {
				Expect(touched).To(Equal(1))
			})
		})

		Context("BeforeEachSkip()", func() {
			touched := 0

			Context("Touch Execution", func() {
				testedSkip.BeforeEachSkip()
				It("Sample Touch", func() {
					touched = 1
				})
			})

			It("Touched should be 1", func() {
				Expect(touched).To(Equal(1))
			})
		})

		Context("Skip()", func() {
			touched := 0

			It("Sample Touch", func() {
				testedSkip.Skip()
				touched = 1
			})

			It("Touched should be 1", func() {
				Expect(touched).To(Equal(1))
			})
		})
	})
})

var _ = Describe("BuildSkipFactory()", func() {
	Context("Skipping", func() {
		It("The type of implementation should be \"shouldSkipFactoryImpl\"", func() {
			skipFactory := BuildSkipFactory(F_JsonRpcClient, FeatureHelpString(F_JsonRpcClient))

			_, ok := interface{}(skipFactory).(*shouldSkipFactoryImpl)
			Expect(ok).To(BeTrue())
		})
	})
})

var _ = Describe("BuildSkipFactoryOfOwlDb()", func() {
	Context("Skipping", func() {
		It("The type of implementation should be \"shouldSkipFactoryImpl\"", func() {
			skipFactory := BuildSkipFactoryOfOwlDb(OWL_DB_LINKS, OwlDbHelpString(OWL_DB_LINKS))

			_, ok := interface{}(skipFactory).(*shouldSkipFactoryImpl)
			Expect(ok).To(BeTrue())
		})
	})
})

var _ = Describe("MatchFlags()", func() {
	Context("Features and flags", gb.NewGinkgoTable().
		Exec(func(viperSetup func(*viper.Viper), sampleFeatures int, expected bool) {
			/**
			 * Set-up TestFlags
			 */
			viperObj := viper.New()
			viperSetup(viperObj)
			sampleFlags := newTestFlags(viperObj)
			// :~)

			testedResult := MatchFlags(sampleFlags, sampleFeatures)
			Expect(testedResult).To(Equal(expected))
		}).
		Case("Nothing(no feature)", func(viperObj *viper.Viper) {}, 0, true).
		Case("Enable of F_MySql", func(viperObj *viper.Viper) {
			viperObj.Set("mysql", "aaa")
		}, F_MySql, true).
		Case("Enable of F_HttpClient", func(viperObj *viper.Viper) {
			viperObj.Set("client.http.host", "acb.com")
			viperObj.Set("client.http.port", "10104")
		}, F_HttpClient, true).
		Case("Enable of F_JsonRpcClient", func(viperObj *viper.Viper) {
			viperObj.Set("client.jsonrpc.host", "acb02.com")
			viperObj.Set("client.jsonrpc.port", "9104")
		}, F_JsonRpcClient, true).
		Case("Enable of F_ItWeb", func(viperObj *viper.Viper) {
			viperObj.Set("it.web.enable", "true")
		}, F_ItWeb, true).
		Case("F_MySql | F_ItWeb. it.web.enable is not enable", func(viperObj *viper.Viper) {
			viperObj.Set("mysql", "aaa")
		}, F_MySql|F_ItWeb, false).
		Case("F_MySql | F_ItWeb. All of flags are enabled", func(viperObj *viper.Viper) {
			viperObj.Set("mysql", "aaa")
			viperObj.Set("it.web.enable", "true")
		}, F_MySql|F_ItWeb, true).
		ToFunc(),
	)
})

var _ = Describe("MatchFlagsOfOwlDb", func() {
	Context("Db and flags", gb.NewGinkgoTable().
		Exec(func(viperSetup func(*viper.Viper), sampleDbs int, expected bool) {
			/**
			 * Set-up TestFlags
			 */
			viperObj := viper.New()
			viperSetup(viperObj)
			sampleFlags := newTestFlags(viperObj)
			// :~)

			testedResult := MatchFlagsOfOwlDb(sampleFlags, sampleDbs)
			Expect(testedResult).To(Equal(expected))
		}).
		Case("Nothing(no feature)", func(viperObj *viper.Viper) {}, 0, true).
		Case("Enable of OWL_DB_GRAPH", func(viperObj *viper.Viper) {
			viperObj.Set("mysql.owl_graph", "conn-ok")
		}, OWL_DB_GRAPH, true).
		Case("Disable of OWL_DB_GRAPH", func(viperObj *viper.Viper) {
			viperObj.Set("mysql.owl_portal", "conn-ok")
		}, OWL_DB_GRAPH, false).
		Case("Enable of OWL_DB_GRAPH and OWL_DB_UIC", func(viperObj *viper.Viper) {
			viperObj.Set("mysql.owl_graph", "conn-ok")
			viperObj.Set("mysql.owl_uic", "conn-ok")
		}, OWL_DB_GRAPH|OWL_DB_UIC, true).
		Case("Enable of OWL_DB_GRAPH and disable of OWL_DB_UIC", func(viperObj *viper.Viper) {
			viperObj.Set("mysql.owl_graph", "conn-ok")
		}, OWL_DB_GRAPH|OWL_DB_UIC, false).
		ToFunc(),
	)
})

var _ = Describe("composeSkipFactoryImpl", func() {
	composeBuilder := func(skip_1 bool, skip_2 bool, others ...bool) SkipFactory {
		factory := BuildSkipFactoryByBool(skip_1, "1st").Compose(
			BuildSkipFactoryByBool(skip_2, "2nd"),
		)

		for i, otherSkip := range others {
			factory = factory.Compose(BuildSkipFactoryByBool(otherSkip, fmt.Sprintf("Other: %d", i+1)))
		}

		return factory
	}

	Context("first is skipped", func() {
		testedSkip := composeBuilder(true, false, false)
		It("Should be skipped", func() {
			testedSkip.Skip()
			Expect(1).To(Equal(2))
		})
	})

	Context("Last is skipped", func() {
		testedSkip := composeBuilder(false, false, true)

		It("Should be skipped", func() {
			testedSkip.Skip()
			Expect(1).To(Equal(2))
		})
	})

	Context("Nothing is skipped", func() {
		var notSkipped = 0

		testedSkip := composeBuilder(false, false, false)

		It("Should not be skipped", func() {
			testedSkip.Skip()
			notSkipped = 1
		})
		It("Value of notSkipped should be 1", func() {
			Expect(notSkipped).To(Equal(1))
		})
	})
})
