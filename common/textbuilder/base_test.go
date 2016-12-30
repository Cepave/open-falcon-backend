package textbuilder

import (
	. "gopkg.in/check.v1"
)

type TestBaseSuite struct{}

var _ = Suite(&TestBaseSuite{})

type sampleStringer bool
func (b sampleStringer) String() string {
	if b {
		return "true"
	}

	return ""
}

// Tests the function for viable checking
func (suite *TestBaseSuite) TestIsViable(c *C) {
	sampleChan := make(chan bool, 2)
	sampleChan <- true
	emptyChan := make(chan bool, 2)

	testCases := []*struct {
		sampleValue interface{}
		expectedResult bool
	} {
		{ "AC01", true },
		{ "", false },
		{ sampleStringer(true), true },
		{ sampleStringer(false), false },
		{ Dsl.S("AC01"), true },
		{ Dsl.S(""), false },
		{ []int { 20, 30 }, true },
		{ []int {}, false },
		{ map[int]bool { 1: true, 2: false }, true },
		{ map[int]bool {}, false },
		{ sampleChan, true },
		{ emptyChan, false },
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		c.Assert(IsViable(testCase.sampleValue), Equals, testCase.expectedResult, comment)
	}
}

type sampleBoolStringer bool
func (b sampleBoolStringer) String() string {
	return "This is bool"
}

// Tests the build of gette by fmt.Sprintf
func (suite *TestBaseSuite) TestTextGetterPrintf(c *C) {
	sampleGetter := TextGetterPrintf("%v - %v", "Your age", 39)

	c.Assert(sampleGetter.String(), Equals, "Your age - 39")
}

// Tests the building of getter from various type
func (suite *TestBaseSuite) TestToTextGetter(c *C) {
	testCases := []*struct {
		sampleValue interface{}
		expectedResult string
	} {
		{ Dsl.S("Hello"), "Hello" },
		{ "Nice", "Nice" },
		{ sampleBoolStringer(true), "This is bool" },
		{ 30, "30" },
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		testedResult := ToTextGetter(testCase.sampleValue).String()
		c.Assert(testedResult, Equals, testCase.expectedResult, comment)
	}
}

// Tests the prefixing function
func (suite *TestBaseSuite) TestPrefix(c *C) {
	testCases := []*struct {
		sampleValue string
		expected string
	} {
		{ "30", "Cool:30" },
		{ "", "" },
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		testedGetter := Prefix(
			Dsl.S("Cool:"),
			Dsl.S(testCase.sampleValue),
		)
		c.Assert(testedGetter.String(), Equals, testCase.expected, comment)
	}
}

// Tests the suffixing function
func (suite *TestBaseSuite) TestSuffix(c *C) {
	testCases := []*struct {
		sampleValue string
		expected string
	} {
		{ "30", "30:HERE" },
		{ "", "" },
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		testedGetter := Suffix(
			Dsl.S(testCase.sampleValue),
			Dsl.S(":HERE"),
		)
		c.Assert(testedGetter.String(), Equals, testCase.expected, comment)
	}
}

// Tests the suffixing function
func (suite *TestBaseSuite) TestSurrounding(c *C) {
	testCases := []*struct {
		prefixString StringGetter
		sampleValue StringGetter
		suffixString StringGetter
		expected string
	} {
		{ "G1", "99", "G2", "G199G2" },
		{ "G1", "", "G2", "" },
		{ "", "99", "", "99" },
		{ "", "", "", "" },
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		testedGetter := Surrounding(
			testCase.prefixString,
			testCase.sampleValue,
			testCase.suffixString,
		)
		c.Assert(testedGetter.String(), Equals, testCase.expected, comment)
	}
}

type stringGetters []string
func (s stringGetters) Get(index int) TextGetter {
	return Dsl.S(s[index])
}
func (s stringGetters) Len() int {
	return len(s)
}
func (s stringGetters) Post() ListPostProcessor {
	return NewListPost(s)
}

// Tests the join function
func (suite *TestBaseSuite) TestJoin(c *C) {
	testCases := []*struct {
		sampleValues stringGetters
		expected string
	} {
		{
			[]string{ "A1", "", "A2", "", "A3" },
			"A1, A2, A3",
		},
		{
			[]string{ "" },
			"",
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		testedResult := JoinTextList(Dsl.S(", "), testCase.sampleValues)

		c.Assert(testedResult.String(), Equals, testCase.expected, comment)
	}
}

type lenValue bool
func (lv lenValue) Len() int {
	if lv {
		return 7
	}

	return 0
}

// Tests the Repeat by len
func (suite *TestBaseSuite) TestRepeatByLen(c *C) {
	sampleChan := make(chan bool, 2)
	sampleChan <- true
	emptyChan := make(chan bool, 2)

	testCases := []*struct {
		sampleLenObject interface{}
		expectedLen int
	} {
		{ lenValue(true), 7 },
		{ lenValue(false), 0 },
		{ "HERE!", 5 },
		{ "", 0 },
		{ []bool { true, true, true, true }, 4 },
		{ []bool {}, 0 },
		{ map[int]bool { 11: true, 13: true, 17: true, 20: true }, 4 },
		{ map[int]bool {}, 0 },
		{ sampleChan, 1 },
		{ emptyChan, 0 },
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		testedList := RepeatByLen(Dsl.S("Cool!"), testCase.sampleLenObject)

		c.Assert(testedList, HasLen, testCase.expectedLen, comment)
	}
}

// Tests the post function
func (suite *TestBaseSuite) TestPost(c *C) {
	testCases := []*struct {
		sample TextGetter
		expected string
	} {
		{
			Dsl.S("HERE").Post().
				Prefix(Dsl.S("Z1 -> ")).
				Suffix(Dsl.S(" <- K3")).
				Surrounding(Dsl.S("<<"), Dsl.S(">>")),
			"<<Z1 -> HERE <- K3>>",
		},
		{
			Dsl.S("atom1").Post().
				Repeat(3).Post().Join(Dsl.S(", ")),
			"atom1, atom1, atom1",
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		testedGetter := testCase.sample

		c.Assert(testedGetter.String(), Equals, testCase.expected, comment)
	}
}
