package testing

import (
	"time"

	ojson "github.com/Cepave/open-falcon-backend/common/json"
	. "github.com/onsi/gomega"
	"gopkg.in/check.v1"
)

func ParseTimeByGinkgo(timeAsString string) time.Time {
	timeValue, err := time.Parse(time.RFC3339, timeAsString)

	Ω(err).Should(Succeed())

	return timeValue
}

func ParseTimeToJsonTimeByGinkgo(timeAsString string) ojson.JsonTime {
	return ojson.JsonTime(ParseTimeByGinkgo(timeAsString))
}

func ParseTime(c *check.C, timeAsString string) time.Time {
	timeValue, err := time.Parse(time.RFC3339, timeAsString)
	c.Assert(err, check.IsNil)

	return timeValue
}
func ParseTimeToJsonTime(c *check.C, timeAsString string) ojson.JsonTime {
	return ojson.JsonTime(ParseTime(c, timeAsString))
}
