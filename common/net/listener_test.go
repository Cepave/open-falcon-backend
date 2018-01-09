package net

import (
	"fmt"
	"net"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Initialize TCP listener", func() {
	Context("Normal initialization", func() {
		It("Listening on 0.0.0.0:19595", func() {
			var err error

			listener, err := InitTcpListener("0.0.0.0:19595")
			defer listener.Close()

			Expect(err).To(Succeed())

			By("[Failed] Listening on 0.0.0.0:19595 again")
			_, err = InitTcpListener("0.0.0.0:19595")
			GinkgoT().Logf("Error content: %v", err)
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("Controller for TCP listener", func() {
	Context("Accept in Loop", func() {
		listener := MustInitTcpListener("0.0.0.0:19596")
		testedCtrl := NewListenerController(listener)

		accepted := false

		BeforeEach(func() {
			go testedCtrl.AcceptLoop(func(conn net.Conn) {
				accepted = true
			})
		})
		AfterEach(func() {
			testedCtrl.Close()
		})

		It("Ensure the accepting is running(sending message)", func() {
			conn, err := net.Dial("tcp", "127.0.0.1:19596")
			defer conn.Close()

			Expect(err).To(Succeed())
			fmt.Fprintf(conn, "Hello World!!")

			Eventually(
				func() bool { return accepted },
				2*time.Second, 200*time.Millisecond,
			).Should(BeTrue())
		})
	})
})
