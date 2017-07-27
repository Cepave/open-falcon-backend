package textbuilder

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type sampleStringer bool

func (b sampleStringer) String() string {
	if b {
		return "true"
	}

	return ""
}

var _ = Describe("Tests IsViable(<value>)", func() {
	caseIt := func(text string, testedValue interface{}, expectedValue bool) {
		It(text, func() {
			Expect(IsViable(testedValue)).To(Equal(expectedValue))
		})
	}

	Context("<value> is viable(non-empty)", func() {
		caseIt("String", "AC01", true)
		caseIt("fmt.Stringer", sampleStringer(true), true)
		caseIt("TextGetter", Dsl.S("GP-01"), true)
		caseIt("[]int", []int{3, 4}, true)
		caseIt("map[int]bool", map[int]bool{3: true, 4: false}, true)

		sampleChan := make(chan bool, 2)
		sampleChan <- true
		caseIt("channal", sampleChan, true)
	})

	Context("<value> is not viable(empty)", func() {
		caseIt("String", "", false)
		caseIt("fmt.Stringer", sampleStringer(false), false)
		caseIt("TextGetter", Dsl.S(""), false)
		caseIt("[]int", []int{}, false)
		caseIt("map[int]bool", map[int]bool{}, false)

		sampleChan := make(chan bool, 2)
		caseIt("channal", sampleChan, false)
	})
})

var _ = Describe("Tests TextGetterPrintf()", func() {
	It(`TextGetterPrintf("%v - %v", "Your age", 39)`, func() {
		Expect(TextGetterPrintf("%v - %v", "Your age", 39).String()).
			To(Equal("Your age - 39"))
	})
})

var _ = Describe("Tests ToTextGetter()", func() {
	caseIt := func(text string, testedValue interface{}, expectedValue string) {
		It(text, func() {
			Expect(ToTextGetter(testedValue).String()).To(Equal(expectedValue))
		})
	}

	caseIt("TextGetter", Dsl.S("Hello"), "Hello")
	caseIt("String", "Nice", "Nice")
	caseIt("fmt.Stringer", sampleStringer(true), "true")
	caseIt("int", 30, "30")
})

var _ = Describe("Tests Prefix(prefix, <value>)", func() {
	caseIt := func(text string, testedValue StringGetter, expectedValue string) {
		It(text, func() {
			Expect(Prefix(Dsl.S("Cool:"), testedValue).String()).
				To(Equal(expectedValue))
		})
	}

	caseIt("<value> is viable", "30", "Cool:30")
	caseIt("<value> is not viable", "", "")
})

var _ = Describe("Tests Suffix(<value>, suffix)", func() {
	caseIt := func(text string, testedValue StringGetter, expectedValue string) {
		It(text, func() {
			Expect(Suffix(testedValue, Dsl.S(":HERE")).String()).
				To(Equal(expectedValue))
		})
	}

	caseIt("<value> is viable", "99", "99:HERE")
	caseIt("<value> is not viable", "", "")
})

var _ = Describe("Tests Surrounding(prefix, <value>, suffix)", func() {
	caseIt := func(text string, prefix StringGetter, testedValue StringGetter, suffix StringGetter, expectedValue string) {
		It(text, func() {
			Expect(Surrounding(prefix, testedValue, suffix).String()).
				To(Equal(expectedValue))
		})
	}

	Context("Surround is non-empty", func() {
		caseIt("<value> is viable", "{", "Hello", "}", "{Hello}")
		caseIt("<value> is not viable", "[", "", "]", "")
	})
	Context("Surround is empty", func() {
		caseIt("<value> is viable", "", "Hello", "", "Hello")
		caseIt("<value> is not viable", "", "", "", "")
	})
})

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

var _ = Describe("Tests JoinTextList(joinChar, <list>)", func() {
	caseIt := func(text string, testedValue stringGetters, expectedValue string) {
		It(text, func() {
			Expect(JoinTextList(Dsl.S(", "), testedValue).String()).
				To(Equal(expectedValue))
		})
	}

	caseIt("<list> is non-empty", []string{"A1", "A2", "A3"}, "A1, A2, A3")
	caseIt("<list> is non-empty with empty elements", []string{"C1", "", "C2", "", "C3"}, "C1, C2, C3")
	caseIt("<list> is empty", []string{}, "")
	caseIt("All of element of <list> are empty", []string{"", "", ""}, "")
})

type lenValue bool

func (lv lenValue) Len() int {
	if lv {
		return 7
	}

	return 0
}

var _ = Describe("Tests RepeatByLen(text, <object>)", func() {
	caseIt := func(text string, testedValue interface{}, expectedSize int) {
		It(text, func() {
			Expect(RepeatByLen(Dsl.S("!Staff!"), testedValue)).
				To(HaveLen(expectedSize))
		})
	}

	Context("<object> is viable", func() {
		caseIt("ObjectLen", lenValue(true), 7)
		caseIt("String", "HERE!", 5)
		caseIt("[]int", []int{10, 17, 66}, 3)
		caseIt("map[int]bool", map[int]bool{11: true, 19: true, 20: false}, 3)

		sampleChan := make(chan bool, 2)
		sampleChan <- true
		caseIt("channel", sampleChan, 1)
	})

	Context("<object> is not viable", func() {
		caseIt("ObjectLen", lenValue(false), 0)
		caseIt("String", "", 0)
		caseIt("[]int", []int{}, 0)
		caseIt("map[int]bool", map[int]bool{}, 0)

		sampleChan := make(chan bool, 2)
		caseIt("channel", sampleChan, 0)
	})
})

var _ = Describe("Tests operations on Post()", func() {
	caseIt := func(text string, testedValue TextGetter, expectedString string) {
		It(text, func() {
			Expect(testedValue.String()).
				To(Equal(expectedString))
		})
	}

	caseIt("Prefix(), Suffix(), and Surrounding()",
		Dsl.S("HERE").Post().
			Prefix(Dsl.S("Z1 -> ")).
			Suffix(Dsl.S(" <- K3")).
			Surrounding(Dsl.S("<<"), Dsl.S(">>")),
		"<<Z1 -> HERE <- K3>>",
	)
	caseIt("Repeat() and Join()",
		Dsl.S("atom1").Post().
			Repeat(3).Post().Join(Dsl.S(", ")),
		"atom1, atom1, atom1",
	)
})
