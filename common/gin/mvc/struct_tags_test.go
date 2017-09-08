package mvc

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Cepave/open-falcon-backend/common/model"
	rt "github.com/Cepave/open-falcon-backend/common/reflect/types"
	ot "github.com/Cepave/open-falcon-backend/common/types"
	"github.com/Cepave/open-falcon-backend/common/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Param loader", func() {
	Context("Query Parameters", func() {
		reqSetup := func(context *gin.Context, sampleValues url.Values) {
			context.Request = httptest.NewRequest("GET", "/my-resource-query?"+sampleValues.Encode(), nil)
		}

		buildDescribeTable("query", reqSetup, SampleSlice|SampleCheckedParam)
	})

	Context("Form Parameters", func() {
		reqSetup := func(context *gin.Context, sampleValues url.Values) {
			context.Request = httptest.NewRequest(
				"POST", "/my-resource-1-form",
				strings.NewReader(sampleValues.Encode()),
			)
			context.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}

		buildDescribeTable("form", reqSetup, SampleSlice|SampleCheckedParam)
	})

	Context("Header Parameters", func() {
		reqSetup := func(context *gin.Context, sampleValues url.Values) {
			context.Request = httptest.NewRequest("GET", "/my-resource-1-header", nil)

			header := context.Request.Header

			for key, values := range sampleValues {
				for _, singleValue := range values {
					header.Add(key, singleValue)
				}
			}
		}

		buildDescribeTable("header", reqSetup, SampleSlice|SampleCheckedParam)
	})

	Context("Context Key/Value", func() {
		reqSetup := func(context *gin.Context, sampleValues url.Values) {
			context.Request = httptest.NewRequest("GET", "/my-resource-1-context", nil)

			for key, values := range sampleValues {
				switch len(values) {
				case 1:
					context.Set(key, values[0])
				default:
					context.Set(key, values)
				}
			}
		}

		buildDescribeTable("key", reqSetup, SampleSlice|SampleCheckedParam)
	})

	Context("Cookie Parameters", func() {
		reqSetup := func(context *gin.Context, sampleValues url.Values) {
			context.Request = httptest.NewRequest("GET", "/my-resource-1-header", nil)

			header := context.Request.Header

			for key, values := range sampleValues {
				for _, singleValue := range values {
					header.Add("Cookie", fmt.Sprintf("%s=%s", key, singleValue))
				}
			}
		}

		buildDescribeTable("cookie", reqSetup, SampleCheckedParam)
	})

	Context("URI Parameters", func() {
		reqSetup := func(context *gin.Context, sampleValues url.Values) {
			context.Request = httptest.NewRequest("GET", "/my-resource-1-uri", nil)

			for key, values := range sampleValues {
				for _, singleValue := range values {
					context.Params = append(context.Params, gin.Param{key, singleValue})
				}
			}
		}

		buildDescribeTable("param", reqSetup, 0)
	})
})

var _ = Describe("Paging object", func() {
	fieldType, _ := reflect.TypeOf(struct {
		P1 *model.Paging `mvc:"pageSize[17] pageOrderBy[ak_1:bd_2]"`
	}{}).FieldByName("P1")
	convSrv := ot.NewDefaultConversionService()

	Context("By Param loader", func() {
		DescribeTable("Check expected properties of page",
			func(reqSetup func(*http.Request), expectedPageSize int, expectedOrderBy []*model.OrderByEntity) {
				context := &gin.Context{
					Request: httptest.NewRequest(http.MethodPost, "/", nil),
				}
				reqSetup(context.Request)
				paramLoader := buildParamLoader(fieldType, convSrv)
				testedPaging := paramLoader(context).(*model.Paging)

				GinkgoT().Logf("Result paging: %#v", testedPaging)
				Expect(testedPaging.Size).To(BeEquivalentTo(expectedPageSize))
				Expect(testedPaging.OrderBy).To(Equal(expectedOrderBy))
			},
			Entry(
				"Default paging", func(req *http.Request) {},
				17,
				[]*model.OrderByEntity{
					{Expr: "ak_1", Direction: utils.DefaultDirection},
					{Expr: "bd_2", Direction: utils.DefaultDirection},
				},
			),
			Entry(
				"Viable paging", func(req *http.Request) {
					req.Header = make(http.Header)
					req.Header.Set("page-size", "39")
					req.Header.Set("order-by", "cp_1#asc:cp_2#desc")
				},
				39,
				[]*model.OrderByEntity{
					{Expr: "cp_1", Direction: utils.Ascending},
					{Expr: "cp_2", Direction: utils.Descending},
				},
			),
		)
	})
})

