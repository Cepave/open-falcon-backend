package testing

import (
	"time"
	"gopkg.in/check.v1"
)

func ParseTime(c *check.C, timeAsString string) time.Time {
	timeValue, err := time.Parse(time.RFC3339, timeAsString)
	c.Assert(err, check.IsNil)

	return timeValue
}
