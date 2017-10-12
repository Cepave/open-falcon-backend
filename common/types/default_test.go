package types

import (
	"reflect"

	"github.com/satori/go.uuid"

	t "github.com/Cepave/open-falcon-backend/common/reflect/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("String To Uuid converter", func() {
	Context("Normal conversion", func() {
		It("As expected UUID value", func() {
			sampleUuidString := "2e196e79-de1a-4459-a8ce-473ad6416938"

			testedUuid := stringToUuid(sampleUuidString).(uuid.UUID)

			Expect(testedUuid.String()).To(Equal(sampleUuidString))
		})
	})

	Context("Convesion service", func() {
		testedSrv := NewDefaultConversionService()
		AddDefaultConverters(testedSrv)

		It("Check if types are convertible", func() {
			Expect(testedSrv.CanConvert(
				t.TypeOfString, reflect.TypeOf(uuid.Nil),
			)).To(BeTrue())
		})

		It("Check result from conversion service", func() {
			sampleUuidString := "db85ec02-75f7-4a35-ade4-583ef059e115"
			testedUuid := testedSrv.ConvertTo(sampleUuidString, reflect.TypeOf(uuid.Nil)).(uuid.UUID)

			Expect(testedUuid.String()).To(Equal(sampleUuidString))
		})
	})
})
