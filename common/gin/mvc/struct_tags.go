package mvc

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"mime/multipart"
	"strconv"

	ogin "github.com/Cepave/open-falcon-backend/common/gin"
	"github.com/Cepave/open-falcon-backend/common/model"
	ot "github.com/Cepave/open-falcon-backend/common/types"
	"github.com/gin-gonic/gin"
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
	"pageSize": true,
	"pageOrderBy": true,
	"default": true,
}

var _t_Paging = reflect.TypeOf(&model.Paging{})
func buildParamLoader(field reflect.StructField, convSrv ot.ConversionService) inputParamLoader {
	/**
	 * Process paging object
	 */
	if field.Type == _t_Paging {
		defaultPaging := loadDefaultPaging(field.Tag)
		return func(c *gin.Context) interface{} {
			return ogin.PagingByHeader(c, defaultPaging)
		}
	}
	// :~)

	tagContext := loadTag(field)
	if tagContext.getterType == 0 {
		return nil
	}

	return tagContext.getLoader(field.Type, convSrv)
}

const (
	paramGetterType = 1
	paramCheckerType = 5

	keyGetterType = 2
	fileGetterType = 3
	pagingGetterType = 4
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
	case paramCheckerType:
		checker := paramCheckers[t.getterName]
		return func(c *gin.Context) interface{} {
			return checker(c, t.paramName)
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
	tagContext := &tagContext{}

	iterateMvcTagProperties(
		field.Tag,
		func(propName string, propValue string) {
			switch propName {
			case "query", "cookie", "param", "form", "header", "req", "basicAuth":
				if strings.HasPrefix(propValue, "?") {
					tagContext.getterType = paramCheckerType
					tagContext.getterName = propName
					tagContext.paramName = strings.TrimLeft(propValue, "?")
				} else {
					tagContext.getterType = paramGetterType
					tagContext.getterName = propName
					tagContext.paramName = propValue
				}
			case "file", "fileHeader":
				tagContext.getterType = fileGetterType
				tagContext.getterName = propName
				tagContext.paramName = propValue
			case "key":
				if strings.HasPrefix(propValue, "?") {
					tagContext.getterType = paramCheckerType
					tagContext.getterName = propName
					tagContext.paramName = strings.TrimLeft(propValue, "?")
				} else {
					tagContext.getterType = keyGetterType
					tagContext.getterName = propName
					tagContext.paramName = propValue
				}
			case "default":
				tagContext.defaultValue = propValue
			default:
				panic(fmt.Sprintf("Cannot recognize property name: [%s]", propName))
			}
		},
	)

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

func loadDefaultPaging(tag reflect.StructTag) *model.Paging {
	paging := &model.Paging {
		Size: 64,
		Position: 1,
	}

	iterateMvcTagProperties(
		tag,
		func(propName string, propValue string) {
			switch propName {
			case "pageSize":
				defaultSize, err := strconv.ParseInt(propValue, 10, 32)
				if err != nil {
					panic(fmt.Sprintf("Cannot parse pageSize[%s]. Error: %v.", propValue, err))
				}
				paging.Size = int32(defaultSize)
			case "pageOrderBy":
				parsedOrderBy, err := ogin.ParseOrderBy(propValue)
				if err != nil {
					panic(fmt.Sprintf("Cannot parse pageOrderBy[%s]. Error: %v.", propValue, err))
				}
				paging.OrderBy = parsedOrderBy
			default:
				panic(fmt.Sprintf("Cannot recognize property name: [%s]", propName))
			}
		},
	)

	return paging
}

func iterateMvcTagProperties(tag reflect.StructTag, propProcessor func(propName string, propValue string)) {
	tagValue := tag.Get(mvcTag)
	if tagValue == "" {
		return
	}

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

		propProcessor(propName, propValue)
	}
}
