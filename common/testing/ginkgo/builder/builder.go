package builder

import (
	"reflect"
	"sync"

	"github.com/juju/errors"

	. "github.com/onsi/ginkgo"
)

func NewGinkgoBuilder(text string) *GinkgoBuilder {
	return &GinkgoBuilder{
		mainText: text,

		testingFuncs: make([]func(), 0, 4),
		testingTable: make([]*GinkgoTable, 0),

		beforeEach:     make([]*eventParams, 0, 1),
		justBeforeEach: make([]*eventParams, 0, 1),
		afterEach:      make([]*eventParams, 0, 1),

		beforeFirst: make([]*eventParams, 0, 1),
		afterLast:   make([]*eventParams, 0, 1),
	}
}

// A cascading pattern for building Ginkgo test environment.
type GinkgoBuilder struct {
	mainText string

	testingFuncs []func()
	testingTable []*GinkgoTable

	beforeEach     []*eventParams
	justBeforeEach []*eventParams
	afterEach      []*eventParams

	beforeFirst []*eventParams
	afterLast   []*eventParams
}

func (b *GinkgoBuilder) BeforeEach(body interface{}, timeout ...float64) *GinkgoBuilder {
	b.beforeEach = append(b.beforeEach, &eventParams{body, timeout})
	return b
}
func (b *GinkgoBuilder) JustBeforeEach(body interface{}, timeout ...float64) *GinkgoBuilder {
	b.justBeforeEach = append(b.justBeforeEach, &eventParams{body, timeout})
	return b
}
func (b *GinkgoBuilder) AfterEach(body interface{}, timeout ...float64) *GinkgoBuilder {
	b.afterEach = append(b.afterEach, &eventParams{body, timeout})
	return b
}
func (b *GinkgoBuilder) BeforeFirst(body interface{}, timeout ...float64) *GinkgoBuilder {
	b.beforeFirst = append(b.beforeFirst, &eventParams{body, timeout})
	return b
}
func (b *GinkgoBuilder) AfterLast(body interface{}, timeout ...float64) *GinkgoBuilder {
	b.afterLast = append(b.afterLast, &eventParams{body, timeout})
	return b
}
func (b *GinkgoBuilder) It(text string, body interface{}, timeout ...float64) *GinkgoBuilder {
	b.appendTesting(func() {
		It(text, body, timeout...)
	})
	return b
}
func (b *GinkgoBuilder) FIt(text string, body interface{}, timeout ...float64) *GinkgoBuilder {
	b.appendTesting(func() {
		FIt(text, body, timeout...)
	})
	return b
}
func (b *GinkgoBuilder) PIt(text string, args ...interface{}) *GinkgoBuilder {
	b.appendTesting(func() {
		PIt(text, args...)
	})
	return b
}
func (b *GinkgoBuilder) XIt(text string, args ...interface{}) *GinkgoBuilder {
	b.appendTesting(func() {
		XIt(text, args)
	})
	return b
}
func (b *GinkgoBuilder) Measure(text string, body interface{}, samples int) *GinkgoBuilder {
	b.appendTesting(func() {
		Measure(text, body, samples)
	})
	return b
}
func (b *GinkgoBuilder) FMeasure(text string, body interface{}, samples int) *GinkgoBuilder {
	b.appendTesting(func() {
		FMeasure(text, body, samples)
	})
	return b
}
func (b *GinkgoBuilder) PMeasure(text string, args ...interface{}) *GinkgoBuilder {
	b.appendTesting(func() {
		PMeasure(text, args...)
	})
	return b
}
func (b *GinkgoBuilder) XMeasure(text string, args ...interface{}) *GinkgoBuilder {
	b.appendTesting(func() {
		XMeasure(text, args...)
	})
	return b
}
func (b *GinkgoBuilder) Specify(text string, body interface{}, timeout ...float64) *GinkgoBuilder {
	b.appendTesting(func() {
		Specify(text, body, timeout...)
	})
	return b
}
func (b *GinkgoBuilder) FSpecify(text string, body interface{}, timeout ...float64) *GinkgoBuilder {
	b.appendTesting(func() {
		FSpecify(text, body, timeout...)
	})
	return b
}
func (b *GinkgoBuilder) PSpecify(text string, args ...interface{}) *GinkgoBuilder {
	b.appendTesting(func() {
		PSpecify(text, args...)
	})
	return b
}
func (b *GinkgoBuilder) XSpecify(text string, args ...interface{}) *GinkgoBuilder {
	b.appendTesting(func() {
		XSpecify(text, args...)
	})
	return b
}
func (b *GinkgoBuilder) Table(table *GinkgoTable) *GinkgoBuilder {
	b.testingTable = append(b.testingTable, table)
	return b
}
func (b *GinkgoBuilder) ToFunc() func() {
	return func() {
		b.Expose()
	}
}

