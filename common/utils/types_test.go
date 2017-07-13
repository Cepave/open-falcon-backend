package utils

import (
	. "gopkg.in/check.v1"
	"reflect"
)

type TestTypesSuite struct{}

var _ = Suite(&TestTypesSuite{})

type car struct{}
type func1 func()

type g1Err struct{}

func (e *g1Err) Error() string {
	return "OK"
}

// Tests types of any value
func (suite *TestTypesSuite) TestIsViable(c *C) {
	ch1, ch2 := make(chan bool, 1), make(chan bool, 1)
	ch1 <- true

	var nilErr1 error = (*g1Err)(nil)
	var nilErr2 *g1Err = (*g1Err)(nil)

	testCases := []*struct {
		sampleValue interface{}
		expected    bool
	}{
		{30, true},
		{0, true},
		{&car{}, true},
		{(*car)(nil), false},
		{[]int{20}, true},
		{[]int{}, false},
		{[]string(nil), false},
		{[]*car{{}}, true},
		{[]*car{}, false},
		{map[int]bool{20: true}, true},
		{map[int]bool{}, false},
		{func1(func() {}), true},
		{func1(nil), false},
		{ch1, true},
		{ch2, false},
		{nilErr1, false},
		{nilErr2, false},
		{(error)(nil), false},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedValue := reflect.ValueOf(testCase.sampleValue)
		testedResult := ValueExt(testedValue).IsViable()

		c.Assert(testedResult, Equals, testCase.expected, comment)
	}
}

// Tests the type conversion
func (suite *TestTypesSuite) TestConvertToForPointer(c *C) {
	type weight int
	var v1 int = 20

	// Asserts the nil value
	c.Assert(ConvertToTargetType((*string)(nil), new(weight)), Equals, (*weight)(nil))

	/**
	 * Asserts the value of pointer
	 */
	convertedValue := ConvertToTargetType(&v1, new(weight)).(*weight)
	c.Assert(*convertedValue, Equals, weight(20))
	// :~)
}

// Tests the conversion for integer types
func (suite *TestTypesSuite) TestConvertToForReal(c *C) {
	testCases := []*struct {
		sourceValue interface{}
		targetValue interface{}
	}{
		{int8(10), int16(10)},
		{int8(11), int32(11)},
		{int8(12), int64(12)},
		{uint8(13), uint16(13)},
		{uint8(14), uint32(14)},
		{uint8(15), uint64(15)},
		{int8(16), uint16(16)},
		{uint64(33), int8(33)},
		{float64(31.77), float32(31.77)},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := ConvertToTargetType(
			testCase.sourceValue, testCase.targetValue,
		)

		c.Assert(testedResult, Equals, testCase.targetValue, comment)
	}
}

func (suite *TestTypesSuite) TestConvertStruct(c *C) {
	type car struct {
		age  int
		name string
	}

	type carV1 car

	c1 := car{20, "LBUE-98"}

	testedCar := ConvertToTargetType(c1, carV1{}).(carV1)
	c.Assert(testedCar.age, Equals, c1.age)
	c.Assert(testedCar.name, Equals, c1.name)
}
