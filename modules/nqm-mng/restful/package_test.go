package restful

import (
	"testing"

	testingHttp "github.com/Cepave/open-falcon-backend/common/testing/http"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

var httpClientConfig = testingHttp.NewHttpClientConfigByFlag()