type SliceAliasG1 []int
type SliceAliasG2 []int

func setupLoaderAndGetValue(sampleType reflect.Type, sampleTag string, contextSetup func(context *gin.Context)) interface{} {
	/**
	 * Set-up the conversion service with converter
	 */
	convSrv := ot.NewDefaultConversionService()
	convSrv.AddConverter(
		rt.TypeOfString, reflect.TypeOf(SliceAliasG1{}),
		func(v interface{}) interface{} {
			return SliceAliasG1{88, 91}
		},
	)
	convSrv.AddConverter(
		rt.STypeOfString, reflect.TypeOf(SliceAliasG2{}),
		func(v interface{}) interface{} {
			return SliceAliasG2{101, 131}
		},
	)
	// :~)

	/**
	 * Build loader
	 */
	fieldType := reflect.StructField{
		Name: "FieldValue1",
		Type: sampleType,
		Tag:  reflect.StructTag(sampleTag),
	}
	testedParamLoader := buildParamLoader(fieldType, convSrv)
	// :~)

	if testedParamLoader == nil {
		return nil
	}

	/**
	 * Gets the value from loader
	 */
	context := &gin.Context{}
	contextSetup(context)
	return testedParamLoader(context)
	// :~)
}

const (
	SampleSlice        = 0x01
	SampleCheckedParam = 0x02
)

func buildDescribeTable(
	paramHolderName string,
	contextSetup func(context *gin.Context, sampleValues url.Values),
	sampleSet int,
) {
	sampleValues := url.Values{
		"vv-1": []string{"hello3"},
		"vv-2": []string{"12"},
		"sv-1": []string{"s1", "s2"},
		"sv-2": []string{"11", "21"},
	}

	rawContextSetup := func(context *gin.Context) {
		contextSetup(context, sampleValues)
	}

	entries := []TableEntry{
		Entry("Viable value", rt.TypeOfString, `mvc:"%s[vv-1]"`, "hello3"),
		Entry("Viable value(type converted by default)", rt.TypeOfInt16, `mvc:"%s[vv-2]"`, int16(12)),
		Entry("Viable value(converted by string -> <type>)", reflect.TypeOf(SliceAliasG1{}), `mvc:"%s[vv-1]"`, SliceAliasG1{88, 91}),
		Entry("Viable value(converted by []string -> <type>)", reflect.TypeOf(SliceAliasG2{}), `mvc:"%s[sv-1]"`, SliceAliasG2{101, 131}),
		Entry("By default value", rt.TypeOfString, `mvc:"%s[non-1] default[gc1]"`, "gc1"),
		Entry("By default value value(slice)", rt.STypeOfInt64, `mvc:"%s[non-1] default[31,41,51]"`, []int64{31, 41, 51}),
		Entry("Empty value(golang's default)", rt.TypeOfInt, `mvc:"%s[non-2]"`, 0),
		Entry("Not a mvc tag", rt.TypeOfInt, `abc:"non"`, nil),
	}

	if sampleSet&SampleSlice > 0 {
		entries = append(
			entries,
			Entry("Viable value(slice)", rt.STypeOfString, `mvc:"%s[sv-1]"`, []string{"s1", "s2"}),
			Entry("Viable value(slice, type converted)", rt.STypeOfUint16, `mvc:"%s[sv-2]"`, []uint16{11, 21}),
			Entry("By default value value(slice)", rt.STypeOfInt64, `mvc:"%s[non-1] default[31,41,51]"`, []int64{31, 41, 51}),
		)
	}

	if sampleSet&SampleCheckedParam > 0 {
		entries = append(
			entries,
			Entry("Checked param(true)", rt.TypeOfBool, `mvc:"%s[?vv-1]"`, true),
			Entry("Checked param(false)", rt.TypeOfBool, `mvc:"%s[?non-1]"`, false),
			Entry("Checked param with default value(false)", rt.TypeOfBool, `mvc:"%s[?non-1] default[d1]"`, false),
		)
	}

	DescribeTable("Match expected value",
		func(sampleType reflect.Type, tagTemplate string, expectedValue interface{}) {
			testedValue := setupLoaderAndGetValue(
				sampleType, fmt.Sprintf(tagTemplate, paramHolderName), rawContextSetup,
			)

			if expectedValue == nil {
				Expect(testedValue).To(BeNil())
				return
			}

			Expect(testedValue).To(Equal(expectedValue))
		},
		entries...,
	)
}
