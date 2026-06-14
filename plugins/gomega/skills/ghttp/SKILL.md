---
name: ghttp
description: The ghttp test HTTP server for testing HTTP clients — NewServer/NewTLSServer, AppendHandlers, the Verify* assertions (VerifyRequest/VerifyHeader/VerifyHeaderKV/VerifyJSON/VerifyForm/VerifyBasicAuth/VerifyContentType/VerifyBody), RespondWith/RespondWithJSONEncoded(Ptr), CombineHandlers, RouteToHandler for unordered MUXed routes, AllowUnhandledRequests, RoundTripper, and TLS. Use when testing code that makes outbound HTTP requests.
---

# ghttp: testing HTTP clients

`ghttp` spins up a real test HTTP server so you can exercise code that makes outbound HTTP requests (clients, SDKs, anything calling `net/http`). You register handlers that **assert on the incoming request and return a canned response**; the server records every request it receives. It wraps `net/http/httptest`.

```go
import (
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"net/http"
)
```

This skill assumes a dot-import of gomega. Docs: <https://onsi.github.io/gomega/#ghttp-testing-http-clients>. See `gomega:matchers` for the matchers used in assertions and `gomega:async` if the client is non-blocking.

## Lifecycle

```go
server := ghttp.NewServer()          // or ghttp.NewTLSServer()
defer server.Close()                 // ALWAYS defer Close

client := myapi.NewClient(server.URL())   // URL is random per run — inject it
```

**The server URL is auto-generated and varies between runs**, so you must inject `server.URL()` into the code under test rather than hard-coding it. With Ginkgo, `server = ghttp.NewServer()` in `BeforeEach` and `server.Close()` in `AfterEach`.

## The handler model

`server.AppendHandlers(handlers...)` registers an **ordered** list of handlers, **one per expected request**. The first request is matched against the first handler, the second against the second, and so on. A handler is just an `http.HandlerFunc`.

```go
server.AppendHandlers(
	ghttp.VerifyRequest("GET", "/sprockets"),
)
```

**Gotchas (the core mental model):**
- **AppendHandlers is strictly ordered.** Requests must arrive in exactly the registered order.
- **Each handler handles exactly one request.** To run several checks against a *single* request, combine them with `ghttp.CombineHandlers` (below).
- **The request count is asserted by default.** If the client sends more requests than there are handlers (or hits an unregistered route), the test fails. Guard the trivial false-positive (client made no request) with `Expect(server.ReceivedRequests()).To(HaveLen(1))`.

## Verifying the incoming request

These return handlers that fail the test if the request doesn't match:

```go
ghttp.VerifyRequest("GET", "/sprockets")                       // method + path
ghttp.VerifyRequest("GET", "/sprockets", "category=encabulators") // + RawQuery
ghttp.VerifyHeader(http.Header{"X-API-Version": []string{"1.0"}})
ghttp.VerifyHeaderKV("X-API-Version", "1.0")                   // convenience for one key
ghttp.VerifyContentType("application/json")
ghttp.VerifyBasicAuth("gopher", "tacoshell")
ghttp.VerifyBody([]byte("raw body"))
ghttp.VerifyJSON(`{"name":"foo"}`)                             // semantic JSON compare
ghttp.VerifyJSONRepresenting(Sprocket{Name: "foo"})           // marshal obj, then compare
ghttp.VerifyForm(url.Values{"q": []string{"x"}})              // parsed form values
ghttp.VerifyFormKV("q", "x")                                  // convenience for one key
```

The `path` argument to `VerifyRequest` accepts a string or a `*regexp.Regexp`. You can write **your own** verify handler — pass any `http.HandlerFunc` that makes assertions against the `*http.Request`.

## Combining handlers for one request

`ghttp.CombineHandlers(...)` wraps several handlers into **one** handler so they all run against a single request. Pass that one wrapper to `AppendHandlers`. This is how you assert on a request *and* respond:

```go
server.AppendHandlers(
	ghttp.CombineHandlers(
		ghttp.VerifyRequest("GET", "/sprockets", "category=encabulators"),
		ghttp.VerifyBasicAuth(username, password),
		ghttp.VerifyHeaderKV("X-API-Version", "1.0"),
		ghttp.RespondWithJSONEncoded(http.StatusOK, sprockets),
	),
)
```

## Providing responses

