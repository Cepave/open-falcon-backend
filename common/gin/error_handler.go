package gin

import (
	"fmt"
	"net/http"
	"gopkg.in/gin-gonic/gin.v1"
	json "github.com/bitly/go-simplejson"
)

type DataConflictError struct {
	ErrorCode int32
	ErrorMessage string
}
func (e DataConflictError) string() string {
	return fmt.Sprintf("[%d] %s", e.ErrorCode, e.ErrorMessage)
}
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

// Type of PanicProcessor, output 500 status with JSON message
func DefaultPanicProcessor(c *gin.Context, panicObject interface{}) {
	switch panicObject.(type) {
	case ValidationError:
		validateErrors := panicObject.(ValidationError)
		c.JSON(
			http.StatusBadRequest,
			map[string]interface{} {
				"http_status": http.StatusBadRequest,
				"error_code": -1,
				"error_message": validateErrors.Error(),
			},
		)
	case BindJsonError:
		jsonError := panicObject.(BindJsonError)
		c.JSON(
			http.StatusBadRequest,
			map[string]interface{} {
				"http_status": http.StatusBadRequest,
				"error_code": -101,
				"error_message": jsonError.Error(),
			},
		)
	default:
		c.JSON(
			http.StatusInternalServerError,
			map[string]interface{} {
				"http_status": http.StatusInternalServerError,
				"error_code": -1,
				"error_message": fmt.Sprintf("%v", panicObject),
			},
		)
	}
}

// Output http.StatusConflict with JSON body
func JsonConflictHandler(c *gin.Context, body interface{}) {
	c.JSON(
		http.StatusConflict,
		body,
	)
}

func JsonNoMethodHandler(c *gin.Context) {
	c.JSON(
		http.StatusNotFound,
		map[string]interface{} {
			"http_status": http.StatusMethodNotAllowed,
			"error_code": -1,
			"method": c.Request.Method,
			"uri": c.Request.RequestURI,
		},
	)
}

func JsonNoRouteHandler(c *gin.Context) {
	c.JSON(
		http.StatusNotFound,
		map[string]interface{} {
			"http_status": http.StatusNotFound,
			"error_code": -1,
			"uri": c.Request.RequestURI,
		},
	)
}
