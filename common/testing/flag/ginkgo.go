//
// SkipFactory - Features
//
// "SkipFactory" provides interfaces to build skip functions, which
// use Ginkgo's "Skip()" function.
//
// The "BuildSkipFactory()" function generate a new "SkipFactory" by
// refresh "NewTestFlags()" instance.
//
// SkipFactory - Owl Databases of MySql
//
// You could use "BuildSkipFactoryOfOwlDb()" to retrieve the "SkipFactory" by
// some of build-in databases.
//
// Features
//
// There are various constants, like "F_HttpClient" or "F_MySql",
// to be used in "BuildSkipFactory()".
//
// Skipping with Ginkgo
//
// You could use "FeatureHelpString()" to generate default message for skipping.
//
//	features := F_HttpClient | F_MySql
// 	sf := BuildSkipFactory(features, FeatureHelpString(features))
//
// 	Context("Sometest", sf.PrependBeforeEach(func() {
// 		/* Your test... */
//
// 		It("Something...", func() {
// 			/* Your test... */
// 		})
// 	}))
//
// Skipping with Ginkgo Builder
//
// You could use "FeatureHelpString" to generate default message for skipping.
//
//	features := F_HttpClient | F_MySql
// 	sf := BuildSkipFactory(features, FeatureHelpString(features))
//
// 	NewGinkgoBuilder("Your Context").
// 		It("Test 1", func() {
// 			sf.Skip()
// 			/* Your test... */
// 		}).
// 		ToContext()
// 	}))
//
// Compose SkipFactory
//
// Some modules are depend on complex environments. For example, the "query" module is using
// multiple databases and other modules.
// You could "Compose()" multiple "SkipFactory"s to perform complex checking of testing environments.
//
//	features := F_HttpClient
//	db := OWL_DB_PORTAL | OWL_DB_UIC
//
// 	sf := BuildSkipFactory(features, FeatureHelpString(features))
// 	sdb := BuildSkipFactoryOfOwlDb(db, OwlDbHelpString(db)).Compose(sf)
package flag

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
)

// Defines the interfaces could be used to skip tests in various situations.
//
// 1. Generates a "BeforeEach()":
//
// 	sf.BeforeEachSkip()
//
// 2. Used with "Describe(string, interface{})" or "Context(string, interface{})"
//
//  Context("Context", sf.PrependBeforeEach(func() {
//  	/* Your test... */
//  }))
//
// 3. Used with "It()", "Specify()":
//
// 	It("Context...", func() {
// 		sf.Skip()
// 		/* Your test... */
// 	})
type SkipFactory interface {
	// Generates a function with prepending of "BeforeEach()" block
	PrependBeforeEach(func()) func()
	// Generates a "BeforeEach()" function with skipping
	BeforeEachSkip()
	// Skips current execution directly
	Skip()
	// Composes another SkipFactory
	Compose(SkipFactory) SkipFactory
}

// Builds factory of skipping process.
//
// This function would auto-load "*TestFlags".
func BuildSkipFactory(matchFeatures int, message string, callerSkip ...int) SkipFactory {
	shouldSkip := !MatchFlags(NewTestFlags(), matchFeatures)
	return BuildSkipFactoryByBool(shouldSkip, message, callerSkip...)
}

// Builds Factory of skipping process
func BuildSkipFactoryOfOwlDb(matchDb int, message string, callerSkip ...int) SkipFactory {
	shouldSkip := !MatchFlagsOfOwlDb(NewTestFlags(), matchDb)
	return BuildSkipFactoryByBool(shouldSkip, message, callerSkip...)
}

var objNotSkipFactory = &notSkipFactoryImpl{}

func BuildSkipFactoryByBool(shouldSkip bool, message string, callerSkip ...int) SkipFactory {
	if shouldSkip {
		return &shouldSkipFactoryImpl{message, callerSkip}
	}

	return objNotSkipFactory
}

type shouldSkipFactoryImpl struct {
	message    string
	callerSkip []int
}

