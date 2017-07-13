package mvc

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var paramGetters = map[string]paramGetter{
	"query":     queryGetterImpl(true),
	"cookie":    cookieGetterImpl(true),
	"param":     uriParamGetterImpl(true),
	"header":    headerGetterImpl(true),
	"form":      formGetterImpl(true),
	"req":       reqGetterImpl(true),
	"basicAuth": basicAuthImpl(true),
}

var paramCheckers = map[string]boolParamChecker{
	"query": func(context *gin.Context, paramName string) bool {
		return isStringViable(context.Query(paramName))
	},
	"form": func(context *gin.Context, paramName string) bool {
		return isStringViable(context.PostForm(paramName))
	},
	"cookie": func(context *gin.Context, paramName string) bool {
		cookie, _ := context.Cookie(paramName)
		return isStringViable(cookie)
	},
	"header": func(context *gin.Context, paramName string) bool {
		return isStringViable(context.Request.Header.Get(paramName))
	},
	"key": func(context *gin.Context, paramName string) bool {
		v, ok := context.Get(paramName)

		if !ok {
			return false
		}

		s, isString := v.(string)
		if isString {
			return isStringViable(s)
		}

		return v != nil
	},
}

func isStringViable(v string) bool {
	return strings.TrimSpace(v) != ""
}

var keyGetter = keyGetterImpl(true)

type basicAuthImpl bool

func (b basicAuthImpl) getParam(context *gin.Context, paramName string, defaultValue string) string {
	username, password, ok := context.Request.BasicAuth()

	if !ok {
		panic("No basic authentication")
	}

	switch paramName {
	case "username":
		return username
	case "password":
		return password
	}

	panic(fmt.Sprintf("Unknown param name for BasicAuth: [%s]", paramName))
}
func (b basicAuthImpl) getParamAsArray(context *gin.Context, paramName string, defaultValue []string) []string {
	return []string{b.getParam(context, paramName, "")}
}

type keyGetterImpl bool

func (k keyGetterImpl) getValue(context *gin.Context, key string, notExistsValue interface{}) interface{} {
	v, ok := context.Get(key)
	if !ok {
		v = notExistsValue
	}

	return v
}

type paramGetter interface {
	getParam(context *gin.Context, paramName string, defaultValue string) string
	getParamAsArray(context *gin.Context, paramName string, defaultValue []string) []string
}

type formGetterImpl bool

func (f formGetterImpl) getParam(context *gin.Context, paramName string, defaultValue string) string {
	v, _ := context.GetPostForm(paramName)
	if v == "" {
		v = defaultValue
	}

	return v
}
func (f formGetterImpl) getParamAsArray(context *gin.Context, paramName string, defaultValue []string) []string {
	req := context.Request

	req.ParseForm()
	req.ParseMultipartForm(32 << 20) // 32 MB

	if values := req.PostForm[paramName]; len(values) > 0 {
		return values
	}

	if req.MultipartForm != nil && req.MultipartForm.File != nil {
		if values := req.MultipartForm.Value[paramName]; len(values) > 0 {
			return values
		}
	}

	return defaultValue
}

type headerGetterImpl bool

func (h headerGetterImpl) getParam(context *gin.Context, paramName string, defaultValue string) string {
	v := context.Request.Header.Get(paramName)
	if v == "" {
		v = defaultValue
	}

	return v
}
func (h headerGetterImpl) getParamAsArray(context *gin.Context, paramName string, defaultValue []string) []string {
	headerKey := http.CanonicalHeaderKey(paramName)

	headerValue, ok := context.Request.Header[headerKey]
	if !ok {
		return defaultValue
	}

	return headerValue
}

type queryGetterImpl bool

func (g queryGetterImpl) getParam(context *gin.Context, paramName string, defaultValue string) string {
	v := context.Query(paramName)
	if v == "" {
		v = defaultValue
	}

	return v
}
func (g queryGetterImpl) getParamAsArray(context *gin.Context, paramName string, defaultValue []string) []string {
	req := context.Request
	if values, ok := req.URL.Query()[paramName]; ok && len(values) > 0 {
		return values
	}

	return defaultValue
}

type cookieGetterImpl bool

func (c cookieGetterImpl) getParam(context *gin.Context, paramName string, defaultValue string) string {
	v, _ := context.Cookie(paramName)

	if v == "" {
		v = defaultValue
	}

	return v
}
func (c cookieGetterImpl) getParamAsArray(context *gin.Context, paramName string, defaultValue []string) []string {
	v := c.getParam(context, paramName, "")

	if v == "" {
		return defaultValue
	}

	return []string{v}
}

type uriParamGetterImpl bool

func (u uriParamGetterImpl) getParam(context *gin.Context, paramName string, defaultValue string) string {
	v := context.Param(paramName)
	if v == "" {
		v = defaultValue
	}

	return v
}
func (u uriParamGetterImpl) getParamAsArray(context *gin.Context, paramName string, defaultValue []string) []string {
	v := u.getParam(context, paramName, "")
	if v == "" {
		return defaultValue
	}

	return []string{v}
}

var defReqName = map[string]bool{
	"ClientIp":      true,
	"ContentType":   true,
	"Referer":       true,
	"UserAgent":     true,
	"Method":        true,
	"Url":           true,
	"Proto":         true,
	"ProtoMajor":    true,
	"ProtoMinor":    true,
	"ContentLength": true,
	"Host":          true,
	"RemoteAddr":    true,
	"RequestURI":    true,
}

type reqGetterImpl bool

func (u reqGetterImpl) getParam(context *gin.Context, paramName string, defaultValue string) string {
	v := ""

	switch paramName {
	case "ClientIp":
		v = context.ClientIP()
	case "ContentType":
		v = context.ContentType()
	case "Referer":
		v = context.Request.Referer()
	case "UserAgent":
		v = context.Request.UserAgent()
	case "Method":
		v = context.Request.Method
	case "Url":
		v = context.Request.URL.String()
	case "Proto":
		v = context.Request.Proto
	case "ProtoMajor":
		v = strconv.FormatInt(int64(context.Request.ProtoMajor), 10)
	case "ProtoMinor":
		v = strconv.FormatInt(int64(context.Request.ProtoMinor), 10)
	case "ContentLength":
		v = strconv.FormatInt(context.Request.ContentLength, 10)
	case "Host":
		v = context.Request.Host
	case "RemoteAddr":
		v = context.Request.RemoteAddr
	case "RequestURI":
		v = context.Request.RequestURI
	default:
		panic(fmt.Sprintf("Unknown name on request: [%s]", paramName))
	}

	v = strings.TrimSpace(v)
	if v == "" {
		v = defaultValue
	}

	return v
}
func (u reqGetterImpl) getParamAsArray(context *gin.Context, paramName string, defaultValue []string) []string {
	value := u.getParam(context, paramName, "")
	if value == "" {
		return defaultValue
	}

	return []string{value}
}

type boolParamChecker func(context *gin.Context, paramName string) bool
