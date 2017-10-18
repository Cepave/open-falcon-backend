package flag

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	gb "github.com/Cepave/open-falcon-backend/common/testing/ginkgo/builder"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("New object of *TestFlags", func() {
	Context("Calling multiple times", func() {
		It("1st calling, everything should be fine", func() {
			testedFlags := NewTestFlags()
			Expect(testedFlags.HasItWeb()).To(BeFalse())
		})

		It("2nd calling, should be same of 1st calling", func() {
			testedFlags := NewTestFlags()
			Expect(testedFlags.HasItWeb()).To(BeFalse())
		})
	})
})

var _ = Describe("Set-up TestFlags by Viper", func() {
	var (
		viperObj  *viper.Viper
		testFlags *TestFlags
	)

	var propDesc = func(propName string, _ interface{}) string {
		return fmt.Sprintf("Check prop: [%s]", propName)
	}

	gb.NewGinkgoBuilder("Build-in Flags").
		BeforeFirst(func() {
			viperObj = viper.New()

			viperObj.Set("mysql", "mysql@ccc")
			viperObj.Set("client.http.host", "pc01.com")
			viperObj.Set("client.http.port", "10")
			viperObj.Set("client.http.ssl", "true")
			viperObj.Set("client.http.resource", "/nqm")
			viperObj.Set("client.jsonrpc.host", "pc02.com")
			viperObj.Set("client.jsonrpc.port", "33")
			viperObj.Set("it.web.enable", "true")

			testFlags = newTestFlags(viperObj)
		}).Table(gb.NewGinkgoTable().
		Exec(func(propName string, expectedValue interface{}) {
			Expect(testFlags.typedFlags[propName]).To(Equal(expectedValue))
		}).
		Case(propDesc, "mysql", "mysql@ccc").
		Case(propDesc, "client.http.host", "pc01.com").
		Case(propDesc, "client.http.port", uint16(10)).
		Case(propDesc, "client.http.ssl", true).
		Case(propDesc, "client.http.resource", "/nqm").
		Case(propDesc, "client.jsonrpc.host", "pc02.com").
		Case(propDesc, "client.jsonrpc.port", uint16(33)).
		Case(propDesc, "it.web.enable", true),
	).
		ToContext()

	gb.NewGinkgoBuilder("Other Flags").
		BeforeFirst(func() {
			viperObj = viper.New()

			viperObj.Set("client.http.host", "zz.pc01.com")
			viperObj.Set("something.cc", "cool!")

			testFlags = newTestFlags(viperObj)
		}).Table(gb.NewGinkgoTable().
		Exec(func(propName string, expectedValue string) {
			Expect(testFlags.GetViper().GetString(propName)).To(Equal(expectedValue))
		}).
		Case(propDesc, "client.http.host", "zz.pc01.com").
		Case(propDesc, "something.cc", "cool!"),
	).
		ToContext()
})

var _ = Describe("multiPropLoader", func() {
	Context("calling of loadProperties()", func() {
		gb.NewGinkgoBuilder("Different Separators").
			Table(gb.NewGinkgoTable().
				Exec(func(propString, separator string) {
					testedLoader := newMultiPropLoader()
					testedLoader.loadProperties(
						"", propString, separator,
					)

					testedViper := testedLoader.viperObj

					expectFunc := func(propName, expected string) {
						GinkgoT().Logf("Asserts prop[%s]", propName)
						Expect(testedViper.GetString(propName)).To(Equal(expected))
					}

					expectFunc("a.v1", "33")
					expectFunc("a.v2", "99")
					expectFunc("a.v3", "hh")
					expectFunc("a.v4", "jj")
				}).
				Case(fmt.Sprintf("Default separator(%s)", DEFAULT_SEPARATOR), "a.v1=33 a.v2=99   a.v3=hh\na.v4=jj", DEFAULT_SEPARATOR).
				Case("Customized separator(!!)", "a.v1=33!!a.v2=99!!a.v3=hh!!a.v4=jj", "!!"),
			).
			ToContext()

		gb.NewGinkgoBuilder("String properties overrides property file").
			It("Overrode value should match expectd", func() {
				testedLoader := newMultiPropLoader()
				testedLoader.loadProperties(
					"./sample.properties", "a.v2=99", DEFAULT_SEPARATOR,
				)

				testedViper := testedLoader.viperObj

				expectFunc := func(propName, expected string) {
					GinkgoT().Logf("Asserts prop[%s]", propName)
					Expect(testedViper.GetString(propName)).To(Equal(expected))
				}

				expectFunc("a.v2", "99") // Overrides property file
				expectFunc("b.v1", "40") // From property file
			}).
			ToContext()
	})

	gb.NewGinkgoBuilder("Load from environment variables").
		BeforeEach(func() {
			os.Setenv(ENV_OWL_TEST_PROPS_FILE, "sample.properties")
			os.Setenv(ENV_OWL_TEST_PROPS, "a1=11!!a2=33")
			os.Setenv(ENV_OWL_TEST_PROPS_SEP, "!!")
		}).
		AfterEach(func() {
			os.Unsetenv(ENV_OWL_TEST_PROPS_FILE)
			os.Unsetenv(ENV_OWL_TEST_PROPS)
			os.Unsetenv(ENV_OWL_TEST_PROPS_SEP)
		}).
		It("Final properties should match expected", func() {
			testedLoader := newMultiPropLoader()
			testedLoader.loadFromEnv()

			testedViper := testedLoader.viperObj

			expectFunc := func(propName, expected string) {
				GinkgoT().Logf("Asserts prop[%s]", propName)
				Expect(testedViper.GetString(propName)).To(Equal(expected))
			}

			expectFunc("a1", "11")
			expectFunc("a2", "33")
			expectFunc("b.v1", "40")
		}).
		ToContext()
})