func (s *shouldSkipFactoryImpl) BeforeEachSkip() {
	BeforeEach(func() {
		Skip(s.message, s.callerSkip...)
	})
}
func (s *shouldSkipFactoryImpl) Skip() {
	Skip(s.message, s.callerSkip...)
}
func (s *shouldSkipFactoryImpl) PrependBeforeEach(target func()) func() {
	return func() {
		BeforeEach(s.Skip)
		target()
	}
}
func (s *shouldSkipFactoryImpl) Compose(anotherSkip SkipFactory) SkipFactory {
	return &composeSkipFactoryImpl{[]SkipFactory{s, anotherSkip}}
}

type notSkipFactoryImpl struct{}

func (s *notSkipFactoryImpl) BeforeEachSkip() {}
func (s *notSkipFactoryImpl) Skip()           {}
func (s *notSkipFactoryImpl) PrependBeforeEach(target func()) func() {
	return target
}
func (s *notSkipFactoryImpl) Compose(anotherSkip SkipFactory) SkipFactory {
	return &composeSkipFactoryImpl{[]SkipFactory{anotherSkip}}
}

type composeSkipFactoryImpl struct {
	skips []SkipFactory
}

func (s *composeSkipFactoryImpl) BeforeEachSkip() {
	for _, skip := range s.skips {
		skip.BeforeEachSkip()
	}
}
func (s *composeSkipFactoryImpl) Skip() {
	for _, skip := range s.skips {
		skip.Skip()
	}
}
func (s *composeSkipFactoryImpl) PrependBeforeEach(target func()) func() {
	return func() {
		BeforeEach(s.Skip)
		target()
	}
}
func (s *composeSkipFactoryImpl) Compose(anotherSkip SkipFactory) SkipFactory {
	return &composeSkipFactoryImpl{
		append(s.skips, anotherSkip),
	}
}

// Gets help of features, every feature has a corresponding message.
func FeatureHelp(matchFeatures int) []string {
	message := make([]string, 0)

	if matchFeatures&F_HttpClient > 0 {
		message = append(message, "client.http.host=<host> client.http.port=<port>")
	}
	if matchFeatures&F_JsonRpcClient > 0 {
		message = append(message, "client.jsonrpc.host=<host> client.jsonrpc.port=<port>")
	}
	if matchFeatures&F_MySql > 0 {
		message = append(message, "mysql=<dsn>")
	}
	if matchFeatures&F_ItWeb > 0 {
		message = append(message, "it.web.enable=true")
	}

	return message
}

// Gets help of features of string.
//
// This function likes "FeatureHelp(int)" beside joining the messages with a space character.
func FeatureHelpString(matchFeatures int) string {
	return fmt.Sprintf("Need test flags: %s", strings.Join(FeatureHelp(matchFeatures), " "))
}

// Gets help of properties about OWL databases
func OwlDbHelp(matchDbs int) []string {
	messages := make([]string, 0)

	for db, propName := range owlDbMap {
		if matchDbs&db > 0 {
			messages = append(messages, fmt.Sprintf("%s=<db_conn>"), propName)
		}
	}

	return messages
}

// Gets help of properties about OWL databases as string
//
// This function likes "OwlDbHelp(int)" beside joining the messages with a space character.
func OwlDbHelpString(matchDbs int) string {
	return fmt.Sprintf("Need test flags: %s", strings.Join(OwlDbHelp(matchDbs), " "))
}

// Checks the match features on "*TestFlags"
func MatchFlags(sourceFlags *TestFlags, matchFeatures int) bool {
	if (matchFeatures&F_HttpClient > 0) && !sourceFlags.HasHttpClient() {
		return false
	}
	if (matchFeatures&F_JsonRpcClient > 0) && !sourceFlags.HasJsonRpcClient() {
		return false
	}
	if (matchFeatures&F_MySql > 0) && !sourceFlags.HasMySql() {
		return false
	}
	if (matchFeatures&F_ItWeb > 0) && !sourceFlags.HasItWeb() {
		return false
	}

	return true
}

// Checks the match db on "*TestFlags"
func MatchFlagsOfOwlDb(sourceFlags *TestFlags, matchDbs int) bool {
	for db := range owlDbMap {
		currentFlag := matchDbs & db
		if currentFlag > 0 && !sourceFlags.HasMySqlOfOwlDb(currentFlag) {
			return false
		}
	}

	return true
}
