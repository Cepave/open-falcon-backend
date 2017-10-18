//
// Provides unified interface to access needed flags when you are testing.
//
// Property file
//
// "NewTestFlags()" has two way to load property file:
//
// 	OWL_TEST_PROPS_FILE - Environment Variable
// 	"-owl.test.propfile" - flag of "go test"
//
// Entry Environment Variables
//
// Following environment variables are supported when calling "NewTestFlags()".
//
// 	OWL_TEST_PROPS - As same as "-owl.test"
// 	OWL_TEST_PROPS_SEP - As same as "-owl.test.sep"
//
// Entry Flags
//
// There are only two flags of Golang needed:
//
// 	-owl.test=<properties>
// 	-owl.test.sep=<separator for properties>
//
// The format used by "owl.test" is property file:
//
// https://en.wikipedia.org/wiki/.properties
//
// In order to separate properties, "owl.test.sep"(as regular expression) would be used to
// recognize a record of property file.
//
// See "DEFAULT_SEPARATOR" constant for default separator.
//
// Loading priority
//
// When you execute go test with viable values for both of the environment variables and the flags,
// the priority is:
//
// 	1. Load file from environment variable(OWL_TEST_PROPS_FILE)
// 	2. Load properties of environment variable(OWL_TEST_PROPS)
// 	3. Load file from flag("-owl.test.propfile")
// 	4. Load properties from flag("-owl.test")
//
// Pre-defined Properties - Features
//
// There are some pre-defined properties:
//
// 	mysql - MySql connection
//
// 	client.http.host - HTTP client
// 	client.http.port - HTTP client
// 	client.http.ssl - HTTP client
// 	client.http.resource - HTTP client
//
// 	client.jsonrpc.host - JSONRPC Client
// 	client.jsonrpc.port - JSONRPC Client
//
// 	it.web.enable - IT to Web
//
// The object of "*TestFlags" provides various functions to check
// whether or not some configuration for testing are enabled.
//
// For example, "HasMySql()" would let you know whether "mysql=<conn>" is viable.
//
// Pre-defined Properties - Owl Databases of MySql
//
// Following list shows build-in supporting databases of Owl Database:
//
// 	mysql.owl_portal - MySql connection on OWL-Portal
// 	mysql.owl_graph - MySql connection on OWL-Graph
// 	mysql.owl_uic - MySql connection on OWL-Uic
// 	mysql.owl_links - MySql connection on OWL-Links
// 	mysql.owl_grafana - MySql connection on OWL-Grafana
// 	mysql.owl_dashboard - MySql connection on OWL-Dashboard
// 	mysql.owl_boss - MySql connection on OWL-Boss
//
// You could use "HasMySqlOfOwlDb(int)" or "GetMysqlOfOwlDb(int)" to retrieve value of properties.
//
// Constraint
//
// The empty string of property value would be considered as non-viable.
package flag

import (
	"flag"
	"fmt"
	"regexp"
	"strings"

	"github.com/juju/errors"
	"github.com/spf13/viper"
)

const (
	// Default separator
	DEFAULT_SEPARATOR = "\\s+"
)

/*
Bit reservation principals:

	Bits (0~7): For clients of various protocols
	Bits (8~15): For databases
	Bits (16~23): For misc(e.x. mocking server)
*/
const (
	// Feature of HTTP client
	F_HttpClient = 0x01
	// Feature of JSONRPC client
	F_JsonRpcClient = 0x02
	// Feature of MySql
	F_MySql = 0x100
	// Feature of IT web
	F_ItWeb = 0x10000
)

const (
	ENV_OWL_TEST_PROPS      = "OWL_TEST_PROPS"
	ENV_OWL_TEST_PROPS_SEP  = "OWL_TEST_PROPS_SEP"
	ENV_OWL_TEST_PROPS_FILE = "OWL_TEST_PROPS_FILE"
)

const (
	OWL_DB_PORTAL    = 0x01
	OWL_DB_GRAPH     = 0x02
	OWL_DB_UIC       = 0x04
	OWL_DB_LINKS     = 0x8
	OWL_DB_GRAFANA   = 0x10
	OWL_DB_DASHBOARD = 0x20
	OWL_DB_BOSS      = 0x40
)

var owlDbMap = map[int]string{
	OWL_DB_PORTAL:    "mysql.owl_portal",
	OWL_DB_GRAPH:     "mysql.owl_graph",
	OWL_DB_UIC:       "mysql.owl_uic",
	OWL_DB_LINKS:     "mysql.owl_links",
	OWL_DB_GRAFANA:   "mysql.owl_grafana",
	OWL_DB_DASHBOARD: "mysql.owl_dashboard",
	OWL_DB_BOSS:      "mysql.owl_boss",
}

