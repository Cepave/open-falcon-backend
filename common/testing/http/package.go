package http

import (
	"github.com/Cepave/open-falcon-backend/common/logruslog"

	oflag "github.com/Cepave/open-falcon-backend/common/testing/flag"
)

var logger = logruslog.NewDefaultLogger("INFO")

var testFlags *oflag.TestFlags

func getTestFlags() *oflag.TestFlags {
	if testFlags == nil {
		testFlags = oflag.NewTestFlags()
	}

	return testFlags
}
