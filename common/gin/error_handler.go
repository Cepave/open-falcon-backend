package gin

import (
	"fmt"
	or "github.com/Cepave/open-falcon-backend/common/runtime"
	json "github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Defines the error used to represent the "409 Conflict" error
//
// Usage
//
// Case 1: The conflict of unique key on a new or modified data
//
// Case 2: The conflict of complex logical of business on modified data
//
// HTTP Specification
//
// Following paragraph comes from RFC-2616(HTTP/1.1)
//
// 	The request could not be completed due to a conflict with the current state of the resource.
// 	This code is only allowed in situations where it is expected that
// 	the user might be able to resolve the conflict and resubmit the request.
//
// 	The response body SHOULD include enough information for the user to
// 	recognize the source of the conflict.
// 	Ideally, the response entity would include enough information for the user or
// 	user agent to fix the problem; however,
// 	that might not be possible and is not required.
//
// 	Conflicts are most likely to occur in response to a PUT request.
// 	For example, if versioning were being used and the entity being PUT included changes to
// 	a resource which conflict with those made by an earlier (third-party) request,
// 	the server might use the 409 response to indicate that it can't complete the request.
//
// 	In this case, the response entity would likely contain a list of
// 	the differences between the two versions in a format defined by the response Content-Type.
//
// See: https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html
type DataConflictError struct {
	ErrorCode    int32
	ErrorMessage string
}

// Implements the error interface
func (e DataConflictError) Error() string {
	return fmt.Sprintf("[%d] %s", e.ErrorCode, e.ErrorMessage)
}

// Marshal this type of error to:
//
// 	{
// 		"http_status": 409,
// 		"error_code": e.ErrorCode,
// 		"error_message": e.ErrorMessage,
// 	}
func (e DataConflictError) MarshalJSON() ([]byte, error) {
	jsonObject := json.New()

	jsonObject.Set("http_status", http.StatusConflict)
	jsonObject.Set("error_code", e.ErrorCode)
	jsonObject.Set("error_message", e.ErrorMessage)

	return jsonObject.MarshalJSON()
}

// This callback function is used to process panic object
type PanicProcessor func(c *gin.Context, panic interface{})

// Builds a gin.HandlerFunc, which is used to handle not-nil object of panic
func BuildJsonPanicProcessor(panicProcessor PanicProcessor) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			p := recover()
			if p == nil {
				return
			}

			panicProcessor(c, p)
		}()

		c.Next()
	}
}

// Process various of panic object with corresponding HTTP status
//
// ValidationError - Use http.StatusBadRequest as output status
//
// BindJsonError - Use http.StatusBadRequest as output status
//
// Otherwise, use http.StatusInternalServerError as output status
func DefaultPanicProcessor(c *gin.Context, panicObject interface{}) {
	stack := or.GetCallerInfoStack(2, 16).AsStringStack()

	logger.Warnf("[GIN] Got panic: %v", panicObject)
	for _, messageInStack := range stack {
		logger.Warnf("%s", messageInStack)
	}

	switch errObject := panicObject.(type) {
	case ValidationError:
		c.JSON(
			http.StatusBadRequest,
			map[string]interface{}{
				"http_status":   http.StatusBadRequest,
				"error_code":    -1,
				"error_message": errObject.Error(),
			},
		)
	case BindJsonError:
		c.JSON(
			http.StatusBadRequest,
			map[string]interface{}{
				"http_status":   http.StatusBadRequest,
				"error_code":    -101,
				"error_message": errObject.Error(),
			},
		)
	case DataConflictError:
		c.JSON(
			http.StatusConflict,
			map[string]interface{}{
				"http_status":   http.StatusConflict,
				"error_code":    errObject.ErrorCode,
				"error_message": errObject.ErrorMessage,
			},
		)
	default:
		c.JSON(
			http.StatusInternalServerError,
			map[string]interface{}{
				"http_status":   http.StatusInternalServerError,
				"error_code":    -1,
				"error_message": fmt.Sprintf("%v", panicObject),
				"error_stack":   stack,
			},
		)
	}
}

// Output http.StatusConflict as JSON.
func JsonConflictHandler(c *gin.Context, body interface{}) {
	c.JSON(
		http.StatusConflict,
		body,
	)
}

// Output http.StatusMethodNotAllowed as JSON;
//
// 	{
// 		"http_status": 405,
// 		"error_code": -1,
// 		"method": c.Request.Method,
// 		"uri": c.Request.RequestURI,
// 	}
func JsonNoMethodHandler(c *gin.Context) {
	c.JSON(
		http.StatusNotFound,
		map[string]interface{}{
			"http_status": http.StatusMethodNotAllowed,
			"error_code":  -1,
			"method":      c.Request.Method,
			"uri":         c.Request.RequestURI,
		},
	)
}

// Output http.StatusNotFound as JSON;
//
// 	{
// 		"http_status": 404,
// 		"error_code": -1,
// 		"uri": c.Request.RequestURI,
// 	}
func JsonNoRouteHandler(c *gin.Context) {
	c.JSON(
		http.StatusNotFound,
		map[string]interface{}{
			"http_status": http.StatusNotFound,
			"error_code":  -1,
			"uri":         c.Request.RequestURI,
		},
	)
}
