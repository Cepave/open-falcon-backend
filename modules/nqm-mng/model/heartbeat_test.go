package model

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TestHeartbeatSuite struct{}

var _ = Suite(&TestHeartbeatSuite{})