var (
	owlTest         = flag.String("owl.test", "", "Owl typedFlags for testing properties")
	owlTestSep      = flag.String("owl.test.sep", DEFAULT_SEPARATOR, "Owl typedFlags for separator of properties")
	owlTestPropFile = flag.String("owl.test.propfile", "", "Owl property file for testing")
)

var singletonTestFlags *TestFlags

// Initializes the object of "*TestFlags" by parsing flag automatically.
//
// This function parses os.Args every time it is get called.
func NewTestFlags() *TestFlags {
	if singletonTestFlags != nil {
		return singletonTestFlags.clone()
	}

	viperLoader := newMultiPropLoader()
	viperLoader.loadFromEnv()
	viperLoader.loadFromFlag()

	/**
	 * Setup Flags of testing
	 */
	newFlags := newTestFlags(viperLoader.viperObj)
	// :~)

	singletonTestFlags = newFlags

	return singletonTestFlags.clone()
}

// Convenient type used to access specific testing environment of OWL.
type TestFlags struct {
	typedFlags map[string]interface{}

	viperObj *viper.Viper
}

func (f *TestFlags) GetViper() *viper.Viper {
	return f.viperObj
}

// Gets property value of "mysql"
func (f *TestFlags) GetMySql() string {
	if f.HasMySql() {
		return f.typedFlags["mysql"].(string)
	}

	return ""
}

func (f *TestFlags) GetMysqlOfOwlDb(owlDb int) string {
	propName, ok := owlDbMap[owlDb]
	if !ok {
		panic(fmt.Sprintf("Unsupported OWL Db: %v", owlDb))
	}

	return strings.TrimSpace(f.viperObj.GetString(propName))
}

// Gets property values of:
// 	client.http.host
// 	client.http.port
// 	client.http.resource
// 	client.http.ssl
func (f *TestFlags) GetHttpClient() (string, uint16, string, bool) {
	if f.HasHttpClient() {
		return f.typedFlags["client.http.host"].(string), f.typedFlags["client.http.port"].(uint16),
			f.typedFlags["client.http.resource"].(string), f.typedFlags["client.http.ssl"].(bool)
	}

	return "", 0, "", false
}

// Gets property values of "client.jsonrpc.host" and "client.jsonrpc.port"
func (f *TestFlags) GetJsonRpcClient() (string, uint16) {
	if f.HasJsonRpcClient() {
		return f.typedFlags["client.jsonrpc.host"].(string), f.typedFlags["client.jsonrpc.port"].(uint16)
	}

	return "", 0
}

// Gives "true" if and only if following properties are viable:
//
// 	client.jsonrpc.host=
// 	client.jsonrpc.port=
//
// Example:
// 	"-owl.flag=client.jsonrpc.host=127.0.0.1 client.jsonrpc.port=3396"
func (f *TestFlags) HasJsonRpcClient() bool {
	_, hostOk := f.typedFlags["client.jsonrpc.host"]
	_, portOk := f.typedFlags["client.jsonrpc.port"]

	return hostOk && portOk
}

// Gives "true" if and only if following properties are viable:
//
// 	client.http.host=
// 	client.http.port=
//
// Example:
// 	"-owl.flag=client.http.host=127.0.0.1 client.http.port=3396"
func (f *TestFlags) HasHttpClient() bool {
	_, hostOk := f.typedFlags["client.http.host"]
	_, portOk := f.typedFlags["client.http.port"]

	return hostOk && portOk
}

// Gives "true" if and only if "mysql" property is non-empty
//
// Example:
// 	"-owl.flag=mysql=root:cepave@tcp(192.168.20.50:3306)/falcon_portal_test?parseTime=True&loc=Local"
func (f *TestFlags) HasMySql() bool {
	_, ok := f.typedFlags["mysql"]
	return ok
}

// Gives "true" if and only if "mysql.<db>" property is non-empty
//
// Example:
// 	"-owl.flag=mysql.portal=root:cepave@tcp(192.168.20.50:3306)/falcon_portal_test?parseTime=True&loc=Local"
func (f *TestFlags) HasMySqlOfOwlDb(owlDb int) bool {
	return f.GetMysqlOfOwlDb(owlDb) != ""
}

// Gives "true" if and only if "it.web.enable" property is true
//
// Example:
// 	"-owl.flag=it.web.enable=true"
func (f *TestFlags) HasItWeb() bool {
	hasItWeb, ok := f.typedFlags["it.web.enable"].(bool)
	return ok && hasItWeb
}

