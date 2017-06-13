package http

import (
	"net/http"
	"net/http/httptest"

	sjson "github.com/bitly/go-simplejson"
	"github.com/dghubble/sling"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tests the initialization of *ResponseResult by \"*sling.Sling\"", func() {
	var sampleServer *httptest.Server
	var slingClient *sling.Sling

	BeforeEach(func() {
		sampleServer = httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte("[3, 55, 17]"))
			}),
		)
		slingClient = sling.New().Base(sampleServer.URL).Get("")
	})
	AfterEach(func() {
		sampleServer.Close()
	})

	It("Get body as string", func() {
		testedResult := NewResponseResultBySling(slingClient)
		Expect(testedResult.Response.StatusCode).To(Equal(200))

		Expect(testedResult.GetBodyAsString()).To(Equal("[3, 55, 17]"))
		By("Get body again")
		Expect(testedResult.GetBodyAsString()).To(Equal("[3, 55, 17]"))
	})

	It("Get body as JSON", func() {
		testedResult := NewResponseResultBySling(slingClient)
		Expect(testedResult.Response.StatusCode).To(Equal(200))

		expectedJson, _ := sjson.NewJson([]byte(`[3, 55, 17]`))
		Expect(testedResult.GetBodyAsJson()).To(Equal(expectedJson))
		By("Get body again")
		Expect(testedResult.GetBodyAsJson()).To(Equal(expectedJson))
	})
})
