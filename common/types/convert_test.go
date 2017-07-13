package types

import (
	"fmt"
	"reflect"

	oreflect "github.com/Cepave/open-falcon-backend/common/reflect"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"

	. "github.com/Cepave/open-falcon-backend/common/reflect/types"
	. "gopkg.in/check.v1"
)

type TestConvertSuite struct{}

var _ = Suite(&TestConvertSuite{})

// Tests the defalt service of conversion on added conversion
func (suite *TestConvertSuite) TestAddConverter(c *C) {
	srv := NewDefaultConversionService()
	srv.AddConverter(
		TypeOfInt64,
		TypeOfInt32,
		func(o interface{}) interface{} {
			return int32(o.(int64) + 33)
		},
	)

	c.Assert(srv.CanConvert(TypeOfInt64, TypeOfInt32), Equals, true)
	c.Assert(srv.ConvertTo(int64(77), TypeOfInt32), Equals, int32(110))
}

// Tests the conversion by buildin converters
func (suite *TestConvertSuite) TestConvertByBuildinConverters(c *C) {
	v := "781"
	v1 := &v

	testCases := []*struct {
		source        interface{}
		desiredType   reflect.Type
		expectedValue interface{}
	}{
		{int32(981), TypeOfInt32, int32(981)}, // Same type
		{"45", TypeOfInt16, int16(45)},
		{"true", TypeOfBool, true},
		{int32(46), TypeOfInt64, int64(46)},
		{int32(46), TypeOfString, "46"},
		{&v1, TypeOfUint32, uint32(781)}, // Source is pointer to value
	}

	srv := NewDefaultConversionService()
	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedValue := srv.ConvertTo(testCase.source, testCase.desiredType)

		c.Assert(testedValue, Equals, testCase.expectedValue, comment)
	}
}

// Tests the conversion for complex data type
//
// pointer, array, slice
func (suite *TestConvertSuite) TestConvertToComplexType(c *C) {
	type mapOfThings map[int32]string

	testString := "Ptr-String-1"
	testStringPtr1 := &testString
	testStringPtr2 := &testStringPtr1

	testHello := "Hello-Ptr"
	testGc01Slice := []string{"GC-01"}

	testCases := []*struct {
		source        interface{}
		desiredType   reflect.Type
		expectedValue interface{}
	}{
		{testString, PTypeOfString, &testString},                          // To pointer
		{testString, reflect.TypeOf(new(**string)), &testStringPtr2},      // To pointer(***)
		{"Hello", STypeOfString, []string{"Hello"}},                       // To slice
		{"39", STypeOfInt16, []int16{39}},                                 // To slice of specific type
		{"King-1", reflect.TypeOf([3]string{}), [3]string{"King-1"}},      // To array
		{"91", reflect.TypeOf([3]uint8{}), [3]uint8{91}},                  // To array of specific type
		{"Hello-Ptr", reflect.TypeOf([]*string{}), []*string{&testHello}}, // To slice of pointers
		{"GC-01", reflect.TypeOf(new([]string)), &testGc01Slice},          // To pointer(of slice)
		{
			mapOfThings{32: "Cool", 77: "nothing"},
			reflect.TypeOf(make(map[int32]string)),
			map[int32]string{32: "Cool", 77: "nothing"},
		},
	}

	srv := NewDefaultConversionService()
	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedValue := srv.ConvertTo(testCase.source, testCase.desiredType)

		c.Assert(testedValue, DeepEquals, testCase.expectedValue, comment)
	}
}