// By the definition of Ginkgo, this function exposes the corresponding Ginkgo function.
func (b *GinkgoBuilder) Expose() {
	countingOfTests := b.getNumberOfTests()

	/**
	 * Before First
	 */
	for _, targetParam := range b.beforeFirst {
		BeforeEach(
			newBeforeFirstFunc(targetParam.body.(func()), countingOfTests),
			targetParam.timeout...,
		)
	}
	// :~)

	/**
	 * Before Each/After Each
	 */
	for _, beforeEachParam := range b.beforeEach {
		BeforeEach(beforeEachParam.body, beforeEachParam.timeout...)
	}
	for _, afterEachParam := range b.afterEach {
		AfterEach(afterEachParam.body, afterEachParam.timeout...)
	}
	// :~)

	/**
	 * After Last
	 */
	for _, targetParam := range b.afterLast {
		AfterEach(
			newAfterLastFunc(targetParam.body.(func()), countingOfTests),
			targetParam.timeout...,
		)
	}
	// :~)

	/**
	 * JustBeforeEach
	 */
	for _, targetParam := range b.justBeforeEach {
		JustBeforeEach(targetParam.body, targetParam.timeout...)
	}
	// :~)

	/**
	 * Executes defined tests by table cases
	 */
	for _, table := range b.testingTable {
		table.Expose()
	}
	// :~)

	/**
	 * Executes defined tests individually
	 */
	for _, execFunc := range b.testingFuncs {
		execFunc()
	}
	// :~)
}
func (b *GinkgoBuilder) ToDescribe() bool {
	return Describe(b.mainText, b.ToFunc())
}
func (b *GinkgoBuilder) ToFDescribe() bool {
	return FDescribe(b.mainText, b.ToFunc())
}
func (b *GinkgoBuilder) ToPDescribe() bool {
	return PDescribe(b.mainText, b.ToFunc())
}
func (b *GinkgoBuilder) ToXDescribe() bool {
	return XDescribe(b.mainText, b.ToFunc())
}
func (b *GinkgoBuilder) ToContext() bool {
	return Context(b.mainText, b.ToFunc())
}
func (b *GinkgoBuilder) ToFContext() bool {
	return FContext(b.mainText, b.ToFunc())
}
func (b *GinkgoBuilder) ToPContext() bool {
	return PContext(b.mainText, b.ToFunc())
}
func (b *GinkgoBuilder) ToXContext() bool {
	return XContext(b.mainText, b.ToFunc())
}
func (b *GinkgoBuilder) appendTesting(testingFunc func()) {
	b.testingFuncs = append(b.testingFuncs, testingFunc)
}
func (b *GinkgoBuilder) getNumberOfTests() int {
	finalCount := len(b.testingFuncs)

	for _, table := range b.testingTable {
		finalCount += table.GetTotalNumberOfFuncs()
	}

	return finalCount
}

func NewGinkgoTable() *GinkgoTable {
	return &GinkgoTable{
		testingType: tt_it,

		execBodies: make([]reflect.Value, 0, 1),
		cases:      make([]*caseContent, 0, 2),
	}
}

/**
 * Testing type
 */
const (
	tt_it = iota
	tt_fit
	tt_pit
	tt_xit
	tt_specify
	tt_fspecify
	tt_pspecify
	tt_xspecify
)

// :~)

var mapOfBuildingTestFunc = map[int]func(text string, targetFunc reflect.Value, params []reflect.Value) func(){
	tt_it: func(text string, targetFunc reflect.Value, params []reflect.Value) func() {
		return func() {
			It(text, func() {
				targetFunc.Call(params)
			})
		}
	},
	tt_fit: func(text string, targetFunc reflect.Value, params []reflect.Value) func() {
		return func() {
			FIt(text, func() {
				targetFunc.Call(params)
			})
		}
	},
	tt_pit: func(text string, targetFunc reflect.Value, params []reflect.Value) func() {
		return func() {
			PIt(text, func() {
				targetFunc.Call(params)
			})
		}
	},
	tt_xit: func(text string, targetFunc reflect.Value, params []reflect.Value) func() {
		return func() {
			XIt(text, func() {
				targetFunc.Call(params)
			})
		}
	},
	tt_specify: func(text string, targetFunc reflect.Value, params []reflect.Value) func() {
		return func() {
			Specify(text, func() {
				targetFunc.Call(params)
			})
		}
	},
	tt_fspecify: func(text string, targetFunc reflect.Value, params []reflect.Value) func() {
		return func() {
			FSpecify(text, func() {
				targetFunc.Call(params)
			})
		}
	},
	tt_pspecify: func(text string, targetFunc reflect.Value, params []reflect.Value) func() {
		return func() {
			PSpecify(text, func() {
				targetFunc.Call(params)
			})
		}
	},
	tt_xspecify: func(text string, targetFunc reflect.Value, params []reflect.Value) func() {
		return func() {
			XSpecify(text, func() {
				targetFunc.Call(params)
			})
		}
	},
}

