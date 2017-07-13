// A MVC binder for free-style of function handler with *gin.Context
//
// Abstract
//
// There are may tedious processes for coding on web service:
//
// 	1. Type conversion from HTTP query parameter to desired type in GoLang.
// 	2. Binding body of HTTP POST to JSON object in GoLang.
// 		2.1. Perform post-process(e.x. trim text) of binding data
// 		2.2. Perform data validation of binding data
// 	3. Convert the result data to JSON response.
//
// Gin has provided foundation features on simple and versatile web application,
// this framework try to enhance the building of web application on instinct way.
//
// MVC Handler
//
// A MVC handler can be any function with ruled types of parameter and defined returned types.
//
// 	type MvcHandler interface{}
//
// You can define handler of supported:
//
// 	func(req *http.Request, params *gin.Params) OutputBody {
// 		/* ... */
// 		return TextOutputBody("Hello World")
// 	}
//
// 	func(
// 		data *struct {
// 			Name string `mvc:"query[name]"`
// 			Age int `mvc:"query[age]"`
// 			SessionId string `mvc:"header[session_id"`
//		}
// 	) OutputBody {
// 		/* ... */
// 		return TextOutputBody("Hello World")
// 	}
//
// Build Gin HandlerFunc
//
// After you define the MVC handler, you could use "MvcBuilder.BuildHandler()" to
// convert your handler to "gin.HandlerFunc".
//
// 	mvcBuilder := NewMvcBuilder(NewDefaultMvcConfig())
// 	engine.Get("/your-web-service", mvcBuilder.BuildHandler(your_mvc_handler))
//
// Parameters of Handler
//
// Supported types:
//
// "ContextBinder" - Feeds the context to implementation of Bind(*gin.Context) function.
//  This type of value woule be checked by ogin.ConformAndValidateStruct automatically.
//
// 'json.Unmarshaler' - If the type of value is json.Unmarshaler, use the UnmarshalJSON([]byte) function of the value
//  This type of value woule be checked by ogin.ConformAndValidateStruct automatically.
//
// "<struct>" - See parameter tags for automatic binding
//  This type of value woule be checked by ogin.ConformAndValidateStruct automatically.
//
// "*gin.Context" - The context object of current request
//
// "gin.ResponseWriter" - See gin.ResponseWriter
//
// "gin.Params" - See gin.Params
//
// "*http.Request" - See http.Request
//
// "http.ResponseWriter" - See http.ResponseWriter
//
// "*url.URL" - See url.URL
//
// "*multipart.Reader" - See multipart.Reader; Once you use *multipart.Form, the reader would reach EOF.
//
// "*multipart.Form" - See multipart.Form
//
// "*validator.Validate" - See go-playground/validator.v9
//
// "types.ConversionService" - See ConversionService
//
// Return value of Handler
//
// "OutputBody" is the main definition for output of web service, it has build-in functions for certain types of output:
//
//     "JsonOutputBody()" - Uses gin.Context.JSON function to perform output
//     "TextOutputBody()" - Uses gin.Context.String function to generate output body(by fmt.Sprintf("%v"))
//     "HtmlOutputBody()," XmlOutputBody(), YamlOutpuBody() - Calls function of gin.Context, respectively.
//
// "json.Marshaler" - If the type of returned value is json.Marshaler, use JsonOutputBody() as output type
//
// "string" - If the type of returned value is string, use TextOutputBody() as output type
//
// "fmt.Stringer" - As same as string
//
// "*model.Paging" - Output the content of paging to HTTP header
//
// Tagging Struct
//
// There are various definition of tags could be used on struct:
//
// 	type MyData struct {
// 		Name string `mvc:"query[name] default[NO-NAME]"`
// 		H1 string `mvc:"header[h1]"`
// 		FormV1 int `mvc:"form[v1]"`
// 		FormV3 []string `mvc:"form[v2]"` // As slice
// 		File1 multipart.File `mvc:"file[f1]"`
// 		File2 multipart.File `mvc:"file[f2]"`
// 	}
//
// Default Value
//
// 	mvc:"query[param_name_1] default[20]" - Gives value 20 if the value of binding is empty
// 	mvc:"query[param_name_1] default[20,40,30]" - Gives value [20, 40, 30](as array, no space)if the value of binding is empty
//
// Parameters, Heaer, Cookie, and From
//
//  mvc:"query[param_name_1]" - Use query parameter param_name_1 as binding value
//  mvc:"query[?param_name_1]" - Must be bool type, used to indicate whether or not has viable value for this parameter
//  mvc:"cookie[ck_1]" - Use the value of cookie ck_1 as binding value
//  mvc:"cookie[?ck_2]" - Must be bool type, used to indicate whether or not has viable value for this parameter
//  mvc:"param[pm_1]" - Use the value of URI parameter pm_1 as binding value
//  mvc:"form[in_1]" - Use the form value of in_1 as binding value
//  mvc:"form[?in_2]" - Must be bool type, used to indicate whether or not has viable value for this parameter
//  mvc:"header[Content-Type]" - Use the header value of Content-Type as binding value
//  mvc:"header[?pg_id]" - Must be bool type, used to indicate whether or not has viable value for this parameter
//  mvc:"key[key-1]" - Use the key value of key-1 as binding value
//  mvc:"key[?key-3]" - Must be bool type, used to indicate whether or not has viable value for this parameter
//
// By default, if the value of binding is existing, the framework would use the default value of binding type.
//
// HTTP
//  mvc:"req[ClientIp]" - The IP of client, the type of value could be string or ​net.IP
//  mvc:"req[ContentType]" - The content type of request, must be string
//  mvc:"req[Referer]" - The "Referer" of request, must be string
//  mvc:"req[UserAgent]" - The "User-Agent" of request, must be string
//  mvc:"req[Method]" - The method of request, must be string
//  mvc:"req[Url]" - The url of reqeust, must be string or ​url.URL
//  mvc:"req[Proto]" - The protocol version for incoming server requests, must be string
//  mvc:"req[ProtoMajor]" - The protocol version for incoming server requests, must be int
//  mvc:"req[ProtoMinor]" - The protocol version for incoming server requests, must be int
//  mvc:"req[ContentLength]" - The ContentLength? records the length of the associated content, must be int64
//  mvc:"req[Host]" - For server requests Host specifies the host on which the URL is sought, must be string
//  mvc:"req[RemoteAddr]" - RemoteAddr allows HTTP servers and other software to record the network address that
//  	sent the request, usually for logging, must be string
//  mvc:"req[RequestURI]" - RequestURI is the unmodified Request-URI of the Request-Line (RFC 2616, Section 5.1) as
//  	sent by the client to a server, must be string
// PAGING
//
// Must be type of "*model.Paging"
//
// 	mvc:"pageSize[50]" - The default value of page size is 50
// 	mvc:"pageOrderBy[name:age]" - The default value of "orderBy" property of paging object is 'name:age'
//
// SECURITY
//
//  mvc:"basicAuth[username]" - The username of BasicAuth?, See RFC-2617
//  mvc:"basicAuth[password]" - The password of BasicAuth?, See RFC-2617
//
// FILE UPLOAD
//
//  mvc:"file[f1]" - The file of request by key value, must be "multipart.File"(or "[]multipart.File")
//      You don't have to close this resource, Gin MVC would do your favour.
//  mvc:"fileHeader[f1]" - The file header of request by key value, must be "*multipart.FileHeader"(or "[]*multipart.FileHeader")
//
// Data Conversion
//
// This framework depends on ConversionService to perform type conversion
//
// Data Validation
//
// For struct type of input parameter, this framework would use "MvcConfig.Validator" and "common/conform"
// to perform post-process on the parameter and validate it.
//
// See go-playground/validator: https://godoc.org/gopkg.in/go-playground/validator.v9
//
// See leebenson/conform: https://github.com/leebenson/conform
package mvc