// Test the conversion for channel type
func (suite *TestConvertSuite) TestConvertToChannel(c *C) {
	var v1, v2 uint32 = 20, 40

	testCases := []*struct {
		source        interface{}
		desiredType   reflect.Type
		expectedValue interface{}
	}{
		{"Go!!", channelTypeOf(TypeOfString), makeChannelWithValues("Go!!")}, // To channel
		{"48", channelTypeOf(TypeOfByte), makeChannelWithValues(byte(48))},   // To channel of specific type
		/**
		 * From complex types to channel
		 */
		{
			[]int16{20, 40},
			channelTypeOf(TypeOfUint32),
			makeChannelWithValues(uint32(20), uint32(40)),
		},
		{ // To pointers
			[]int16{20, 40},
			reflect.TypeOf(make(chan *uint32)),
			makeChannelWithValues(&v1, &v2),
		},
		{
			[2]int16{71, 73},
			channelTypeOf(TypeOfUint32),
			makeChannelWithValues(uint32(71), uint32(73)),
		},
		{
			makeChannelWithValues(int16(35), int16(57)),
			channelTypeOf(TypeOfUint32),
			makeChannelWithValues(uint32(35), uint32(57)),
		},
		// :~)
	}

	srv := NewDefaultConversionService()
	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedValue := srv.ConvertTo(testCase.source, testCase.desiredType)

		c.Assert(testedValue, ocheck.ChannelEquals, testCase.expectedValue, comment)
	}
}

func channelTypeOf(elemType reflect.Type) reflect.Type {
	return reflect.ChanOf(reflect.BothDir, elemType)
}
func makeChannelWithValues(values ...interface{}) interface{} {
	valueOfElem := reflect.ValueOf(values[0])

	channelType := channelTypeOf(valueOfElem.Type())
	channelValue := reflect.MakeChan(channelType, len(values))

	for _, v := range values {
		channelValue.Send(reflect.ValueOf(v))
	}

	return channelValue.Interface()
}

// Tests complex type to slice
func (suite *TestConvertSuite) TestComplexTypesToSlice(c *C) {
	var v1, v2 uint32 = 20, 40

	testCases := []*struct {
		source      interface{}
		desiredType reflect.Type
		expected    interface{}
	}{
		{
			[]int16{20, 40},
			STypeOfUint32,
			[]uint32{20, 40},
		},
		{ // Nil value of slice
			[]int16(nil),
			STypeOfUint32,
			[]uint32{},
		},
		{ // To pointers
			[]int16{20, 40},
			reflect.TypeOf([]*uint32{}),
			[]*uint32{&v1, &v2},
		},
		{
			[2]int16{71, 73},
			STypeOfUint32,
			[]uint32{71, 73},
		},
		{
			makeChannelWithValues(int16(35), int16(57)),
			STypeOfUint32,
			[]uint32{35, 57},
		},
	}

	srv := NewDefaultConversionService()
	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedResult := srv.ConvertTo(testCase.source, testCase.desiredType)

		c.Assert(testedResult, DeepEquals, testCase.expected, comment)
	}
}

// Tests complex type to array
func (suite *TestConvertSuite) TestComplexTypesToArray(c *C) {
	var v1, v2 uint32 = 20, 40

	testCases := []*struct {
		source      interface{}
		desiredType reflect.Type
		expected    interface{}
	}{
		{
			[]int16{20, 40},
			reflect.TypeOf([2]uint32{}),
			[2]uint32{20, 40},
		},
		{
			[]int16(nil),
			reflect.TypeOf([2]uint32{}),
			[2]uint32{},
		},
		{ // To pointers
			[]int16{20, 40},
			reflect.TypeOf([2]*uint32{}),
			[2]*uint32{&v1, &v2},
		},
		{
			[]int16{20, 40},
			reflect.TypeOf([4]uint32{}),
			[4]uint32{20, 40},
		},
		{
			[2]int16{71, 73},
			reflect.TypeOf([2]uint32{}),
			[2]uint32{71, 73},
		},
		{
			[2]int16{71, 73},
			reflect.TypeOf([4]uint32{}),
			[4]uint32{71, 73},
		},
		{
			makeChannelWithValues(int16(35), int16(57)),
			reflect.TypeOf([2]uint32{}),
			[2]uint32{35, 57},
		},
		{
			makeChannelWithValues(int16(80), int16(81)),
			reflect.TypeOf([4]uint32{}),
			[4]uint32{80, 81},
		},
	}

	srv := NewDefaultConversionService()
	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedResult := srv.ConvertTo(testCase.source, testCase.desiredType)

		c.Assert(testedResult, DeepEquals, testCase.expected, comment)
	}
}

