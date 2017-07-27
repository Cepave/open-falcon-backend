package conform

import (
	"reflect"

	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	"github.com/Cepave/open-falcon-backend/common/utils"
	. "gopkg.in/check.v1"
)

type TestConformSuite struct{}

var _ = Suite(&TestConformSuite{})

type myText string

func (t *myText) ConformSelf(s *ConformService) {
	*t = "Cool!"
}

type NCar struct {
	Name     string `conform:"trim"`
	CaonName string
}

// Tests the conform service
func (suite *TestConformSuite) TestMustConform(c *C) {
	ps := func(v string) *string {
		return &v
	}
	pps := func(v string) **string {
		pp := ps(v)
		return &pp
	}
	ts := func(v string) *myText {
		t := myText(v)
		return &t
	}

	testCases := []*struct {
		value    interface{}
		expected interface{}
		checker  Checker
	}{
		/**
		 * Plain value
		 */
		{" Hello ", "Hello", Equals},
		{[]string{" Ck1 ", " Ck2 "}, []string{"Ck1", "Ck2"}, DeepEquals},
		{[]string(nil), []string(nil), DeepEquals},
		{&([]string{" Ck1 ", " Ck2 "}), &([]string{"Ck1", "Ck2"}), DeepEquals},
		{(*[]string)(nil), (*[]string)(nil), DeepEquals},
		{[2]string{" Ak1 ", " Ak2 "}, [2]string{"Ak1", "Ak2"}, DeepEquals},
		{&([2]string{" Ak1 ", " Ak2 "}), &([2]string{"Ak1", "Ak2"}), DeepEquals},
		// :~)
		/**
		 * Conformer
		 */
		{ts("Hello"), ts("Cool!"), DeepEquals},
		// :~)
		/**
		 * Pointer values
		 */
		{ps(" pc1 "), ps("pc1"), DeepEquals},
		{(*string)(nil), (*string)(nil), DeepEquals},
		{[]*string{ps(" zc1 "), ps(" zc2 ")}, []*string{ps("zc1"), ps("zc2")}, DeepEquals},
		{&([]*string{ps(" zc1 "), ps(" zc2 ")}), &([]*string{ps("zc1"), ps("zc2")}), DeepEquals},
		{[2]*string{ps(" kc1 "), ps(" kc2 ")}, [2]*string{ps("kc1"), ps("kc2")}, DeepEquals},
		{[]**string{pps(" io1 "), pps(" io2 ")}, []**string{pps("io1"), pps("io2")}, DeepEquals},
		{[2]**string{pps(" se1 "), pps(" se2 ")}, [2]**string{pps("se1"), pps("se2")}, DeepEquals},
		{[]*myText{ts("ts-1"), ts("ts-2")}, []*myText{ts("Cool!"), ts("Cool!")}, DeepEquals},
		{&([]*myText{ts("ts-1"), ts("ts-2")}), &([]*myText{ts("Cool!"), ts("Cool!")}), DeepEquals},
		// :~)
		/**
		 * Struct
		 */
		{NCar{" GC-01 ", " CI-20 "}, NCar{"GC-01", " CI-20 "}, DeepEquals},
		{&NCar{" pGC-01 ", " pCI-20 "}, &NCar{"pGC-01", " pCI-20 "}, DeepEquals},
		{(*NCar)(nil), (*NCar)(nil), DeepEquals},
		{
			[]*NCar{{" 2l-car ", " 3l-car "}, {" BG1 ", "C-N1"}},
			[]*NCar{{"2l-car", " 3l-car "}, {"BG1", "C-N1"}},
			DeepEquals,
		},
		{
			&([]*NCar{{" 2l-car ", " 3l-car "}, {" BG1 ", "C-N1"}}),
			&([]*NCar{{"2l-car", " 3l-car "}, {"BG1", "C-N1"}}),
			DeepEquals,
		},
		// :~)
	}

	assertFunc := func(sampleValue interface{}, checker Checker, expectedValue interface{}, comment CommentInterface) {
		newStructType := reflect.StructOf([]reflect.StructField{
			{
				Name: "TestedField",
				Type: reflect.TypeOf(sampleValue),
				Tag:  `conform:"trim"`,
			},
		})
		structPtrValue := reflect.New(newStructType)
		structPtrValue.Elem().Field(0).Set(
			reflect.ValueOf(sampleValue),
		)

		MustConform(structPtrValue.Interface())

		c.Assert(structPtrValue.Elem().Field(0).Interface(), checker, expectedValue, comment)
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		// Asserts value
		assertFunc(testCase.value, testCase.checker, testCase.expected, comment)
	}
}

// Tests the trimming to nil
func (suite *TestConformSuite) TestTrimToNil(c *C) {
	ps := utils.PointerOfCloneString

	testCases := []*struct {
		source   interface{}
		expected interface{}
	}{
		{ps(""), (*string)(nil)},
		{ps("    "), (*string)(nil)},
		{ps(" ACB "), ps("ACB")},
		{[]*string{ps("  "), ps(""), ps(" K-1 ")}, []*string{nil, nil, ps("K-1")}},
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		newStructType := reflect.StructOf([]reflect.StructField{
			{
				Name: "PtrField",
				Type: reflect.TypeOf(testCase.expected),
				Tag:  `conform:"trimToNil"`,
			},
		})
		structPtrValue := reflect.New(newStructType)
		targetField := structPtrValue.Elem().Field(0)

		targetField.Set(reflect.ValueOf(testCase.source))
		MustConform(structPtrValue.Interface())

		c.Logf("Conformed value: %v", targetField.Interface())
		c.Assert(targetField.Interface(), DeepEquals, testCase.expected, comment)
	}
}

// Tests the conform service with not-touched data
func (suite *TestConformSuite) TestMustConformForNotTouchedData(c *C) {
	v1 := &struct {
		Name    string
		Age     int
		BinData []byte
	}{
		Name:    " Nice ",
		Age:     22,
		BinData: []byte{0x43, 0x71, 0x88},
	}

	MustConform(v1)

	c.Assert(v1.Name, Equals, " Nice ")
	c.Assert(v1.Age, Equals, 22)
	c.Assert(v1.BinData, DeepEquals, []byte{0x43, 0x71, 0x88})
}
