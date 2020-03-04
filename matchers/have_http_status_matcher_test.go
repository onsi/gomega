package matchers_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HaveHTTPStatus", func() {
	When("ACTUAL is *http.Response", func() {
		When("EXPECTED is integer", func() {
			It("matches the StatusCode", func() {
				resp := &http.Response{StatusCode: http.StatusOK}
				Expect(resp).To(HaveHTTPStatus(http.StatusOK))
				Expect(resp).NotTo(HaveHTTPStatus(http.StatusNotFound))
			})
		})
		When("EXPECTED is string", func() {
			It("matches the Status", func() {
				resp := &http.Response{Status: "200 OK"}
				Expect(resp).To(HaveHTTPStatus("200 OK"))
				Expect(resp).NotTo(HaveHTTPStatus("404 Not Found"))
			})
		})
		When("EXPECTED is anything else", func() {
			It("does not match", func() {
				failures := InterceptGomegaFailures(func() {
					resp := &http.Response{StatusCode: http.StatusOK}
					Expect(resp).NotTo(HaveHTTPStatus(true))
				})
				Expect(failures).To(ConsistOf("HaveHTTPStatus matcher must be passed an int or a string. Got:\n    <bool>: true"))
			})
		})
	})

	When("ACTUAL is *httptest.ResponseRecorder", func() {
		When("EXPECTED is integer", func() {
			It("matches the StatusCode", func() {
				resp := &httptest.ResponseRecorder{Code: http.StatusOK}
				Expect(resp).To(HaveHTTPStatus(http.StatusOK))
				Expect(resp).NotTo(HaveHTTPStatus(http.StatusNotFound))
			})
		})
		When("EXPECTED is string", func() {
			It("matches the Status", func() {
				resp := &httptest.ResponseRecorder{Code: http.StatusOK}
				Expect(resp).To(HaveHTTPStatus("200 OK"))
				Expect(resp).NotTo(HaveHTTPStatus("404 Not Found"))
			})
		})
		When("EXPECTED is anything else", func() {
			It("does not match", func() {
				failures := InterceptGomegaFailures(func() {
					resp := &httptest.ResponseRecorder{Code: http.StatusOK}
					Expect(resp).NotTo(HaveHTTPStatus(nil))
				})
				Expect(failures).To(ConsistOf("HaveHTTPStatus matcher must be passed an int or a string. Got:\n    <nil>: nil"))
			})
		})
	})

	When("ACTUAL is neither *http.Response nor *httptest.ResponseRecorder", func() {
		It("errors", func() {
			failures := InterceptGomegaFailures(func() {
				Expect("foo").To(HaveHTTPStatus(http.StatusOK))
			})
			Expect(failures).To(ConsistOf("HaveHTTPStatus matcher expects *http.Response or *httptest.ResponseRecorder. Got:\n    <string>: foo"))
		})
	})

	Describe("FailureMessage", func() {
		It("returns message", func() {
			failures := InterceptGomegaFailures(func() {
				resp := &http.Response{StatusCode: http.StatusBadGateway}
				Expect(resp).To(HaveHTTPStatus(http.StatusOK))
			})
			Expect(failures).To(ConsistOf(MatchRegexp("Expected(.|\n)*StatusCode: 502(.|\n)*to have HTTP status\n    <int>: 200")))
		})
	})
	Describe("NegatedFailureMessage", func() {
		It("returns message", func() {
			failures := InterceptGomegaFailures(func() {
				resp := &http.Response{StatusCode: http.StatusOK}
				Expect(resp).NotTo(HaveHTTPStatus(http.StatusOK))
			})
			Expect(failures).To(ConsistOf(MatchRegexp("Expected(.|\n)*StatusCode: 200(.|\n)*not to have HTTP status\n    <int>: 200")))
		})
	})
})
