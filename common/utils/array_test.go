package utils

import (
	. "gopkg.in/check.v1"
	"reflect"
)

type TestArraySuite struct{}

var _ = Suite(&TestArraySuite{})

// Tests the filtering for array
func (suite *TestArraySuite) TestFilterWith(c *C) {
	sampleData := []int{2, 4, 6, 9}

	testedResult := MakeAbstractArray(sampleData).
		FilterWith(func(v interface{}) bool {
			if v.(int) > 5 {
				return false
			}

			return true
		})

	c.Assert(testedResult.GetArray(), DeepEquals, []int{2, 4})
}

// Tests the mapping for array
func (suite *TestArraySuite) TestMapperTo(c *C) {
	sampleData := []int{1, 3, 5}

	testedResult := MakeAbstractArray(sampleData).
		MapTo(
			func(v interface{}) interface{} {
				return v.(int) + 3
			},
			reflect.TypeOf(int(0)),
		)

	c.Assert(testedResult.GetArray(), DeepEquals, []int{4, 6, 8})
}

// Tests the conversion from typed function to filter
func (suite *TestArraySuite) TestTypedFuncToFilter(c *C) {
	testCases := []*struct {
		testedFunc   FilterFunc
		sampleData   interface{}
		expectedData interface{}
	}{
		{
			TypedFuncToFilter(func(v string) bool { return v == "ok" }),
			[]string{"ok", "skip", "ok"}, []string{"ok", "ok"},
		},
		{
			TypedFuncToFilter(func(v int32) bool { return v > 20 }),
			[]int16{7, 18, 22, 98}, []int16{22, 98},
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := MakeAbstractArray(testCase.sampleData).
			FilterWith(testCase.testedFunc)

		c.Assert(testedResult.GetArray(), DeepEquals, testCase.expectedData, comment)
	}
}

// Tests the conversion from typed function to mapper
func (suite *TestArraySuite) TestTypedFuncToMapper(c *C) {
	testCases := []*struct {
		testedFunc   MapperFunc
		targetType   reflect.Type
		sampleData   interface{}
		expectedData interface{}
	}{
		{
			TypedFuncToMapper(func(v int8) int8 { return v + 3 }),
			reflect.TypeOf(int8(0)),
			[]int8{1, 3, 5}, []int8{4, 6, 8},
		},
		{ // Tests the type conversion
			TypedFuncToMapper(func(v int64) int8 { return int8(v + 2) }),
			reflect.TypeOf(int8(0)),
			[]int64{11, 12}, []int8{13, 14},
		},
		{
			TypedFuncToMapper(func(v string) string { return v + "ok" }),
			reflect.TypeOf(""),
			[]string{"g1:", "g2:"}, []string{"g1:ok", "g2:ok"},
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := MakeAbstractArray(testCase.sampleData).
			MapTo(testCase.testedFunc, testCase.targetType)
		c.Assert(testedResult.GetArray(), DeepEquals, testCase.expectedData, comment)
	}
}

// Tests the unique filter
func (suite *TestArraySuite) TestNewUniqueFilter(c *C) {
	testCases := []*struct {
		targetType   reflect.Type
		sampleData   interface{}
		expectedData interface{}
	}{
		{
			TypeOfString,
			[]string{"A1", "B1", "A2", "B1"},
			[]string{"A1", "B1", "A2"},
		},
		{
			TypeOfInt16,
			[]int16{10, 20, 10, 20},
			[]int16{10, 20},
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := MakeAbstractArray(testCase.sampleData).
			FilterWith(NewUniqueFilter(testCase.targetType))
		c.Assert(testedResult.GetArray(), DeepEquals, testCase.expectedData, comment)
	}
}

// Tests the domain filter
func (suite *TestArraySuite) TestNewDomainFilter(c *C) {
	testCases := []*struct {
		domain       interface{}
		sampleData   interface{}
		expectedData interface{}
	}{
		{
			map[int]bool{1: true, 2: false},
			[]int{1, 2, 3, 4},
			[]int{1, 2},
		},
		{
			map[string]bool{"G1": true, "G2": false},
			[]string{"G1", "G3", "G2"},
			[]string{"G1", "G2"},
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := MakeAbstractArray(testCase.sampleData).
			FilterWith(NewDomainFilter(testCase.domain))
		c.Assert(testedResult.GetArray(), DeepEquals, testCase.expectedData, comment)
	}
}

// Tests the getting array by type conversion
func (suite *TestArraySuite) TestGetArrayAsType(c *C) {
	testCases := []*struct {
		sourceArray    interface{}
		targetValue    interface{}
		expectedResult interface{}
	}{
		{[]int16{11, 16}, uint32(0), []uint32{11, 16}},
		{[]int64{-1, -11}, int8(0), []int8{-1, -11}},
		{[]int32{-13, 109}, int32(0), []int32{-13, 109}},
		{[]int32{}, int8(0), []int8{}},
		{[]int32(nil), int16(0), []int16{}},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := MakeAbstractArray(testCase.sourceArray).
			GetArrayAsTargetType(testCase.targetValue)

		c.Assert(testedResult, DeepEquals, testCase.expectedResult, comment)
	}
}
