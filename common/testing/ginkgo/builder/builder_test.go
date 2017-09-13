package builder

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GinkgoBuilder: For To?Describe()", func() {
	var describe = 0
	NewGinkgoBuilder("ToDescribe()").
		It("Sample It", func() {
			describe = 1
		}).
		ToDescribe()
	It("Flag for Describe() should be 1", func() {
		Expect(describe).To(Equal(1))
	})

	var fdescribe = 1
	NewGinkgoBuilder("ToFDescribe()").
		It("Sample It", func() {
			fdescribe = 1
		}).
		ToFDescribe()
	It("Flag for FDescribe() should be 1", func() {
		Expect(fdescribe).To(Equal(1))
	})

	NewGinkgoBuilder("ToPDescribe()").
		// This should be pending
		It("Sample It", func() {
			Expect(1).To(Equal(2))
		}).
		ToPDescribe()

	NewGinkgoBuilder("ToXDescribe()").
		// This should be pending
		It("Sample It", func() {
			Expect(1).To(Equal(2))
		}).
		ToXDescribe()
})

var _ = Describe("GinkgoBuilder: For To?Context()", func() {
	var context = 0
	NewGinkgoBuilder("ToContext()").
		It("Sample It", func() {
			context = 1
		}).
		ToContext()
	It("Flag for Context() should be 1", func() {
		Expect(context).To(Equal(1))
	})

	var fcontext = 0
	NewGinkgoBuilder("ToFContext()").
		It("Sample It", func() {
			fcontext = 1
		}).
		ToFContext()
	It("Flag for FContext() should be 1", func() {
		Expect(fcontext).To(Equal(1))
	})

	NewGinkgoBuilder("ToPContext()").
		It("Sample It", func() {
			Expect(1).To(Equal(2))
		}).
		ToPContext()

	NewGinkgoBuilder("ToXContext()").
		It("Sample It", func() {
			Expect(1).To(Equal(2))
		}).
		ToXContext()
})

var _ = Describe("GinkgoBuilder: For building of testing function", func() {
	Context("?It related functions",
		func() {
			var (
				it  = 0
				fit = 0
			)

			NewGinkgoBuilder("").
				It("For It", func() {
					it = 1
				}).
				FIt("For FIt", func() {
					fit = 1
				}).
				PIt("For PIt", func() {
					// This should be pending
					Expect(1).To(Equal(2))
				}).
				XIt("For XIt", func() {
					// This should be pending
					Expect(1).To(Equal(2))
				}).
				Expose()

			It("Flag for It() should be set to 1", func() {
				Expect(it).To(Equal(1))
			})
			It("Flag for FIt() should be set to 1", func() {
				Expect(fit).To(Equal(1))
			})
		},
	)

	Context("?Specify related functions",
		func() {
			var (
				specify  = 0
				fspecify = 0
			)

			NewGinkgoBuilder("").
				Specify("For Specify", func() {
					specify = 1
				}).
				FSpecify("For FSpecify", func() {
					fspecify = 1
				}).
				PSpecify("For PSpecify", func() {
					// This should be pending
					Expect(1).To(Equal(2))
				}).
				XSpecify("For XSpecify", func() {
					// This should be pending
					Expect(1).To(Equal(2))
				}).
				Expose()

			It("Flag for Specify() should be set to 1", func() {
				Expect(specify).To(Equal(1))
			})
			It("Flag for FSpecify() should be set to 1", func() {
				Expect(fspecify).To(Equal(1))
			})
		},
	)

	Context("?Measure related functions",
		func() {
			var (
				measure  = 0
				fmeasure = 0
			)

			NewGinkgoBuilder("").
				Measure("For Measure", func(b Benchmarker) {
					measure = 1
				}, 1).
				FMeasure("For FMeasure", func(b Benchmarker) {
					fmeasure = 1
				}, 1).
				PMeasure("For PMeasure", func(b Benchmarker) {
					// This should be pending
					Expect(1).To(Equal(2))
				}, 1).
				XMeasure("For XMeasure", func(b Benchmarker) {
					// This should be pending
					Expect(1).To(Equal(2))
				}, 1).
				Expose()

			It("Flag for Measure() should be set to 1", func() {
				Expect(measure).To(Equal(1))
			})
			It("Flag for FMeasure() should be set to 1", func() {
				Expect(fmeasure).To(Equal(1))
			})
		},
	)
})

