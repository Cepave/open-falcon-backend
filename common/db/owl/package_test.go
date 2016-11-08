package owl

import (
	"flag"
	"testing"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

func init() {
	flag.Parse()
}