var _ = Describe("Checking function of TestFlags", func() {
	const nonViable = "<!Non-Viable!>"
	setIfViable := func(viperObj *viper.Viper, propName, propValue string) {
		if propValue != nonViable {
			viperObj.Set(propName, propValue)
		}
	}

	gb.NewGinkgoBuilder("HasHttpClient()").
		Table(gb.NewGinkgoTable().
			Exec(func(sampleHost string, samplePort string, expectedResult bool) {
				viperObj := viper.New()

				setIfViable(viperObj, "client.http.host", sampleHost)
				setIfViable(viperObj, "client.http.port", samplePort)

				testFlags := newTestFlags(viperObj)

				Expect(testFlags.HasHttpClient()).To(Equal(expectedResult))
			}).
			Case("viable host and viable port", "cp01.hu.com", "6181", true).
			Case("Non-viable host and viable port", nonViable, "6182", false).
			Case("Viable host and non-viable port", "rpc02.kc.com", nonViable, false).
			Case("Empty host and viable port", "", "6183", false).
			Case("Viable host and empty port", "cp03.aa.com", "", false).
			Case("Viable host and \"0\" port", "cp03.aa.com", "0", false),
		).
		ToContext()

	gb.NewGinkgoBuilder("HasMySql()").
		Table(gb.NewGinkgoTable().
			Exec(func(sampleMysql string, expectedResult bool) {
				viperObj := viper.New()

				setIfViable(viperObj, "mysql", sampleMysql)

				testFlags := newTestFlags(viperObj)

				Expect(testFlags.HasMySql()).To(Equal(expectedResult))
			}).
			Case("mysql is viable", "root:ccc@tcp(127.0.0.1:3306)/ddbb", true).
			Case("mysql is empty", nonViable, false).
			Case("mysql is non-viable", "", false),
		).
		ToContext()

	gb.NewGinkgoBuilder("HasItWeb()").
		Table(gb.NewGinkgoTable().
			Exec(func(sampleMysql string, expectedResult bool) {
				viperObj := viper.New()

				setIfViable(viperObj, "it.web.enable", sampleMysql)

				testFlags := newTestFlags(viperObj)

				Expect(testFlags.HasItWeb()).To(Equal(expectedResult))
			}).
			Case("it.web.enable is true", "true", true).
			Case("it.web.enable is false", "false", false).
			Case("it.web.enable is non-viable", nonViable, false).
			Case("it.web.enable is empty", "", false),
		).
		ToContext()

	gb.NewGinkgoBuilder("HasJsonRpcClient()").
		Table(gb.NewGinkgoTable().
			Exec(func(sampleHost string, samplePort string, expectedResult bool) {
				viperObj := viper.New()

				setIfViable(viperObj, "client.jsonrpc.host", sampleHost)
				setIfViable(viperObj, "client.jsonrpc.port", samplePort)

				testFlags := newTestFlags(viperObj)

				Expect(testFlags.HasJsonRpcClient()).To(Equal(expectedResult))
			}).
			Case("viable host and viable port", "rpc01.hu.com", "10022", true).
			Case("Non-viable host and viable port", nonViable, "10023", false).
			Case("Viable host and non-viable port", "rpc02.kc.com", nonViable, false).
			Case("Empty host and viable port", "", "10101", false).
			Case("Viable host and empty port", "rpc03.aa.com", "", false).
			Case("Viable host and \"0\" port", "kp03.aa.com", "0", false),
		).
		ToContext()
})