func (f *TestFlags) clone() *TestFlags {
	newViper := viper.New()

	for k, v := range f.viperObj.AllSettings() {
		newViper.Set(k, v)
	}

	newTypedFlags := make(map[string]interface{})
	for k, v := range f.typedFlags {
		newTypedFlags[k] = v
	}

	return &TestFlags{
		newTypedFlags, newViper,
	}
}

func (f *TestFlags) setupByViper() {
	viperObj := f.viperObj

	/**
	 * MySql
	 */
	if viperObj.IsSet("mysql") {
		setNonEmptyString(f.typedFlags, "mysql", viperObj)
	}
	// :~)

	/**
	 * HTTP client
	 */
	if viperObj.IsSet("client.http.host") {
		setNonEmptyString(f.typedFlags, "client.http.host", viperObj)
	}
	if viperObj.IsSet("client.http.port") {
		setValidPort(f.typedFlags, "client.http.port", viperObj)
	}
	viperObj.SetDefault("client.http.ssl", false)
	f.typedFlags["client.http.ssl"] = viperObj.GetBool("client.http.ssl")
	viperObj.SetDefault("client.http.resource", "")
	f.typedFlags["client.http.resource"] = viperObj.GetString("client.http.resource")
	// :~)

	/**
	 * JSONRPC Client
	 */
	if viperObj.IsSet("client.jsonrpc.host") {
		setNonEmptyString(f.typedFlags, "client.jsonrpc.host", viperObj)
	}
	if viperObj.IsSet("client.jsonrpc.port") {
		setValidPort(f.typedFlags, "client.jsonrpc.port", viperObj)
	}
	// :~)

	/**
	 * Start web for integration test
	 */
	if viperObj.IsSet("it.web.enable") {
		f.typedFlags["it.web.enable"] = viperObj.GetBool("it.web.enable")
	}
	// :~)
}

func setValidPort(props map[string]interface{}, key string, viperObj *viper.Viper) {
	v := viperObj.GetInt(key)
	if v > 0 {
		props[key] = uint16(v)
	}
}
func setNonEmptyString(props map[string]interface{}, key string, viperObj *viper.Viper) {
	value := strings.TrimSpace(viperObj.GetString(key))
	if len(value) > 0 {
		props[key] = value
	}
}

func newTestFlags(viperObj *viper.Viper) *TestFlags {
	testFlag := &TestFlags{
		make(map[string]interface{}),
		viperObj,
	}
	testFlag.setupByViper()
	return testFlag
}

func newMultiPropLoader() *multiPropLoader {
	viperObj := viper.New()
	viperObj.SetConfigType("properties")

	return &multiPropLoader{viperObj}
}

type multiPropLoader struct {
	viperObj *viper.Viper
}

func (l *multiPropLoader) loadFromEnv() {
	envLoader := viper.New()
	envLoader.SetDefault(ENV_OWL_TEST_PROPS_SEP, DEFAULT_SEPARATOR)
	envLoader.BindEnv(ENV_OWL_TEST_PROPS)
	envLoader.BindEnv(ENV_OWL_TEST_PROPS_SEP)
	envLoader.BindEnv(ENV_OWL_TEST_PROPS_FILE)

	l.loadProperties(
		envLoader.GetString(ENV_OWL_TEST_PROPS_FILE),
		envLoader.GetString(ENV_OWL_TEST_PROPS),
		envLoader.GetString(ENV_OWL_TEST_PROPS_SEP),
	)
}
func (l *multiPropLoader) loadFromFlag() {
	if !flag.Parsed() {
		flag.Parse()
	}

	l.loadProperties(*owlTestPropFile, *owlTest, *owlTestSep)
}

func (l *multiPropLoader) loadProperties(filename string, propertiesString string, separator string) {
	/**
	 * Loads property file into viper object
	 */
	if filename != "" {
		l.viperObj.SetConfigFile(filename)
		if err := l.viperObj.MergeInConfig(); err != nil {
			logger.Warnf("Load property file[%s] has error: %v", filename, err)
		}
	}
	// :~)

	propertiesString = strings.TrimSpace(propertiesString)
	if propertiesString == "" {
		return
	}

	splitRegExp := regexp.MustCompile(separator)
	properties := splitRegExp.ReplaceAllString(propertiesString, "\n")

	/**
	 * Loads properties into viper object
	 */
	if err := errors.Annotate(
		l.viperObj.MergeConfig(strings.NewReader(properties)),
		"Read owl.test as format of property file has error",
	); err != nil {
		panic(errors.Details(err))
	}
	// :~)
}
