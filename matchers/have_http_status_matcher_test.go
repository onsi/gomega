package matchers_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

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
				resp := &http.Response{
					StatusCode: http.StatusBadGateway,
					Status:     "502 Bad Gateway",
					Body:       ioutil.NopCloser(strings.NewReader("did not like it")),
				}
				Expect(resp).To(HaveHTTPStatus(http.StatusOK))
			})
			Expect(failures).To(HaveLen(1))
			Expect(failures[0]).To(Equal(`Expected
    <*http.Response>: {
        Status:     <string>: "502 Bad Gateway"
        StatusCode: <int>: 502
        Body:       <string>: "did not like it"
    }
to have HTTP status
    <int>: 200`), failures[0])
		})
	})

	Describe("NegatedFailureMessage", func() {
		It("returns message", func() {
			failures := InterceptGomegaFailures(func() {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Status:     "200 OK",
					Body:       ioutil.NopCloser(strings.NewReader("got it!")),
				}
				Expect(resp).NotTo(HaveHTTPStatus(http.StatusOK))
			})
			Expect(failures).To(HaveLen(1))
			Expect(failures[0]).To(Equal(`Expected
    <*http.Response>: {
        Status:     <string>: "200 OK"
        StatusCode: <int>: 200
        Body:       <string>: "got it!"
    }
not to have HTTP status
    <int>: 200`), failures[0])
		})
	})
})
