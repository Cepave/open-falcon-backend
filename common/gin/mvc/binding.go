package mvc

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"

	oreflect "github.com/Cepave/open-falcon-backend/common/reflect"
)

// Alias name of http status
type HttpStatus int

// Main interface for binding context
type ContextBinder interface {
	Bind(*gin.Context)
}
type ContextBinderFunc func(*gin.Context)
func (f ContextBinderFunc) Bind(context *gin.Context) {
	f(context)
}

var _t_ContextBinder = oreflect.TypeOfInterface((*ContextBinder)(nil))

// Main interface for generating response
type OutputBody interface {
	Output(*gin.Context)
}

// Function versino of OutputBody
type OutputBodyFunc func(*gin.Context)
func (f OutputBodyFunc) Output(context *gin.Context) {
	f(context)
}

func JsonOutputBody(v interface{}) OutputBody {
	return JsonOutputBody2(http.StatusOK, v)
}
func JsonOutputBody2(code int, v interface{}) OutputBody {
	return OutputBodyFunc(func(context *gin.Context) {
		context.JSON(code, v)
	})
}

func TextOutputBody(v interface{}) OutputBody {
	return TextOutputBody2(http.StatusOK, v)
}
func TextOutputBody2(code int, v interface{}) OutputBody {
	return OutputBodyFunc(func(context *gin.Context) {
		context.String(code, "%s", v)
	})
}

func HtmlOutputBody(name string, v interface{}) OutputBody {
	return HtmlOutputBody2(http.StatusOK, name, v)
}
func HtmlOutputBody2(code int, name string, v interface{}) OutputBody {
	return OutputBodyFunc(func(context *gin.Context) {
		context.HTML(code, name, v)
	})
}

func XmlOutputBody(v interface{}) OutputBody {
	return XmlOutputBody2(http.StatusOK, v)
}
func XmlOutputBody2(code int, v interface{}) OutputBody {
	return OutputBodyFunc(func(context *gin.Context) {
		context.XML(code, v)
	})
}

func YamlOutputBody(v interface{}) OutputBody {
	return YamlOutputBody2(http.StatusOK, v)
}
func YamlOutputBody2(code int, v interface{}) OutputBody {
	return OutputBodyFunc(func(context *gin.Context) {
		context.YAML(code, v)
	})
}

// Could be any type of function
type MvcHandler interface{}
