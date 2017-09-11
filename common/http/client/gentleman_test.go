package client_test

import (
	"errors"
	"net/http"
	"time"

	"github.com/h2non/gentleman/plugins/timeout"
	je "github.com/juju/errors"
	gt "gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gock.v1"

	tl "github.com/Cepave/open-falcon-backend/common/http/client"
	mock "github.com/Cepave/open-falcon-backend/common/testing/http/gock"
	gp "github.com/Cepave/open-falcon-backend/common/testing/http/gock_plugin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var gockConfig = mock.GockConfigBuilder.NewConfigByRandom()

var _ = Describe("Get JSON from response", func() {
	var client *gt.Client

	AfterEach(func() {
		gockConfig.Off()
		client = nil
	})

	Context("Normal JSON", func() {
		BeforeEach(func() {
			gockConfig.New().Get("/r1").
				Reply(http.StatusOK).JSON(
				map[string]interface{}{
					"v1": 30,
					"v2": "good",
				},
			)

			client = gockConfig.GentlemanT.SetupClient(
				tl.CommonGentleman.NewDefaultClient(),
			).Path("/r1")
		})

		It("JSON matching expected", func() {
			resp, err := client.Get().Send()

			Expect(err).To(Succeed())

			respExt := tl.ToGentlemanResp(resp)
			jsonBody := respExt.MustGetJson()

			Expect(jsonBody.Get("v1").MustInt()).To(BeEquivalentTo(30))
			Expect(jsonBody.Get("v2").MustString()).To(Equal("good"))
		})
	})

	Context("Non-JSON Body", func() {
		BeforeEach(func() {
			gockConfig.New().Get("/txt1").
				Reply(http.StatusOK).
				Type("text").
				BodyString("Non-JSON Body")

			client = gockConfig.GentlemanT.SetupClient(
				tl.CommonGentleman.NewDefaultClient(),
			).Path("/txt1")
		})

		It("Error should occur", func() {
			resp, err := client.Get().Send()

			Expect(err).To(Succeed())

			respExt := tl.ToGentlemanResp(resp)
			_, err = respExt.GetJson()

			Expect(err).To(HaveOccurred())
			GinkgoT().Logf("Error for parsing JSON: %v", je.Details(err))
			Expect(err.Error()).To(ContainSubstring("to JSON has error"))
		})
	})
})

var _ = Describe("Send request and check status", func() {
	var client *gt.Client

	AfterEach(func() {
		gockConfig.Off()
		client = nil
	})

	Context("Normal response", func() {
		BeforeEach(func() {
			gockConfig.New().Get("/some1").
				Reply(http.StatusOK).
				Type("text").
				BodyString("Hello Man !!")

			client = gockConfig.GentlemanT.SetupClient(
				tl.CommonGentleman.NewDefaultClient(),
			).Path("/some1")
		})

		It("Equals to 200", func() {
			resp, err := tl.ToGentlemanReq(client.Get()).SendAndStatusMatch(http.StatusOK)

			Expect(err).To(Succeed())
			Expect(resp.String()).To(Equal("Hello Man !!"))
		})

		It("Not Equals to 404", func() {
			resp, err := tl.ToGentlemanReq(client.Get()).SendAndStatusMatch(http.StatusNotFound)

			Expect(err).To(HaveOccurred())
			GinkgoT().Logf("Error: %v", je.Details(err))
			Expect(resp).To(BeNil())
		})
	})

	Context("Request Error", func() {
		BeforeEach(func() {
			gockConfig.New().Get("/some1").
				ReplyError(errors.New("Sampele Error 1"))

			client = gockConfig.GentlemanT.SetupClient(
				tl.CommonGentleman.NewDefaultClient(),
			).Path("/some1")
		})

		It("The error message contains \"Send()\"", func() {
			resp, err := tl.ToGentlemanReq(client.Get()).SendAndStatusMatch(http.StatusNotFound)

			Expect(err).To(HaveOccurred())
			GinkgoT().Logf("Error: %v", je.Details(err))
			Expect(err.Error()).To(ContainSubstring("Send()"))
			Expect(resp).To(BeNil())
		})
	})

	Context("Timeout", func() {
		BeforeEach(func() {
			gockConfig.New().Get("/some1").
				ReplyFunc(func(resp *gock.Response) {
					time.Sleep(250 * time.Millisecond)
				})

			client = gockConfig.GentlemanT.SetupClient(
				tl.CommonGentleman.NewDefaultClient(),
			).
				Use(timeout.Request(1 * time.Millisecond)).
				Path("/some1")
		})

		It("The error message contains \"timetout\"", func() {
			resp, err := tl.ToGentlemanReq(client.Get()).SendAndStatusMatch(http.StatusNotFound)

			Expect(err).To(HaveOccurred())
			GinkgoT().Logf("Error: %v", je.Details(err))
			Expect(err.Error()).To(ContainSubstring("timeout"))
			Expect(resp).To(BeNil())
		})
	})
})

func newClientByGock() *gt.Client {
	return gt.New().Use(gp.GockPlugin)
}
