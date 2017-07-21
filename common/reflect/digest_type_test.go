package reflect

import (
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	. "gopkg.in/check.v1"
	"reflect"
)

type TestDigestSuite struct{}

var _ = Suite(&TestDigestSuite{})

// Tests the getting of has value for reflect.Type
func (suite *TestDigestSuite) Test(c *C) {
	type Gd1 struct{}
	type Gd2 struct{}

	type v1Int8 int8
	type v2Int8 int8

	testCases := []*struct {
		type1Value interface{}
		type2Value interface{}
		isSame     Checker
	}{
		/**
		 * Concrete value
		 */
		{int(20), int(33), Equals},
		{int8(71), int16(92), Not(Equals)},
		{v1Int8(71), v2Int8(92), Not(Equals)},
		// :~)
		/**
		 * Slice/Array
		 */
		{[]int{5}, []int{4, 7}, Equals},
		{[]int{5}, []int16{4, 7}, Not(Equals)},
		{[3]int{}, [3]int{}, Equals},
		{[3]int{}, [2]int{}, Not(Equals)},
		// :~)
		/**
		 * Map
		 */
		{map[string]int{}, map[string]int{}, Equals},
		{map[string]int{}, map[string]string{}, Not(Equals)},
		// :~)
		/**
		 * Struct
		 */
		{Gd1{}, Gd1{}, Equals},
		{Gd1{}, Gd2{}, Not(Equals)},
		// :~)
		/**
		 * Channel
		 */
		{make(chan int, 0), make(chan int, 2), Equals},
		{make(chan int), make(chan string), Not(Equals)},
		// :~)
		/**
		 * Functions
		 */
		{func() {}, func() {}, Equals},
		{func(int, int) string { return "" }, func(int, int) string { return "OK" }, Equals},
		{func(int) {}, func(int, int) {}, Not(Equals)},
		{func() int { return 0 }, func() string { return "" }, Not(Equals)},
		// :~)
		/**
		 * Pointers
		 */
		{new(int16), new(int16), Equals},
		{new(*int16), new(*int16), Equals},
		{new(int16), new(int32), Not(Equals)},
		{new(*int16), new(**int16), Not(Equals)},
		{&Gd1{}, &Gd1{}, Equals},
		{&Gd1{}, &Gd2{}, Not(Equals)},
		{new([]string), new([]string), Equals},
		{new([]string), new([]int), Not(Equals)},
		{new(map[string]int), new(map[string]int), Equals},
		{new(map[string]int), new(map[string]string), Not(Equals)},
		// :~)
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		type1 := reflect.TypeOf(testCase.type1Value)
		type2 := reflect.TypeOf(testCase.type2Value)

		c.Logf(
			"Type1: [%s]%s. Type 2: [%s]%s",
			type1.String(), type1.PkgPath(),
			type2.String(), type2.PkgPath(),
		)

		type1HashCode := DigestType(type1)
		type2HashCode := DigestType(type2)

		c.Check(type1HashCode, testCase.isSame, type2HashCode, comment)
	}
}
