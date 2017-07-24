package service

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"testing"
)

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}

var _ = Describe("[Unit] Test resolveUrl(...)", func() {
	DescribeTable("should panic when",
		func(host string, resource string) {
			Expect(
				func() {
					_ = resolveUrl(host, resource)
				},
			).To(Panic())
		},
		Entry("Empty host url", "", ""),
		Entry("Invalid scheme", "h@ttp://www.owl.com", "mysqlapi"),
	)

	DescribeTable("should parse successfully when",
		func(host, resource, expectedURL string) {
			url := resolveUrl(host, resource)
			Expect(url).To(Equal(expectedURL))
		},
		Entry("No tailing slash", "http://owl.com", "mysqlapi", "http://owl.com/mysqlapi"),
		Entry("1 tailing slash at host", "http://owl.com/", "mysqlapi", "http://owl.com/mysqlapi"),
		Entry("1 tailing slash at resource", "http://owl.com", "/mysqlapi", "http://owl.com/mysqlapi"),
		Entry("2 tailing slash", "http://owl.com/", "/mysqlapi", "http://owl.com/mysqlapi"),
		Entry("IP:port with no tailing slash", "http://127.0.0.1:5566", "mysqlapi", "http://127.0.0.1:5566/mysqlapi"),
	)
})
