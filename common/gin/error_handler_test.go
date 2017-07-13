package gin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"

	. "gopkg.in/check.v1"
)

type TestErrorHandlerSuite struct{}

var _ = Suite(&TestErrorHandlerSuite{})

func ExampleDataConflictError() {
	engine := NewDefaultJsonEngine(&GinConfig{Mode: gin.ReleaseMode})
	engine.GET(
		"/data-conflict",
		func(c *gin.Context) {
			panic(DataConflictError{33, "Sample Conflict"})
		},
	)

	req := httptest.NewRequest(http.MethodGet, "/data-conflict", nil)
	resp := httptest.NewRecorder()

	engine.ServeHTTP(resp, req)

	fmt.Println(resp.Body)

	// Output:
	// {"error_code":33,"error_message":"Sample Conflict","http_status":409}
}