```go
ghttp.RespondWith(http.StatusOK, `[{"name":"foo"}]`)          // status + body (+ optional http.Header)
ghttp.RespondWith(http.StatusOK, body, http.Header{"X-Foo": []string{"bar"}})
ghttp.RespondWithJSONEncoded(http.StatusOK, anyObject)        // JSON-encodes the object for you
ghttp.RespondWithJSONEncodedPtr(&statusCode, &object)         // reads through pointers at request time
ghttp.RespondWithPtr(&statusCode, &body)
```

## End-to-end example

```go
var server *ghttp.Server
var client *sprockets.Client

BeforeEach(func() {
	server = ghttp.NewServer()
	client = sprockets.NewClient(server.URL())
})
AfterEach(func() { server.Close() })

It("fetches sprockets", func() {
	returned := []Sprocket{{Name: "entropic decoupler", Color: "red"}}
	server.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/sprockets", "category=encabulators"),
			ghttp.VerifyHeaderKV("X-API-Version", "1.0"),
			ghttp.RespondWithJSONEncoded(http.StatusOK, returned),
		),
	)

	sprockets, err := client.FetchSprockets("encabulators")
	Expect(err).NotTo(HaveOccurred())
	Expect(sprockets).To(Equal(returned))
	Expect(server.ReceivedRequests()).To(HaveLen(1))
})
```

## Testing different response scenarios

To exercise 200/401/404 against the *same* handler setup, capture a status (and/or body) variable in a closure via the `*Ptr` responders and mutate it per sub-context. The pointer is dereferenced when the request arrives:

```go
var statusCode int
var returned []Sprocket

BeforeEach(func() {
	server.AppendHandlers(ghttp.CombineHandlers(
		ghttp.VerifyRequest("GET", "/sprockets"),
		ghttp.RespondWithJSONEncodedPtr(&statusCode, &returned),
	))
})

Context("when unauthorized", func() {
	BeforeEach(func() { statusCode = http.StatusUnauthorized })
	It("errors", func() {
		_, err := client.FetchSprockets("")
		Expect(err).To(MatchError(SprocketsErrorUnauthorized))
	})
})
```

## Handling multiple requests

Pass multiple handlers to `AppendHandlers` for clients that make several requests (e.g. pagination). They're consumed in order:

```go
server.AppendHandlers(
	ghttp.CombineHandlers(
		ghttp.VerifyRequest("GET", "/sprockets", "category=encabulators"),
		ghttp.RespondWithJSONEncoded(http.StatusOK, firstPage),
	),
	ghttp.CombineHandlers(
		ghttp.VerifyRequest("GET", "/sprockets", "category=encabulators&pagination-token=next"),
		ghttp.RespondWithJSONEncoded(http.StatusOK, secondPage),
	),
)
```

Because the server fails when requests exceed handlers, this also asserts the client *stops* after the last page.

### Mutating registered handlers by index

```go
server.SetHandler(0, newHandler)            // replace handler at index
h := server.GetHandler(0)                   // read it back
server.WrapHandler(0, extraVerifyHandler)   // run an extra handler around the existing one
```

## MUXing unordered routes

`AppendHandlers` is for ordered assertions. When requests can arrive in any order or any number of times, use `server.RouteToHandler(method, path, handler)` instead. The `path` is a string or a `*regexp.Regexp`.

```go
server.RouteToHandler("POST", "/login", ghttp.CombineHandlers(
	ghttp.VerifyRequest("POST", "/login", "user=bob&password=password"),
	ghttp.RespondWith(http.StatusOK, "your-auth-token"),
))
```

**Routing precedence:** on each request the server first checks `RouteToHandler` routes; if none match it pops the next `AppendHandlers` handler; if that stack is empty it consults the unhandled-request policy below. You can mix both: ordered assertions via `AppendHandlers` for the requests you care about, MUXed routes for incidental traffic (auth, health checks) that interleaves.

## Allowing unhandled requests

By default an unhandled request fails the test. To return a fixed status instead:

```go
server.SetAllowUnhandledRequests(true)
server.SetUnhandledRequestStatusCode(http.StatusOK) // default would be 500
```

Unhandled requests are still recorded in `server.ReceivedRequests()`, so you can assert on them after the fact.

## Escape hatches

- `server.HTTPTestServer` is the underlying `*httptest.Server` — use it for custom listeners (via `ghttp.NewUnstartedServer()` + `server.Start()`), TLS config, or anything ghttp doesn't wrap.
- `server.RoundTripper(rt)` returns an `http.RoundTripper` that redirects requests to the test server. Use it when you must hand a configured `*http.Client` to the code under test rather than a URL. Pass `nil` to use `http.DefaultTransport`:

```go
httpClient := &http.Client{Transport: server.RoundTripper(nil)}
```
