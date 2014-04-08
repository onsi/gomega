package ghttp

import (
	"encoding/base64"
	"fmt"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
)

func CombineHandlers(handlers ...http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		for _, handler := range handlers {
			handler(w, req)
		}
	}
}

func VerifyRequest(method string, path string, rawQuery ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		Ω(req.Method).Should(Equal(method), "Method mismatch")
		Ω(req.URL.Path).Should(Equal(path), "Path mismatch")
		if len(rawQuery) > 0 {
			Ω(req.URL.RawQuery).Should(Equal(rawQuery[0]), "RawQuery mismatch")
		}
	}
}

func VerifyContentType(contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		Ω(req.Header.Get("Content-Type")).Should(Equal(contentType))
	}
}

func VerifyBasicAuth(username string, password string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")
		decoded, err := base64.StdEncoding.DecodeString(auth[6:])
		Ω(err).ShouldNot(HaveOccurred())

		Ω(string(decoded)).Should(Equal(fmt.Sprintf("%s:%s", username, password)), "Authorization mismatch")
	}
}

func VerifyHeader(header http.Header) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		for key, values := range header {
			key = http.CanonicalHeaderKey(key)
			Ω(req.Header[key]).Should(Equal(values), "Header mismatch for key: %s", key)
		}
	}
}

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

func RespondWith(statusCode int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(body))
	}
}

func RespondWithPtr(statusCode *int, body *string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(*statusCode)
		w.Write([]byte(*body))
	}
}