import (
	"net/http"
	"reflect"

	"github.com/Cepave/open-falcon-backend/common/utils"
	"github.com/gin-gonic/gin"

	oreflect "github.com/Cepave/open-falcon-backend/common/reflect"
)

// Could be any type of function
type MvcHandler interface{}

// Main interface for binding context
type ContextBinder interface {
	Bind(*gin.Context)
}

// Function type of "ContextBuilder"
type ContextBinderFunc func(*gin.Context)

func (f ContextBinderFunc) Bind(context *gin.Context) {
	f(context)
}

var _t_ContextBinder = oreflect.TypeOfInterface((*ContextBinder)(nil))

// Main interface for generating response
type OutputBody interface {
	Output(*gin.Context)
}

// Function version of "OutputBody"
type OutputBodyFunc func(*gin.Context)

// As implementation of "OutputBody"
func (f OutputBodyFunc) Output(context *gin.Context) {
	f(context)
}

// Uses "(*gin.Context).JSON(http.StatusOK, v)" to perform response
func JsonOutputBody(v interface{}) OutputBody {
	return JsonOutputBody2(http.StatusOK, v)
}

// Uses "(*gin.Context).JSON(code, v)" to perform response
func JsonOutputBody2(code int, v interface{}) OutputBody {
	return OutputBodyFunc(func(context *gin.Context) {
		context.JSON(code, v)
	})
}

