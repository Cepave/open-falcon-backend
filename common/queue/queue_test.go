package queue

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tests putting/draining object on queue(single thread)", func() {
	type sampleCar struct {
		name string
		age  int
	}

	DescribeTable("Testing different types:",
		func(input interface{}) {
			testedQueue := New()
			testedQueue.Enqueue(input)

			drainResult := testedQueue.DrainNWithDuration(&Config{Num: 10})
			Expect(drainResult).To(HaveLen(1))
			Expect(drainResult[0]).To(Equal(input))
		},
		Entry("struct value", &sampleCar{"GT-210", 3}),
		Entry("primitive value", 123),
	)
})

var _ = DescribeTable("Tests draining objects(with correct sequence)(single thread):",
	func(input []int, num int, expected []int) {
		testedQueue := New()
		for _, e := range input {
			testedQueue.Enqueue(e)
		}
		Expect(testedQueue.DrainNWithDurationByType(&Config{Num: num}, int(0))).To(Equal(expected))
	},
	Entry("take 2 from a queue of 3 elements", []int{1, 2, 3}, 2, []int{1, 2}),
	Entry("take 3 from a queue of 3 elements", []int{37, 91, 3}, 3, []int{37, 91, 3}),
	Entry("take 4 from a queue of 3 elements", []int{65, 33, 44}, 4, []int{65, 33, 44}),
	Entry("take 3 from an empty queue", []int{}, 3, []int{}),
	Entry("take 0 from a queue of 3 elements", []int{91, 92}, 0, []int{}),
	Entry("take 0 from an empty queue", []int{}, 0, []int{}),
)

var _ = Describe("Tests Enqueue() by multiple go-routines", func() {
	numberOfElements := 10240
	numberOfProducers := 16

	Measure("Put element(multiple producers, concurrently)", func(b Benchmarker) {
		GinkgoT().Logf("Elements: %d. Producers: %d", numberOfElements, numberOfProducers)
		testedQueue := New()
		b.Time("runtime", func() {
			wg := &sync.WaitGroup{}
			numberOfGoRoutines := make(chan bool, numberOfProducers)
			for i := 0; i < numberOfElements; i++ {
				numberOfGoRoutines <- true
				wg.Add(1)
				go func() {
					defer func() {
						wg.Done()
						<-numberOfGoRoutines
					}()
					testedQueue.Enqueue(rand.Int31n(50000))
				}()
			}

			wg.Wait()
		})

		Expect(testedQueue.Len()).To(Equal(numberOfElements))
		GinkgoT().Logf("[0:4] Elements: %#v", testedQueue.DrainNWithDuration(&Config{Num: 5}))
	}, 8)
})

var _ = Describe("Tests DrainNWithDuration() by multiple go-routines(consumer)", func() {
	numberOfElements := 20240
	drainingBatch := 13
	numberOfConsumers := 8

	buildQueue := func(numberOfElements int) *Queue {
		queue := New()
		for i := 0; i < numberOfElements; i++ {
			queue.Enqueue(rand.Int31n(3450))
		}
		return queue
	}

	Measure("Drains", func(b Benchmarker) {
		queue := buildQueue(numberOfElements)
		wg := &sync.WaitGroup{}

		numberOfDrainedElements := int32(0)
		b.Time("runtime", func() {
			config := &Config{Num: drainingBatch}

			for i := 0; i < numberOfConsumers; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for {
						drainedElements := queue.DrainNWithDuration(config)
						if len(drainedElements) == 0 {
							break
						}

						atomic.AddInt32(&numberOfDrainedElements, int32(len(drainedElements)))
					}
				}()
			}

			wg.Wait()
		})

		Expect(int(numberOfDrainedElements)).To(Equal(numberOfElements))
	}, 8)
})

var _ = Describe("Tests draining before enqueuing(with waiting)", func() {
	It("Waits for certain time when draining", func() {
		testedQueue := New()
		testedQueue.Enqueue(30)
		testedQueue.Enqueue(40)

		var waitingTimeOfDraining time.Duration
		var testedResult []interface{}

		wg := &sync.WaitGroup{}

		/**
		 * Drains data with waiting for 3 seconds
		 */
		wg.Add(1)
		go func() {
			defer wg.Done()
			beforeDraining := time.Now()
			testedResult = testedQueue.DrainNWithDuration(&Config{Num: 3, Dur: 1 * time.Second})
			waitingTimeOfDraining = time.Now().Sub(beforeDraining)
		}()
		// :~)

		/**
		 * Enqueues data after 1 second(the consuemr should un-lock the queue while it is waiting)
		 */
		wg.Add(1)
		go func() {
			defer wg.Done()

			time.Sleep(300 * time.Millisecond)
			testedQueue.Enqueue(9)
		}()
		// :~)

		wg.Wait()

		// If the lock is kept holding by draining, the draining should get only 2 elements.
		Expect(testedResult).To(HaveLen(3))
		// Asserts that the wating has occurred
		Expect(waitingTimeOfDraining).To(BeNumerically(">=", 1*time.Second))
	})
})

var _ = Describe("Tests the type conversion", func() {
	type myFoot struct {
		color  string
		length int
	}

	It("Type conversion", func() {
		testedQueue := New()
		testedQueue.Enqueue(&myFoot{"blue", 331})

		testedResult := testedQueue.DrainNWithDurationByType(&Config{Num: 1}, &myFoot{}).([]*myFoot)
		Expect(testedResult).To(HaveLen(1))
	})
})
