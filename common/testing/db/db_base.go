//
// The base environment for RDB testing
//
// Flags
//
// This package has pre-defined flags of command:
//
// 	-dsn_mysql - MySQL DSN used to intialize configuration of mysql connection
package db

import (
	"fmt"

	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	tflag "github.com/Cepave/open-falcon-backend/common/testing/flag"
)

// This callback is used to setup a viable database configuration for testing.
type ViableDbConfigFunc func(config *commonDb.DbConfig)

var testFlags *tflag.TestFlags

func getTestFlags() *tflag.TestFlags {
	if testFlags == nil {
		testFlags = tflag.NewTestFlags()
	}

	return testFlags
}

var flagMessage = fmt.Sprintf("Skip MySql Test. Missed property: \"mysql.owl_portal=<dsn>\" or properties in flag \"-owl.test=%s\"", tflag.FeatureHelp(tflag.F_MySql)[0])
