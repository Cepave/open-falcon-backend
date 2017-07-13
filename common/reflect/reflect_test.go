package reflect

import (
	"reflect"

	. "github.com/Cepave/open-falcon-backend/common/reflect/types"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	. "gopkg.in/check.v1"
)

type TestReflectSuite struct{}

var _ = Suite(&TestReflectSuite{})

// Tests the getting value of field by strings
func (suite *TestReflectSuite) TestGetValueOfField(c *C) {
	type sampleBox struct {
		Age        int
		Name       string
		ChildBox   *sampleBox
		ChildBoxL2 **sampleBox
	}

	box1 := sampleBox{Age: 30, Name: "b1"}
	box2 := &sampleBox{Age: 33, ChildBox: &box1}
	box3 := &sampleBox{ChildBoxL2: &box2}

	c.Assert(GetValueOfField(box1, "Age"), Equals, 30)

	c.Assert(GetValueOfField(box2, "Age"), Equals, 33)
	c.Assert(GetValueOfField(box2, "ChildBox", "Age"), Equals, 30)

	c.Assert(GetValueOfField(box3, "ChildBoxL2", "Age"), Equals, 33)
	c.Assert(GetValueOfField(box3, "ChildBox"), IsNil)
}

// Tests the getting value of field by strings
func (suite *TestReflectSuite) TestSetValueOfField(c *C) {
	type blueBox struct {
		Age int
	}
	type sampleBox struct {
		Age        int
		Name       string
		ChildBox   *sampleBox
		ChildBoxL2 **sampleBox

		BlueBox blueBox
	}

	box1 := &sampleBox{Age: 35}
	box2 := &sampleBox{Age: 77}
	box3 := &sampleBox{Age: 491}

	/**
	 * Sets the field value by pointer
	 */
	SetValueOfField(box1, 77, "Age")
	c.Assert(box1.Age, Equals, 77)

	SetValueOfField(box1, blueBox{Age: 62}, "BlueBox")
	c.Assert(box1.BlueBox.Age, Equals, 62)

	SetValueOfField(box1, 53, "BlueBox", "Age")
	c.Assert(box1.BlueBox.Age, Equals, 53)

	SetValueOfField(box1, box2, "ChildBox")
	c.Assert(box1.ChildBox.Age, Equals, 77)

	SetValueOfField(box1, &box3, "ChildBoxL2")
	c.Assert((*(box1.ChildBoxL2)).Age, Equals, 491)
	// :~)
}

// Test the extraction for final pointed type
func (suite *TestReflectSuite) TestFinalPointedType(c *C) {
	testCases := []*struct {
		sampleType   interface{}
		expectedType reflect.Type
	}{
		{int32(45), TypeOfInt32},
		{new(**int32), TypeOfInt32},
		{new(****int32), TypeOfInt32},
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedType := FinalPointedType(reflect.TypeOf(testCase.sampleType))

		c.Assert(testedType, DeepEquals, testCase.expectedType, comment)
	}
}

// Tests the getting of final value from pointed value
func (suite *TestReflectSuite) TestGetPointedValue(c *C) {
	v := int(77)
	v1 := &v
	v2 := &v1

	testCases := []*struct {
		sampleValue   interface{}
		expectedValue interface{}
	}{
		{&v, v},
		{v1, v},
		{v2, v},
		{(*uint16)(nil), (*uint16)(nil)},
		{(**uint16)(nil), (**uint16)(nil)},
		{(****uint16)(nil), (****uint16)(nil)},
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedResult := FinalPointedValue(
			reflect.ValueOf(testCase.sampleValue),
		)

		c.Assert(testedResult.Interface(), Equals, testCase.expectedValue, comment)
	}
}

// Tests the new() of final value from pointer type
func (suite *TestReflectSuite) TestNewFinalValue(c *C) {
	testCases := []*struct {
		sampleType    reflect.Type
		expectedValue interface{}
	}{
		{TypeOfInt8, int8(0)},
		{reflect.TypeOf(new(**string)), ""},
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedValue := NewFinalValue(testCase.sampleType).Interface()

		c.Assert(testedValue, Equals, testCase.expectedValue, comment)
	}
}

// Tests the new() of final value from any type to its pointer type(multi-layer)
func (suite *TestReflectSuite) TestNewFinalValueFrom(c *C) {
	v := int32(97)
	v1 := &v
	v2 := &v1
	v3 := &v2

	testCases := []*struct {
		sampleValue    interface{}
		finalType      reflect.Type
		expectedResult interface{}
	}{
		{int64(99), TypeOfInt64, int64(99)},
		{v, reflect.TypeOf(new(*int32)), v2},
		{v, reflect.TypeOf(new(**int32)), v3},
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedValue := NewFinalValueFrom(
			reflect.ValueOf(testCase.sampleValue),
			testCase.finalType,
		)

		c.Logf("New type[%v]. Result type[%v]", testedValue.Type(), reflect.TypeOf(testCase.expectedResult))
		c.Assert(testedValue.Type().String(), Equals, reflect.TypeOf(testCase.expectedResult).String(), comment)

		testedFinalValue := FinalPointedValue(testedValue)
		expectedFinalValue := FinalPointedValue(reflect.ValueOf(testCase.expectedResult))

		c.Assert(testedFinalValue.Interface(), Equals, expectedFinalValue.Interface(), comment)
	}
}