// Output the value or not found error if the value is not viable
func JsonOutputOrNotFound(v interface{}) OutputBody {
	if !isViableValue(v) {
		return NotFoundOutputBody
	}

	return JsonOutputBody(v)
}

// Uses "(*gin.Context).String(http.StatusOK, v)" to perform response
func TextOutputBody(v interface{}) OutputBody {
	return TextOutputBody2(http.StatusOK, v)
}

// Uses "(*gin.Context).String(code, v)" to perform response
func TextOutputBody2(code int, v interface{}) OutputBody {
	return OutputBodyFunc(func(context *gin.Context) {
		context.String(code, "%s", v)
	})
}

// Output the value or not found error if the value is not viable
func TextOutputOrNotFound(v interface{}) OutputBody {
	if !isViableValue(v) {
		return NotFoundOutputBody
	}

	return TextOutputBody(v)
}

// Uses "(*gin.Context).HTML(http.StatusOK, name, v)" to perform response
func HtmlOutputBody(name string, v interface{}) OutputBody {
	return HtmlOutputBody2(http.StatusOK, name, v)
}

// Uses "(*gin.Context).HTML(code, name, v)" to perform response
func HtmlOutputBody2(code int, name string, v interface{}) OutputBody {
	return OutputBodyFunc(func(context *gin.Context) {
		context.HTML(code, name, v)
	})
}

// Output the value or not found error if the value is not viable
func HtmlOutputOrNotFound(name string, v interface{}) OutputBody {
	if !isViableValue(v) {
		return NotFoundOutputBody
	}

	return HtmlOutputBody(name, v)
}

// Uses "(*gin.Context).XML(http.StatusOK, v)" to perform response
func XmlOutputBody(v interface{}) OutputBody {
	return XmlOutputBody2(http.StatusOK, v)
}

// Uses "(*gin.Context).XML(code, v)" to perform response
func XmlOutputBody2(code int, v interface{}) OutputBody {
	return OutputBodyFunc(func(context *gin.Context) {
		context.XML(code, v)
	})
}

// Output the value or not found error if the value is not viable
func XmlOutputOrNotFound(v interface{}) OutputBody {
	if !isViableValue(v) {
		return NotFoundOutputBody
	}

	return XmlOutputBody(v)
}

// Uses "(*gin.Context).YAML(http.StatusOK, v)" to perform response
func YamlOutputBody(v interface{}) OutputBody {
	return YamlOutputBody2(http.StatusOK, v)
}

// Uses "(*gin.Context).YAML(code, v)" to perform response
func YamlOutputBody2(code int, v interface{}) OutputBody {
	return OutputBodyFunc(func(context *gin.Context) {
		context.YAML(code, v)
	})
}

// Output the value or not found error if the value is not viable
func YamlOutputOrNotFound(v interface{}) OutputBody {
	if !isViableValue(v) {
		return NotFoundOutputBody
	}

	return YamlOutputBody(v)
}

func isViableValue(v interface{}) bool {
	return utils.ValueExt(reflect.ValueOf(v)).IsViable()
}
