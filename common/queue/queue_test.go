package queue

import (
	"math/rand"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

type testStruct struct {
	testStr  string
	testStr2 string
	testStr3 string
	testTime time.Time
}

var _ = Describe("Enqueue(): Put an element in the queue", func() {

	DescribeTable("Testing different types:",
		func(input interface{}) {
			testedQueue := New()
			testedQueue.Enqueue(input)
			Expect(testedQueue.Len()).To(Equal(1))
		},
		Entry("self-defined struct element", &testStruct{"qwer", "asdf", "zxcv", time.Now()}),
		Entry("int element", 123),
	)

	It("Enqueues thread-safely", func() {
		var wg sync.WaitGroup
		testedQueue := New()
		input := []int{ // 100 int elements
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
		}

		for _, elem := range input {
			wg.Add(1)
			go func(e int) {
				defer wg.Done()
				testedQueue.Enqueue(e)
			}(elem)
		}
		wg.Wait()
		Expect(testedQueue.Len()).To(Equal(100))
	})
})

var _ = Describe("DequeueN(): Take up to N elements from the queue", func() {

	DescribeTable("Testing different types:",
		func(input interface{}, expected interface{}) {
			testedQueue := New()
			testedQueue.Enqueue(input)
			r := testedQueue.dequeueN(1)
			Expect(r[0]).To(Equal(expected))
		},
		Entry("int element", &testStruct{"qwer", "asdf", "zxcv", time.Now()}, &testStruct{"qwer", "asdf", "zxcv", time.Now()}),
		Entry("self-defined struct element", 123, 123),
	)

	DescribeTable("Testing different scenarios:",
		func(input []int, num int, expected int) {
			testedQueue := New()
			for _, e := range input {
				testedQueue.Enqueue(e)
			}
			r := testedQueue.dequeueN(num)
			Expect(len(r)).To(Equal(expected))
		},
		Entry("take 2 from a queue of 3 elements", []int{1, 2, 3}, 2, 2),
		Entry("take 3 from a queue of 3 elements", []int{1, 2, 3}, 3, 3),
		Entry("take 4 from a queue of 3 elements", []int{1, 2, 3}, 4, 3),
		Entry("take 0 from a queue of 3 elements", []int{1, 2, 3}, 0, 0),
		Entry("take 1 from an empty queue", []int{}, 1, 0),
		Entry("take 0 from an empty queue", []int{}, 0, 0),
	)

	It("Dequeues thread-safely", func() {
		var wg sync.WaitGroup
		testedQueue := New()
		rand.Seed(time.Now().UTC().UnixNano())
		input := []int{ // 100 int elements
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
		}

		for _, e := range input {
			testedQueue.Enqueue(e)
		}

		t := time.After(300 * time.Millisecond)
		for {
			select {
			case <-t:
				return
			default:
				wg.Add(1)
				go func() {
					defer wg.Done()
					testedQueue.dequeueN(rand.Intn(16))
				}()
			}
		}
		wg.Wait()
		Expect(testedQueue.Len()).To(Equal(0))
	})
})

var _ = Describe("DrainNWithDurationByType(): Take N elements from the queue. It waits up to specific duration for each elements.", func() {
	It("Returns empty slice if Config.Num is zero", func() {
		testedQueue := New()
		input := []int{1, 2, 3}
		for _, e := range input {
			testedQueue.Enqueue(e)
		}
		Expect(testedQueue.DrainNWithDurationByType(&Config{Num: 0, Dur: time.Second}, int(0))).To(Equal([]int{}))
	})

	DescribeTable("Testing empty queue cases:",
		func(c *Config) {
			testedQueue := New()
			Expect(testedQueue.DrainNWithDurationByType(c, int(0))).To(Equal([]int{}))
		},
		Entry("Take 2 elements", &Config{Num: 2, Dur: 300 * time.Millisecond}),
		Entry("Take 0 element", &Config{Num: 0, Dur: 300 * time.Millisecond}),
		Entry("Take 2 elements with zero duration", &Config{Num: 2, Dur: 0}),
		Entry("Take 0 element with zero duration", &Config{Num: 0, Dur: 0}),
	)

	DescribeTable("Testing zero duration cases:",
		func(c *Config, input []int, expected []int) {
			testedQueue := New()
			Expect(testedQueue.DrainNWithDurationByType(c, int(0))).To(Equal([]int{}))
		},
		Entry("Take 2 from queue of 3 elements", &Config{Num: 2, Dur: 0}, []int{1, 2, 3}, []int{1, 2}),
		Entry("Take 3 from queue of 3 elements", &Config{Num: 3, Dur: 0}, []int{1, 2, 3}, []int{1, 2, 3}),
		Entry("Take 4 from queue of 3 elements", &Config{Num: 4, Dur: 0}, []int{1, 2, 3}, []int{1, 2, 3}),
	)

	DescribeTable("Testing generic design with different types:",
		func(c *Config, inputE interface{}, expected interface{}) {
			testedQueue := New()
			testedQueue.Enqueue(inputE)
			r := testedQueue.DrainNWithDurationByType(c, inputE)
			Expect(r).To(Equal(expected))
		},
		Entry("self-defined element", &Config{Num: 1, Dur: 0}, &testStruct{"qwer", "asdf", "zxcv", time.Now()}, []*testStruct{{"qwer", "asdf", "zxcv", time.Now()}}),
		Entry("int element", &Config{Num: 1, Dur: 0}, 123, []int{123}),
	)

	DescribeTable("Testing timeout cases:",
		func(c *Config, input []int, enqInterval time.Duration, expectedResults []int) {
			testedQueue := New()
			go func(i time.Duration) {
				for _, e := range input {
					time.Sleep(i)
					testedQueue.Enqueue(e)
					i *= 2
				}
			}(enqInterval)
			go func() {
				r := testedQueue.DrainNWithDurationByType(c, int(0))
				Expect(r).To(Equal(expectedResults))
			}()
		},
		Entry("take 2 from queue of 3 elements but incapable to take 2nd element", &Config{Num: 2, Dur: 500 * time.Millisecond}, []int{1, 2, 3}, 400*time.Millisecond, []int{1}),
		Entry("take 4 from queue of 3 elements but the 4th never comes", &Config{Num: 4, Dur: 500 * time.Millisecond}, []int{1, 2, 3}, 100*time.Millisecond, []int{1, 2, 3}),
	)

	DescribeTable("Testing cases of multiple enqueuing goroutines one draining",
		func(inputLen int, enqInterval int) {
			testedQueue := New()
			rand.Seed(time.Now().UTC().UnixNano())
			cnt := inputLen
			var cntMutex sync.Mutex
			cntDone := make(chan bool)

			// Randomized input queue
			var inputElems []int
			for i := 0; i < inputLen; i++ {
				inputElems = append(inputElems, rand.Int())
			}

			go func() {
				for i, e := range inputElems {
					go func(i int, e int) {
						// Randomized enqueuing intervals
						time.Sleep(time.Duration(rand.Intn(enqInterval)) * time.Millisecond)
						testedQueue.Enqueue(e)

						cntMutex.Lock()
						if cnt -= 1; cnt == 0 {
							cntDone <- true
						}
						cntMutex.Unlock()
					}(i, e)
				}
			}()

			var r []int
			for {
				select {
				case <-cntDone:
					res := testedQueue.DrainNWithDurationByType(&Config{Num: inputLen, Dur: 0}, int(0)).([]int)
					r = append(r, res...)
					Expect(len(r)).To(Equal(inputLen))
					return
				default:
					// Randomized drain config
					res := testedQueue.DrainNWithDurationByType(&Config{Num: rand.Intn(16), Dur: time.Duration(rand.Intn(enqInterval))}, int(0)).([]int)
					r = append(r, res...)
				}
			}
		},
		Entry("100 enqueuing goroutines", 100, 400),
		Entry("200 enqueuing goroutines", 200, 400),
		Entry("400 enqueuing goroutines", 400, 400),
		Entry("800 enqueuing goroutines", 800, 400),
	)

})