// A cascading pattern for building test cases as table-like paradigm.
//
// The default test block would be "It()".
type GinkgoTable struct {
	testingType int

	execBodies []reflect.Value
	cases      []*caseContent
}

func (t *GinkgoTable) Exec(body interface{}) *GinkgoTable {
	funcValue := reflect.ValueOf(body)
	if funcValue.Kind() != reflect.Func {
		panic(errors.Details(
			errors.New("Need to be function for Exec(body)."),
		))
	}

	t.execBodies = append(t.execBodies, funcValue)
	return t
}
func (t *GinkgoTable) Case(descBody interface{}, params ...interface{}) *GinkgoTable {
	valueOfDescBody := reflect.ValueOf(descBody)

	var finalDesc string

	switch valueOfDescBody.Kind() {
	case reflect.Func:
		valueOfParams := make([]reflect.Value, 0, len(params))
		for _, v := range params {
			valueOfParams = append(valueOfParams, reflect.ValueOf(v))
		}

		finalDesc = valueOfDescBody.Call(valueOfParams)[0].Interface().(string)
	case reflect.String:
		finalDesc = descBody.(string)
	}

	t.cases = append(t.cases, &caseContent{finalDesc, params})
	return t
}
func (t *GinkgoTable) AsIt() *GinkgoTable {
	t.testingType = tt_it
	return t
}
func (t *GinkgoTable) AsFIt() *GinkgoTable {
	t.testingType = tt_fit
	return t
}
func (t *GinkgoTable) AsPIt() *GinkgoTable {
	t.testingType = tt_pit
	return t
}
func (t *GinkgoTable) AsXIt() *GinkgoTable {
	t.testingType = tt_xit
	return t
}
func (t *GinkgoTable) AsSpecify() *GinkgoTable {
	t.testingType = tt_specify
	return t
}
func (t *GinkgoTable) AsFSpecify() *GinkgoTable {
	t.testingType = tt_fspecify
	return t
}
func (t *GinkgoTable) AsPSpecify() *GinkgoTable {
	t.testingType = tt_pspecify
	return t
}
func (t *GinkgoTable) AsXSpecify() *GinkgoTable {
	t.testingType = tt_xspecify
	return t
}
func (t *GinkgoTable) GetTotalNumberOfFuncs() int {
	return len(t.cases) * len(t.execBodies)
}
func (t *GinkgoTable) ToFunc() func() {
	return func() { t.Expose() }
}

// By the definition of Ginkgo, this function exposes the corresponding Ginkgo "It()" or "Specify().
func (t *GinkgoTable) Expose() {
	for _, targetFunc := range t.prebuildFuncs() {
		targetFunc()
	}
}
func (t *GinkgoTable) prebuildFuncs() []func() {
	totalFuncs := t.GetTotalNumberOfFuncs()
	prebuildFuncs := make([]func(), 0, totalFuncs)

	/**
	 * Builds pre-defined function for calling of
	 * It(), Specify() ... by bodies and cases
	 */
	for _, targetBody := range t.execBodies {
		for _, targetCase := range t.cases {
			params := make([]reflect.Value, len(targetCase.params))
			for i, srcParam := range targetCase.params {
				params[i] = reflect.ValueOf(srcParam)
			}

			prebuildFuncs = append(
				prebuildFuncs,
				mapOfBuildingTestFunc[t.testingType](
					targetCase.desc, targetBody, params,
				),
			)
		}
	}
	// :~)

	return prebuildFuncs
}

type caseContent struct {
	desc   string
	params []interface{}
}

type eventParams struct {
	body    interface{}
	timeout []float64
}

func newBeforeFirstFunc(targetFunc func(), totalNumber int) func() {
	eventObject := newCountingEvent(uint32(totalNumber))
	return func() {
		if !eventObject.isBeforeFirstAndIncrease() {
			return
		}

		targetFunc()
	}
}
func newAfterLastFunc(targetFunc func(), totalNumber int) func() {
	eventObject := newCountingEvent(uint32(totalNumber))
	return func() {
		if !eventObject.isAfterLastOrIncrease() {
			return
		}

		targetFunc()
	}
}

func newCountingEvent(totalNumber uint32) *countingEvent {
	return &countingEvent{
		currentCount: 0,
		totalCount:   totalNumber,
		mutex:        &sync.Mutex{},
	}
}

type countingEvent struct {
	currentCount uint32
	totalCount   uint32
	mutex        *sync.Mutex
}

func (c *countingEvent) isBeforeFirstAndIncrease() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	result := c.currentCount == 0
	if result {
		c.currentCount++
	}

	return result
}
func (c *countingEvent) isAfterLastOrIncrease() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.currentCount++
	result := c.currentCount == c.totalCount

	return result
}
