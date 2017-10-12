package mvc

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"reflect"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"

	ogin "github.com/Cepave/open-falcon-backend/common/gin"
	"github.com/Cepave/open-falcon-backend/common/model"
	oreflect "github.com/Cepave/open-falcon-backend/common/reflect"
	ot "github.com/Cepave/open-falcon-backend/common/types"
	ov "github.com/Cepave/open-falcon-backend/common/validate"
)

// Defines configuration of MVC framework
type MvcConfig struct {
	// ConversionServer used to perform type conversion
	ConvertService ot.ConversionService
	// Validator object used to perform data validations
	Validator *validator.Validate
}

// Constructs default configuration of MVC framework
func NewDefaultMvcConfig() *MvcConfig {
	newValidator := validator.New()
	ov.RegisterDefaultValidators(newValidator)

	convService := ot.NewDefaultConversionService()
	ot.AddDefaultConverters(convService)

	return &MvcConfig{
		ConvertService: convService,
		Validator:      newValidator,
	}
}

// Builder object with internal objects
type MvcBuilder struct {
	config *MvcConfig
}

// Constructs a new builder with configuration
func NewMvcBuilder(newConfig *MvcConfig) *MvcBuilder {
	return &MvcBuilder{
		config: newConfig,
	}
}

// Builds gin.HandlerFunc by MVC handler
var _t_IoCloser = oreflect.TypeOfInterface((*io.Closer)(nil))

func (b *MvcBuilder) BuildHandler(handlerFunc MvcHandler) gin.HandlerFunc {
	funcValue := reflect.ValueOf(handlerFunc)
	funcType := funcValue.Type()

	if funcType.Kind() != reflect.Func {
		panic(fmt.Sprintf("Need to be function for \"MvcHandler\". Got: [%T]", handlerFunc))
	}

	/**
	 * Builds binding functions for output/input parameters
	 */
	inputTypes, outputTypes := oreflect.GetAllTypesForFunction(funcType)

	inputFunc := b.buildInputFunc(inputTypes)
	outputFunc := b.buildOutputFunc(outputTypes)
	// :~)

	return func(c *gin.Context) {
		inputParams := inputFunc(c)
		returnedValues := funcValue.Call(inputParams)

		if outputFunc != nil {
			outputFunc(c, returnedValues)
		}

		// Release closable resources binding
		releaseResources(c)
	}
}

type inputParamLoader func(c *gin.Context) interface{}

func (b *MvcBuilder) buildInputFunc(targetTypes []reflect.Type) func(c *gin.Context) []reflect.Value {
	/**
	 * Builds loaders for echo of the input parameters
	 */
	loaders := make([]inputParamLoader, len(targetTypes))
	for i, t := range targetTypes {
		loaders[i] = b.buildInputLoader(t)
	}
	// :~)

	return func(c *gin.Context) []reflect.Value {
		params := loadInputParams(c, loaders)
		valuesOfParams := make([]reflect.Value, len(params))

		for i, p := range params {
			valuesOfParams[i] = reflect.ValueOf(p)
		}

		return valuesOfParams
	}
}

var webObjectFuncs = map[string]inputParamLoader{
	"*gin.Context": func(c *gin.Context) interface{} {
		return c
	},
	"gin.Params": func(c *gin.Context) interface{} {
		return c.Params
	},
	"*http.Request": func(c *gin.Context) interface{} {
		return c.Request
	},
	"http.ResponseWriter": func(c *gin.Context) interface{} {
		return c.Writer
	},
	"gin.ResponseWriter": func(c *gin.Context) interface{} {
		return c.Writer
	},
	"*url.URL": func(c *gin.Context) interface{} {
		return c.Request.URL
	},
	"*multipart.Form": func(c *gin.Context) interface{} {
		return getMultipartForm(c)
	},
	"*multipart.Reader": func(c *gin.Context) interface{} {
		return getMultipartReader(c)
	},
	"http.Header": func(c *gin.Context) interface{} {
		return c.Writer.Header()
	},
}

const (
	_MultipartReader = "_mp_reader_"
	_MultipartForm   = "_mp_form_"
)

func getMultipartReader(c *gin.Context) *multipart.Reader {
	reader, ok := c.Get(_MultipartReader)
	if !ok {
		var err error
		reader, err = c.Request.MultipartReader()
		if err != nil {
			panic(fmt.Sprintf("Multpart has error: %v", err))
		}
		c.Set(_MultipartReader, reader)
	}

	return reader.(*multipart.Reader)
}
func getMultipartForm(c *gin.Context) *multipart.Form {
	form, ok := c.Get(_MultipartForm)
	if !ok {
		r := getMultipartReader(c)

		var err error
		form, err = r.ReadForm(2 * 1024 * 1024) // 2MB
		if err != nil {
			panic(fmt.Sprintf("Multpart(ReadForm) has error: %v", err))
		}

		c.Set(_MultipartForm, form)
	}

	return form.(*multipart.Form)
}

var _t_JsonUnmarshaler = oreflect.TypeOfInterface((*json.Unmarshaler)(nil))

