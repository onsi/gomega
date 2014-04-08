/*

ghttp provides a fake server (simply a thin wrapper around httptest's server) that supports
registering multiple handlers.

ghttp also comes with a collection of handlers for making assertions.  You can use these building
blocks to expressively and efficiently write fake servers that assert against incoming requests.

More documentation coming soon!

*/

package ghttp

import (
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
)

func new() *Server {
	return &Server{
		AllowUnhandledRequests:     false,
		UnhandledRequestStatusCode: http.StatusInternalServerError,
		writeLock:                  &sync.Mutex{},
	}
}

func NewServer() *Server {
	s := new()
	s.HTTPTestServer = httptest.NewServer(s)
	return s
}

func NewTLSServer() *Server {
	s := new()
	s.HTTPTestServer = httptest.NewTLSServer(s)
	return s
}

type Server struct {
	HTTPTestServer *httptest.Server

	AllowUnhandledRequests     bool
	UnhandledRequestStatusCode int

	receivedRequests []*http.Request
	requestHandlers  []http.HandlerFunc

	writeLock *sync.Mutex
	calls     int
}

func (s *Server) URL() string {
	return s.HTTPTestServer.URL
}

func (s *Server) Close() {
	server := s.HTTPTestServer
	s.HTTPTestServer = nil
	server.Close()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.writeLock.Lock()
	defer s.writeLock.Unlock()
	defer func() {
		recover()
	}()

	if s.calls < len(s.requestHandlers) {
		s.requestHandlers[s.calls](w, req)
	} else {
		if s.AllowUnhandledRequests {
			ioutil.ReadAll(req.Body)
			req.Body.Close()
			w.WriteHeader(s.UnhandledRequestStatusCode)
		} else {
			Ω(req).Should(BeNil(), "Received Unhandled Request")
		}
	}
	s.receivedRequests = append(s.receivedRequests, req)
	s.calls++
}

func (s *Server) ReceivedRequests() []*http.Request {
	s.writeLock.Lock()
	defer s.writeLock.Unlock()

	return s.receivedRequests
}

func (s *Server) AppendHandlers(handlers ...http.HandlerFunc) {
	s.writeLock.Lock()
	defer s.writeLock.Unlock()

	s.requestHandlers = append(s.requestHandlers, handlers...)
}

func (s *Server) SetHandler(index int, handler http.HandlerFunc) {
	s.writeLock.Lock()
	defer s.writeLock.Unlock()

	s.requestHandlers[index] = handler
}

func (s *Server) GetHandler(index int) http.HandlerFunc {
	s.writeLock.Lock()
	defer s.writeLock.Unlock()

	return s.requestHandlers[index]
}

func (s *Server) WrapHandler(index int, handler http.HandlerFunc) {
	existingHandler := s.GetHandler(index)
	s.SetHandler(index, CombineHandlers(existingHandler, handler))
}
