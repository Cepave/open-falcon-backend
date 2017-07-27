package nqm

import (
	"flag"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

func init() {
	flag.Parse()
}
