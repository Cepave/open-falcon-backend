package types

import (
	"fmt"

	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"

	. "github.com/Cepave/open-falcon-backend/common/reflect/types"
	. "gopkg.in/check.v1"
)

type TestBindingSuite struct{}

var _ = Suite(&TestBindingSuite{})

type gBox struct {
	name string
}

func (b *gBox) Bind(sourceObject interface{}) {
	switch typedV := sourceObject.(type) {
	case string:
		b.name = "string: " + typedV
	case int32:
		b.name = fmt.Sprintf("int: %d", typedV)
	default:
		panic("Nothing")
	}
}

// Tests the translation from binding to converter
func (suite *TestBindingSuite) TestBindingToConverter(c *C) {
	testCases := []*struct {
		sourceObj    interface{}
		expectedName string
	}{
		{"Easy", "string: Easy"},
		{int32(91), "int: 91"},
	}

	converter, targetType := BindingToConverter(
		func() interface{} {
			return &gBox{}
		},
	)
	srv := NewDefaultConversionService()
	srv.AddConverter(TypeOfString, targetType, converter)
	srv.AddConverter(TypeOfInt32, targetType, converter)

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedValue := srv.ConvertTo(testCase.sourceObj, targetType).(*gBox)

		c.Assert(testedValue.name, Equals, testCase.expectedName, comment)
	}
}

type MyBox struct {
	Name string
}

func (b *MyBox) Bind(source interface{}) {
	switch v := source.(type) {
	case int:
		b.Name = fmt.Sprintf("Name by number: %d", v)
	case string:
		b.Name = fmt.Sprintf("Name by string: %s", v)
	default:
		panic(fmt.Sprintf("Cannot be converted from value of type: %T", source))
	}
}

func ExampleBinding() {
	/*
		type MyBox struct {
			Name string
		}
		func (b *MyBox) Bind(source interface{}) {
			switch v := source.(type) {
			case int:
				b.Name = fmt.Sprintf("Name by number: %d", v)
			case string:
				b.Name = fmt.Sprintf("Name by string: %s", v)
			default:
				panic(fmt.Sprintf("Cannot be converted from value of type: %T", source))
			}
		}
	*/
	box := &MyBox{}

	DoBinding(44, box)
	fmt.Println(box.Name)
	DoBinding("BC-908", box)
	fmt.Println(box.Name)

	// Output:
	// Name by number: 44
	// Name by string: BC-908
}

func ExampleBinding_conversionService() {
	convSrv := NewDefaultConversionService()

	/*
		type MyBox struct {
			Name string
		}
		func (b *MyBox) Bind(source interface{}) {
			switch v := source.(type) {
			case int:
				b.Name = fmt.Sprintf("Name by number: %d", v)
			case string:
				b.Name = fmt.Sprintf("Name by string: %s", v)
			default:
				panic(fmt.Sprintf("Cannot be converted from value of type: %T", source))
			}
		}
	*/

	// Gets the function of conversion and target type
	convertFunc, targetType := BindingToConverter(func() interface{} {
		return &MyBox{}
	})

	convSrv.AddConverter(TypeOfInt, targetType, convertFunc)
	convSrv.AddConverter(TypeOfString, targetType, convertFunc)

	fmt.Println(convSrv.ConvertTo(89, targetType).(*MyBox).Name)
	fmt.Println(convSrv.ConvertTo("CP-091", targetType).(*MyBox).Name)

	// Output:
	// Name by number: 89
	// Name by string: CP-091
}
