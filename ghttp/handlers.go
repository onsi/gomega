package ghttp

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	. "github.com/onsi/gomega"
)

//CombineHandler takes variadic list of handlers and produces one handler
//that calls each handler in order.
func CombineHandlers(handlers ...http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		for _, handler := range handlers {
			handler(w, req)
		}
	}
}

//VerifyRequest returns a handler that verifies that a request uses the specified method to connect to the specified path
//You may also pass in an optional rawQuery string which is tested against the request's `req.URL.RawQuery`
func VerifyRequest(method string, path string, rawQuery ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		Ω(req.Method).Should(Equal(method), "Method mismatch")
		Ω(req.URL.Path).Should(Equal(path), "Path mismatch")
		if len(rawQuery) > 0 {
			Ω(req.URL.RawQuery).Should(Equal(rawQuery[0]), "RawQuery mismatch")
		}
	}
}

//VerifyContentType returns a handler that verifies that a request has a Content-Type header set to the
//specified value
func VerifyContentType(contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		Ω(req.Header.Get("Content-Type")).Should(Equal(contentType))
	}
}

//VerifyBasicAuth returns a handler that verifies the request contains a BasicAuth Authorization header
//matching the passed in username and password
func VerifyBasicAuth(username string, password string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")
		decoded, err := base64.StdEncoding.DecodeString(auth[6:])
		Ω(err).ShouldNot(HaveOccurred())

		Ω(string(decoded)).Should(Equal(fmt.Sprintf("%s:%s", username, password)), "Authorization mismatch")
	}
}

//VerifyHeader returns a handler that verifies the request contains the passed in headers.
//The passed in header keys are first canonicalized via http.CanonicalHeaderKey.
//
//The request must contain *all* the passed in headers, but it is allowed to have additional headers
//beyond the passed in set.
func VerifyHeader(header http.Header) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		for key, values := range header {
			key = http.CanonicalHeaderKey(key)
			Ω(req.Header[key]).Should(Equal(values), "Header mismatch for key: %s", key)
		}
	}
}

//VerifyJSON returns a handler that verifies that the body of the request is a valid JSON representation
//matching the passed in JSON string.  It does this using Gomega's MatchJSON method
//
//VerifyJSON also verifies that the request's content type is application/json
func VerifyJSON(expectedJSON string) http.HandlerFunc {
	return CombineHandlers(
		VerifyContentType("application/json"),
		func(w http.ResponseWriter, req *http.Request) {
			body, err := ioutil.ReadAll(req.Body)
			req.Body.Close()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(body).Should(MatchJSON(expectedJSON), "JSON Mismatch")
		},
	)
}

//VerifyJSONRepresenting is similar to VerifyJSON.  Instead of taking a JSON string, however, it
//takes an arbitrary JSON-encodable object and verifies that the requests's body is a JSON representation
//that matches the object
func VerifyJSONRepresenting(object interface{}) http.HandlerFunc {
	data, err := json.Marshal(object)
	Ω(err).ShouldNot(HaveOccurred())
	return CombineHandlers(
		VerifyContentType("application/json"),
		VerifyJSON(string(data)),
	)
}

//RespondWith returns a handler that responds to a request with the specified status code and body
func RespondWith(statusCode int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(body))
	}
}

//RespondWithPtr returns a handler that responds to a request with the specified status code and body
//
//Unlike RespondWith, you pass RepondWithPtr a pointer to the status code and body allowing different tests
//to share the same setup but specify different status codes and bodies.
func RespondWithPtr(statusCode *int, body *string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(*statusCode)
		if body != nil {
			w.Write([]byte(*body))
		}
	}
}

//RespondWithJSONEncoded returns a handler that responds to a request with the specified status code and a body
//containing the JSON-encoding of the passed in object
func RespondWithJSONEncoded(statusCode int, object interface{}) http.HandlerFunc {
	data, err := json.Marshal(object)
	Ω(err).ShouldNot(HaveOccurred())
	return RespondWith(statusCode, string(data))
}

//RespondWithJSONEncodedPtr behaves like RespondWithJSONEncoded but takes a pointer
//to a status code and object.
//
//This allows different tests to share the same setup but specify different status codes and JSON-encoded
//objects.
func RespondWithJSONEncodedPtr(statusCode *int, object *interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data, err := json.Marshal(*object)
		Ω(err).ShouldNot(HaveOccurred())
		w.WriteHeader(*statusCode)
		w.Write(data)
	}
}
