package integration_test

import (
	"io"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Calendar API", func() {
	BeforeEach(func() {
		// Ждем запуска всех сервисов
		time.Sleep(5 * time.Second)
	})

	It("should respond with pong to ping endpoint", func() {
		resp, err := http.Get("http://calendar:8080/ping")
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(body)).To(Equal("pong"))
	})
})