func (b *MvcBuilder) buildInputLoader(targetType reflect.Type) inputParamLoader {
	typedFunc, ok := webObjectFuncs[targetType.String()]
	if ok {
		return typedFunc
	}

	switch targetType.String() {
	case "*validator.Validate":
		return b.getValidateFunc
	case "types.ConversionService":
		return b.getConversionServiceFunc
	}

	/**
	 * Builds the function for context binder
	 */
	if targetType.Implements(_t_ContextBinder) {
		return func(c *gin.Context) interface{} {
			value := reflect.New(targetType.Elem()).Interface()
			value.(ContextBinder).Bind(c)
			b.validateStruct(value)
			return value
		}
	}
	// :~)

	/**
	 * Binds the value by body of json
	 */
	if targetType.Implements(_t_JsonUnmarshaler) {
		return func(c *gin.Context) interface{} {
			value := reflect.New(targetType.Elem()).Interface()
			ogin.BindJson(c, value)
			b.validateStruct(value)
			return value
		}
	}
	// :~)

	/**
	 * Builds the struct value
	 */
	finalType := oreflect.FinalPointedType(targetType)
	if finalType.Kind() == reflect.Struct {
		pointerFunc := b.buildStructPointerFunc(finalType)
		return func(c *gin.Context) interface{} {
			structValue := pointerFunc(c)

			if targetType.Kind() == reflect.Ptr {
				return oreflect.NewFinalValueFrom(structValue, targetType).Interface()
			}

			return structValue.Elem().Interface()
		}
	}
	// :~)

	panic(fmt.Sprintf("Unknown type for input parameter: [%s]", targetType.String()))
}

func (b *MvcBuilder) buildStructPointerFunc(structType reflect.Type) func(c *gin.Context) reflect.Value {
	setters := b.buildWebParamFunc(structType)
	return func(c *gin.Context) reflect.Value {
		pointerValue := reflect.New(structType)
		structValue := pointerValue.Elem()

		for fieldName, paramFunc := range setters {
			structValue.FieldByName(fieldName).Set(
				reflect.ValueOf(paramFunc(c)),
			)
		}

		b.validateStruct(pointerValue.Interface())
		return pointerValue
	}
}

func (b *MvcBuilder) validateStruct(object interface{}) {
	typeOfValue := reflect.TypeOf(object)

	if typeOfValue.Kind() == reflect.Struct ||
		(typeOfValue.Kind() == reflect.Ptr &&
			typeOfValue.Elem().Kind() == reflect.Struct) {
		ogin.ConformAndValidateStruct(object, b.config.Validator)
	}
}

func (b *MvcBuilder) buildWebParamFunc(structType reflect.Type) map[string]inputParamLoader {
	result := make(map[string]inputParamLoader)

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		fieldLoader := buildParamLoader(field, b.config.ConvertService)
		if fieldLoader == nil {
			/**
			 * Supporting of nested struct
			 *
			 * Could be pointer to struct or concrete struct
			 */
			finalType := oreflect.FinalPointedType(field.Type)
			if finalType.Kind() == reflect.Struct {
				structFunc := b.buildStructPointerFunc(finalType)
				if field.Type.Kind() == reflect.Ptr {
					result[field.Name] = func(c *gin.Context) interface{} {
						return structFunc(c).Interface()
					}
				} else {
					result[field.Name] = func(c *gin.Context) interface{} {
						return structFunc(c).Elem().Interface()
					}
				}
			}
			// :~)

			continue
		}

		result[field.Name] = fieldLoader
	}

	return result
}

func (b *MvcBuilder) getValidateFunc(c *gin.Context) interface{} {
	return b.config.Validator
}
func (b *MvcBuilder) getConversionServiceFunc(c *gin.Context) interface{} {
	return b.config.ConvertService
}

var _t_OutputBody = oreflect.TypeOfInterface((*OutputBody)(nil))
var _t_JsonMarshaler = oreflect.TypeOfInterface((*json.Marshaler)(nil))
var _t_Stringer = oreflect.TypeOfInterface((*fmt.Stringer)(nil))

func (b *MvcBuilder) buildOutputFunc(targetTypes []reflect.Type) func(c *gin.Context, returnedValue []reflect.Value) {
	if len(targetTypes) == 0 {
		return nil
	}

	/**
	 * Builds the processor for every return value
	 */
	valueProcessors := make([]func(c *gin.Context, returnedValue reflect.Value), len(targetTypes))
	for i, targetType := range targetTypes {
		var processor func(c *gin.Context, returnedValue reflect.Value) = nil

		switch {
		case targetType == _t_Paging:
			processor = func(c *gin.Context, returnValue reflect.Value) {
				ogin.HeaderWithPaging(c, returnValue.Interface().(*model.Paging))
			}
		case targetType.Implements(_t_OutputBody):
			processor = func(c *gin.Context, returnValue reflect.Value) {
				returnValue.Interface().(OutputBody).Output(c)
			}
		case targetType.Implements(_t_JsonMarshaler):
			processor = func(c *gin.Context, returnValue reflect.Value) {
				JsonOutputBody(returnValue.Interface()).Output(c)
			}
		case targetType.Kind() == reflect.String ||
			targetType.Implements(_t_Stringer):
			processor = func(c *gin.Context, returnValue reflect.Value) {
				TextOutputBody(returnValue.Interface()).Output(c)
			}
		default:
			panic(fmt.Sprintf("Unknown type for building output: [%s]", targetType))
		}

		valueProcessors[i] = processor
	}
	// :~)

	return func(c *gin.Context, returnValues []reflect.Value) {
		for i, processor := range valueProcessors {
			if !returnValues[i].IsValid() {
				continue
			}

			processor(c, returnValues[i])
		}
	}
}

func loadInputParams(c *gin.Context, loaders []inputParamLoader) []interface{} {
	result := make([]interface{}, len(loaders))

	for i, loader := range loaders {
		result[i] = loader(c)
	}

	return result
}
func releaseResources(c *gin.Context) {
	releaseMultipartFiles(c)
}
