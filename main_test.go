package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// Launch test suite
func TestPkg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "main package test suite")
}

// Start app server and wait for initialization to complete.
var _ = BeforeSuite(func() {

	// Start app server
	go main()

	// Wait  for app server to finish starting
	waitForInitialization()
})

// Stop app server.
var _ = AfterSuite(func() {

	// Shutdown app server after test completion
	shutdown()
})

// BDD test cases
var _ = Describe("main package", func() {

	Context("launches activity tracker", func() {

		// Initialize http client since POST operation isn't available in http helper function
		client := &http.Client{}

		It("should have initialized the application server", func() {
			Expect(appServer).ShouldNot(BeNil())
			Expect(isInitializationComplete).Should(BeTrue())
		})

		It("should report zero activity summary with no prior traffic", func() {

			resp, err := http.Get("http://" + appServerListenAddress + "/metric/test/sum")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).ShouldNot(BeNil())

			// Suppressing Close potential error since operation is only contained within the test
			//goland:noinspection GoUnhandledErrorResult
			defer resp.Body.Close()
			Expect(resp.StatusCode).Should(Equal(fiber.StatusOK))

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())

			activity := activityMetric{}
			err = json.Unmarshal(bodyBytes, &activity)

			Expect(err).NotTo(HaveOccurred())
			Expect(activity.Value).Should(Equal(0))
		})

		It("should report activity summary after receiving traffic", func() {

			bodyJSONReader := strings.NewReader("{\"value\":30}")
			req, err := http.NewRequest(http.MethodPost, "http://"+appServerListenAddress+"/metric/test", bodyJSONReader)
			Expect(err).NotTo(HaveOccurred())
			Expect(req).ShouldNot(BeNil())

			req.Header.Set("Content-Type", "application/json")
			postResp, err := client.Do(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(postResp).ShouldNot(BeNil())

			// Suppressing Close potential error since operation is only contained within the test
			//goland:noinspection GoUnhandledErrorResult
			defer postResp.Body.Close()

			Expect(postResp.StatusCode).Should(Equal(fiber.StatusOK))

			resp, err := http.Get("http://" + appServerListenAddress + "/metric/test/sum")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).ShouldNot(BeNil())

			// Suppressing Close potential error since operation is only contained within the test
			//goland:noinspection GoUnhandledErrorResult
			defer resp.Body.Close()
			Expect(resp.StatusCode).Should(Equal(fiber.StatusOK))

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())

			activity := activityMetric{}
			err = json.Unmarshal(bodyBytes, &activity)

			Expect(err).NotTo(HaveOccurred())
			Expect(activity.Value).Should(Equal(30))
		})

		It("should append activity when using the same activity key", func() {

			bodyJSONReader := strings.NewReader("{\"value\":15}")
			req, err := http.NewRequest(http.MethodPost, "http://"+appServerListenAddress+"/metric/test", bodyJSONReader)
			Expect(err).NotTo(HaveOccurred())
			Expect(req).ShouldNot(BeNil())

			req.Header.Set("Content-Type", "application/json")
			postResp, err := client.Do(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(postResp).ShouldNot(BeNil())

			// Suppressing Close potential error since operation is only contained within the test
			//goland:noinspection GoUnhandledErrorResult
			defer postResp.Body.Close()

			Expect(postResp.StatusCode).Should(Equal(fiber.StatusOK))

			resp, err := http.Get("http://" + appServerListenAddress + "/metric/test/sum")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).ShouldNot(BeNil())

			// Suppressing Close potential error since operation is only contained within the test
			//goland:noinspection GoUnhandledErrorResult
			defer resp.Body.Close()
			Expect(resp.StatusCode).Should(Equal(fiber.StatusOK))

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())

			activity := activityMetric{}
			err = json.Unmarshal(bodyBytes, &activity)

			Expect(err).NotTo(HaveOccurred())
			Expect(activity.Value).Should(Equal(45))
		})
	})
})
