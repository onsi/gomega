package ghttp_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
	"io"
	"net/http"
)

var _ = Describe("Handlers test", func() {
	Describe("RespondWithMultiple", func() {
		value := -1
		handlers := RespondWithMultiple(
			func(resp http.ResponseWriter, req *http.Request) { value = 1 },
			func(resp http.ResponseWriter, req *http.Request) { value = 2 },
			func(resp http.ResponseWriter, req *http.Request) { value = 3 },
			func(resp http.ResponseWriter, req *http.Request) { value = 4 },
			func(resp http.ResponseWriter, req *http.Request) { value = 5 },
		)
		DescribeTable("Should rotate through each handler and repeat last handler",
			func(expected int) {
				handlers(nil, nil)
				Expect(value).To(Equal(expected))
			},
			Entry("call 1", 1),
			Entry("call 2", 2),
			Entry("call 3", 3),
			Entry("call 4", 4),
			Entry("call 5", 5),
			Entry("call 6", 5),
		)
	})

	Describe("RoundRobinWithMultiple", func() {
		value := -1
		handlers := RoundRobinWithMultiple(
			func(resp http.ResponseWriter, req *http.Request) { value = 1 },
			func(resp http.ResponseWriter, req *http.Request) { value = 2 },
			func(resp http.ResponseWriter, req *http.Request) { value = 3 },
			func(resp http.ResponseWriter, req *http.Request) { value = 4 },
			func(resp http.ResponseWriter, req *http.Request) { value = 5 },
		)
		DescribeTable("Should rotate through each and start at the beginning again",
			func(expected int) {
				handlers(nil, nil)
				Expect(value).To(Equal(expected))
			},
			Entry("call 1", 1),
			Entry("call 2", 2),
			Entry("call 3", 3),
			Entry("call 4", 4),
			Entry("call 5", 5),
			Entry("call 6", 1),
		)
	})
})

func ExampleRespondWithMultiple() {
	server := NewServer()
	server.RouteToHandler(http.MethodGet, "/example",
		RespondWithMultiple(
			RespondWith(http.StatusOK, "1"),
			RespondWith(http.StatusOK, "2"),
			RespondWith(http.StatusNotFound, "Not found"),
			RespondWith(http.StatusOK, "3"),
			RespondWith(http.StatusInternalServerError, "Internal server error"),
		))
	for i := 0; i < 6; i++ {
		resp, _ := http.Get(server.URL() + "/example")
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("%d - %s\n", resp.StatusCode, string(body))
	}

	// Output: 200 - 1
	// 200 - 2
	// 404 - Not found
	// 200 - 3
	// 500 - Internal server error
	// 500 - Internal server error
}

func ExampleRoundRobinWithMultiple() {
	server := NewServer()
	server.RouteToHandler(http.MethodGet, "/example",
		RoundRobinWithMultiple(
			RespondWith(http.StatusOK, "1"),
			RespondWith(http.StatusOK, "2"),
			RespondWith(http.StatusNotFound, "Not found"),
			RespondWith(http.StatusOK, "3"),
			RespondWith(http.StatusInternalServerError, "Internal server error"),
		))
	for i := 0; i < 6; i++ {
		resp, _ := http.Get(server.URL() + "/example")
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("%d - %s\n", resp.StatusCode, string(body))
	}

	// Output: 200 - 1
	// 200 - 2
	// 404 - Not found
	// 200 - 3
	// 500 - Internal server error
	// 200 - 1
}