// Tests the conversion of map
func (suite *TestConvertSuite) TestMapConversion(c *C) {
	np := oreflect.NewPointerValue

	testCases := []*struct {
		source     interface{}
		targetType reflect.Type
		expected   interface{}
	}{
		{
			map[string]int8{"91": 32, "104": 17},
			buildMapType(int32(0), uint64(0)),
			map[int32]uint64{91: 32, 104: 17},
		},
		{ // Pointers
			map[string]int8{"91": 32, "104": 17},
			buildMapType(int32(0), new(uint64)),
			map[int32]*uint64{
				91:  np(uint64(32)).(*uint64),
				104: np(uint64(17)).(*uint64),
			},
		},
		{ // nil value
			map[string]int8(nil),
			buildMapType(int32(0), uint64(0)),
			map[int32]uint64{},
		},
	}

	srv := NewDefaultConversionService()
	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedResult := srv.ConvertTo(testCase.source, testCase.targetType)

		c.Assert(testedResult, DeepEquals, testCase.expected, comment)
	}

	/**
	 * Tests the key as complex value
	 */
	testedResult := srv.ConvertTo(
		map[string]int16{"101": 10, "108": 26},
		buildMapType(new(int32), ""),
	).(map[*int32]string)

	c.Logf("Special key of map: %#v", testedResult)
	c.Assert(testedResult, HasLen, 2)
	for key, value := range testedResult {
		switch *key {
		case 101:
			c.Assert(value, Equals, "10")
		case 108:
			c.Assert(value, Equals, "26")
		default:
			c.Errorf("Unknown key: [%v]", key)
		}
	}
	// :~)
}

func buildMapType(keySample interface{}, valueSample interface{}) reflect.Type {
	return reflect.MapOf(
		reflect.TypeOf(keySample),
		reflect.TypeOf(valueSample),
	)
}

func ExampleDefaultConversionService() {
	convSrv := NewDefaultConversionService()

	stringV := convSrv.ConvertTo(45, TypeOfString)
	fmt.Printf("%T:%s\n", stringV, stringV)

	// Output:
	// string:45
}

func ExampleConversionService_customiedConverter() {
	convSrv := NewDefaultConversionService()

	// import ot "github.com/Cepave/open-falcon-backend/common/reflect/types"
	convSrv.AddConverter(
		TypeOfInt8, TypeOfString,
		func(source interface{}) interface{} {
			return fmt.Sprintf("int8:%d", source)
		},
	)

	fmt.Printf("%s\n", convSrv.ConvertTo(int8(44), TypeOfString))

	// Output:
	// int8:44
}

func ExampleConversionService_slice() {
	convSrv := NewDefaultConversionService()

	// import ot "github.com/Cepave/open-falcon-backend/common/reflect/types"
	sourceSlice := []string{"76", "38", "99"}
	targetSlice := convSrv.ConvertTo(sourceSlice, STypeOfInt).([]int)

	fmt.Printf("%#v\n", targetSlice)

	// Output:
	// []int{76, 38, 99}
}

func ExampleConversionService_map() {
	convSrv := NewDefaultConversionService()

	sourceMap := map[string]int32{
		"Key1": 88,
		"Key2": 92,
		"Key3": 63,
	}
	// import ot "github.com/Cepave/open-falcon-backend/common/reflect/types"
	targetMap := convSrv.ConvertTo(
		sourceMap,
		reflect.MapOf(
			TypeOfString, TypeOfUint64,
		),
	).(map[string]uint64)

	fmt.Printf("%d, %d, %d\n", targetMap["Key1"], targetMap["Key2"], targetMap["Key3"])

	// Output:
	// 88, 92, 63
}
