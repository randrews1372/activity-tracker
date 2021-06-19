package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

		It("should have initialized the application server", func() {
			Expect(appServer).ShouldNot(BeNil())
			Expect(isInitializationComplete).Should(BeTrue())
		})
	})
})
