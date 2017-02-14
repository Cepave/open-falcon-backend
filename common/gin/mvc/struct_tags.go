package mvc

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"mime/multipart"

	ot "github.com/Cepave/open-falcon-backend/common/types"
	"gopkg.in/gin-gonic/gin.v1"
)

const mvcTag = "mvc"

var defProp = map[string]bool {
	"file": true,
	"fileHeader": true,
	"query": true,
	"cookie": true,
	"param": true,
	"form": true,
	"header": true,
	"key": true,
	"req": true,
	"basicAuth": true,
	"default": true,
}

func buildParamLoader(field reflect.StructField, convSrv ot.ConversionService) inputParamLoader {
	tagContext := loadTag(field)
	if tagContext == nil {
		return nil
	}

	return tagContext.getLoader(field.Type, convSrv)
}

const (
	paramGetterType = 1
	keyGetterType = 2
	fileGetterType = 3
)

type tagContext struct {
	getterType int
	getterName string
	paramName string
	defaultValue string
}
func (t *tagContext) getDefaultValueAsSlice() []string {
	return strings.Split(t.defaultValue, ",")
}
func (t *tagContext) getLoader(targetType reflect.Type, convSrv ot.ConversionService) inputParamLoader {
	switch t.getterType {
	case paramGetterType:
		paramGetter := paramGetters[t.getterName]

		switch targetType.Kind() {
		case reflect.Array, reflect.Slice:
			return func(c *gin.Context) interface{} {
				return convSrv.ConvertTo(
					paramGetter.getParamAsArray(c, t.paramName, t.getDefaultValueAsSlice()),
					targetType,
				)
			}
		default:
			return func(c *gin.Context) interface{} {
				return convSrv.ConvertTo(
					paramGetter.getParam(c, t.paramName, t.defaultValue),
					targetType,
				)
			}
		}
	case keyGetterType:
		return func(c *gin.Context) interface{} {
			var finalDefaultValue interface{} = t.defaultValue
			switch targetType.Kind() {
			case reflect.Array, reflect.Slice:
				finalDefaultValue = t.getDefaultValueAsSlice()
			}

			return convSrv.ConvertTo(
				keyGetter.getValue(c, t.paramName, finalDefaultValue),
				targetType,
			)
		}
	case fileGetterType:
		switch t.getterName {
		case "file":
			return func(c *gin.Context) interface{} {
				headers := getFileHeaders(c, t.paramName)
				if len(headers) == 0 {
					return []multipart.File(nil)
				}

				switch targetType.Kind() {
				case reflect.Array, reflect.Slice:
					files := make([]multipart.File, len(headers))
					for i, header := range headers {
						files[i] = openMultipartFile(c, header)
					}
					return files
				default:
					return openMultipartFile(c, headers[0])
				}
			}
		case "fileHeader":
			return func(c *gin.Context) interface{} {
				headers := getFileHeaders(c, t.paramName)
				if len(headers) == 0 {
					return []multipart.File(nil)
				}

				switch targetType.Kind() {
				case reflect.Array, reflect.Slice:
					return headers;
				default:
					return headers[0]
				}
			}
		}
	}

	panic(fmt.Sprintf("Unknown type of getter: [%d]", t.getterType))
}

var propRegExp, _ = regexp.Compile(`^(\w+)\[([^]]+)\]$`)
func loadTag(field reflect.StructField) *tagContext {
	tagValue := field.Tag.Get(mvcTag)
	if tagValue == "" {
		return nil
	}

	tagContext := &tagContext{}

	for _, propPair := range strings.Split(tagValue, " ") {
		propPair = strings.TrimSpace(propPair)

		matches := propRegExp.FindStringSubmatch(propPair)
		if matches == nil {
			panic(fmt.Sprintf("Cannot recognize in mvc:\"prop...\": %s", propPair))
		}

		propName := matches[1]
		propValue := matches[2]

		if _, ok := defProp[propName]; !ok {
			panic(fmt.Sprintf("Cannot recognize property name: [%s]", propName))
		}

		switch propName {
		case "query", "cookie", "param", "form", "header", "req", "basicAuth":
			tagContext.getterType = paramGetterType
			tagContext.getterName = propName
			tagContext.paramName = propValue
		case "file", "fileHeader":
			tagContext.getterType = fileGetterType
			tagContext.getterName = propName
			tagContext.paramName = propValue
		case "key":
			tagContext.getterType = keyGetterType
			tagContext.getterName = propName
			tagContext.paramName = propValue
		case "default":
			tagContext.defaultValue = propValue
		default:
			panic(fmt.Sprintf("Unknown context of mvc tag: %s", propPair))
		}
	}

	return tagContext
}

func getFileHeaders(context *gin.Context, fieldName string) []*multipart.FileHeader {
	fileHeaders, ok := getMultipartForm(context).File[fieldName]
	if !ok {
		return nil
	}

	return fileHeaders
}

// This function would put opened file into context's value
const _keyOpenedFiles = "_opened_files_"
func openMultipartFile(context *gin.Context, fileHeader *multipart.FileHeader) multipart.File {
	file, err := fileHeader.Open()
	if err != nil {
		panic(fmt.Sprintf("Cannot open file[%s]. Error: %v", fileHeader.Filename, err))
	}

	openedFiles, ok := context.Get(_keyOpenedFiles)
	if !ok {
		openedFiles = make([]multipart.File, 0)
		context.Set(_keyOpenedFiles, openedFiles)
	}

	openedFiles = append(openedFiles.([]multipart.File), file)
	context.Set(_keyOpenedFiles, openedFiles)

	return file
}
func releaseMultipartFiles(context *gin.Context) {
	openedFiles, ok := context.Get(_keyOpenedFiles)
	if !ok {
		return
	}

	for _, file := range openedFiles.([]multipart.File) {
		if err := file.Close(); err != nil {
			logger.Errorf("Close file has error: [%v]", err)
		}
	}
}
