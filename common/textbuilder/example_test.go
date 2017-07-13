package textbuilder

import (
	"fmt"
)

func ExampleStringGetter() {
	stringGetter := StringGetter("Your string")

	fmt.Printf("%s", stringGetter.String())

	// Output:
	// Your string
}

type Weight int

func (w Weight) String() string {
	return fmt.Sprintf("Your weight is %d", w)
}

func ExampleStringerGetter() {
	// Weight implments "fmt.Stringer" interface
	var w1 = Weight(77)
	var getter = NewStringerGetter(w1)

	fmt.Printf("%s", getter.String())

	// Output:
	// Your weight is 77
}

// You could use "Viable" function to control the final output of a value
func ExampleDefaultPost_Viable() {
	surrounding := Surrounding(
		StringGetter("["),
		StringGetter("Hello"),
		StringGetter("]"),
	)

	shown := surrounding.Post().Viable(true)
	hidden := surrounding.Post().Viable(false)

	fmt.Printf("1 - %s\n", shown)
	fmt.Printf("2 - %s\n", hidden)

	// Output:
	// 1 - [Hello]
	// 2 -
}