var _ = Describe("GinkgoBuilder: For BeforeEach/AfterEach/JustBeforeEach", func() {
	var (
		justBeforeEach = 0
		beforeEach     = 0
		afterEach      = 0
	)

	NewGinkgoBuilder("").
		BeforeEach(func() {
			beforeEach++
		}).
		AfterEach(func() {
			afterEach++
		}).
		JustBeforeEach(func() {
			justBeforeEach++
		}).
		Expose()

	It("1st testing. BeforeEach()/JustBeforeEach() get called 1 time", func() {
		Expect(beforeEach).To(Equal(1))
		Expect(justBeforeEach).To(Equal(1))
		Expect(afterEach).To(Equal(0))
	})
	It("2nd testing. BeforeEach()/JustBeforeEach() get called 2 times. AfterEach() get called once", func() {
		Expect(beforeEach).To(Equal(2))
		Expect(justBeforeEach).To(Equal(2))
		Expect(afterEach).To(Equal(1))
	})
})

var _ = Describe("GinkgoBuilder: For BeforeFirst/AfterLast", func() {
	var (
		beforeCalled = 0
		afterCalled  = 0
	)

	NewGinkgoBuilder("").
		BeforeFirst(func() {
			beforeCalled++
		}).
		AfterLast(func() {
			afterCalled++
		}).
		It("1st testing", func() {
			Expect(beforeCalled).To(Equal(1))
			Expect(afterCalled).To(Equal(0))
		}).
		It("2nd testing", func() {
			Expect(beforeCalled).To(Equal(1))
			Expect(afterCalled).To(Equal(0))
		}).
		Expose()

	It("BeforeFirst should get called only once", func() {
		Expect(beforeCalled).To(Equal(1))
	})
	It("AfterLast should get called only once", func() {
		Expect(afterCalled).To(Equal(1))
	})
})

var _ = Describe("GinkgoTable", func() {
	Context("Simple case", func() {
		counter := 0

		NewGinkgoTable().
			Exec(func(v int) {
				counter += v
			}).
			Exec(func(v int) {
				counter += v
			}).
			Case("case 1 for value 2", 2).
			Case(
				func(v int) string {
					return fmt.Sprintf("case 2 for value %d", v)
				}, 3,
			).
			Expose()

		It("The final counter should be #case * #exec_body", func() {
			Expect(counter).To(Equal(10))
		})
	})
})

var _ = Describe("Use GinkgoTable with GinkgoBuilder", func() {
	var (
		beforeFirst  = 0
		beforeEach   = 0
		tableExec    = 0
		externalExec = 0
	)

	NewGinkgoBuilder("Sample Context 1").
		BeforeFirst(func() {
			beforeFirst++
		}).
		BeforeEach(func() {
			beforeEach++
		}).
		Table(NewGinkgoTable().
			Exec(func() {
				tableExec++
			}).
			Case("Exec Case 1").
			Case("Exec Case 2"),
		).
		It("External It 1", func() {
			externalExec++
		}).
		ToContext()

	It("\"before first\" should be 1", func() {
		Expect(beforeFirst).To(Equal(1))
	})
	It("\"before each\" should be 3", func() {
		Expect(beforeEach).To(Equal(3))
	})
	It("\"table exec\" should be 2", func() {
		Expect(tableExec).To(Equal(2))
	})
	It("\"external exec\" should be 1", func() {
		Expect(externalExec).To(Equal(1))
	})
})
