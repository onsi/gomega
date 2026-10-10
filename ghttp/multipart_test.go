package ghttp_test

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

type multipartTestBody struct {
	io.Reader
	closed   bool
	closeErr error
}

func (b *multipartTestBody) Close() error { b.closed = true; return b.closeErr }

func multipartTestRequest(t *testing.T) (*http.Request, *multipartTestBody) {
	t.Helper()
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	for _, value := range []string{"first", "second"} {
		if err := w.WriteField("label", value); err != nil {
			t.Fatal(err)
		}
	}
	for _, name := range []string{"one.txt", "two.txt"} {
		part, err := w.CreateFormFile("upload", name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := io.WriteString(part, "contents of "+name); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	tracked := &multipartTestBody{Reader: bytes.NewReader(body.Bytes())}
	req := httptest.NewRequest(http.MethodPost, "/?label=query&invalid=a;b", tracked)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req, tracked
}

func TestVerifyMultipartForm(t *testing.T) {
	for _, parsed := range []bool{false, true} {
		t.Run(map[bool]string{false: "fresh", true: "already-parsed"}[parsed], func(t *testing.T) {
			temp := t.TempDir()
			t.Setenv("TMPDIR", temp)
			req, body := multipartTestRequest(t)
			if parsed {
				req.Form = url.Values{}
				if err := req.ParseMultipartForm(0); err != nil {
					t.Fatal(err)
				}
			}
			var failures []string
			g := gomega.NewGomega(func(message string, _ ...int) {
				failures = append(failures, message)
			})
			ghttp.NewGHTTPWithGomega(g).VerifyMultipartForm(
				url.Values{"label": {"first", "second"}},
				map[string][]ghttp.MultipartFile{"upload": {
					{Filename: "one.txt", Content: []byte("contents of one.txt"), Header: textproto.MIMEHeader{"content-type": {"application/octet-stream"}}},
					{Filename: "two.txt", Content: []byte("contents of two.txt")},
				}}, ghttp.MultipartFormOptions{Exact: true},
			)(httptest.NewRecorder(), req)
			if len(failures) != 0 {
				t.Fatalf("unexpected failures: %v", failures)
			}
			if !body.closed {
				t.Fatal("body was not closed")
			}
			entries, err := os.ReadDir(temp)
			if err != nil || len(entries) != 0 {
				t.Fatalf("temporary files remain: %v, %v", entries, err)
			}
		})
	}
}

func TestVerifyMultipartFormFailures(t *testing.T) {
	for _, test := range []struct {
		name  string
		alter func(*http.Request, *multipartTestBody)
		want  string
	}{
		{"content-type", func(r *http.Request, _ *multipartTestBody) { r.Header.Set("Content-Type", "text/plain") }, "content type"},
		{"missing-boundary", func(r *http.Request, _ *multipartTestBody) { r.Header.Set("Content-Type", "multipart/form-data") }, "boundary"},
		{"truncated-body", func(_ *http.Request, b *multipartTestBody) { b.Reader = strings.NewReader("invalid multipart") }, "multipart parse"},
		{"close-error", func(_ *http.Request, b *multipartTestBody) { b.closeErr = errors.New("close failed") }, "close failed"},
	} {
		t.Run(test.name, func(t *testing.T) {
			req, body := multipartTestRequest(t)
			test.alter(req, body)
			var failures []string
			g := gomega.NewGomega(func(message string, _ ...int) { failures = append(failures, message) })
			ghttp.NewGHTTPWithGomega(g).VerifyMultipartForm(nil, nil, ghttp.MultipartFormOptions{})(httptest.NewRecorder(), req)
			if !strings.Contains(strings.Join(failures, "\n"), test.want) {
				t.Fatalf("missing %q in %v", test.want, failures)
			}
			if !body.closed {
				t.Fatal("body not closed on failure")
			}
		})
	}
}

func TestVerifyMultipartFormMismatchAndExact(t *testing.T) {
	for _, exact := range []bool{false, true} {
		req, body := multipartTestRequest(t)
		var failures []string
		g := gomega.NewGomega(func(message string, _ ...int) {
			if !body.closed {
				t.Error("assertion reported before cleanup")
			}
			failures = append(failures, message)
		})
		ghttp.NewGHTTPWithGomega(g).VerifyMultipartForm(nil, nil, ghttp.MultipartFormOptions{Exact: exact})(httptest.NewRecorder(), req)
		if (len(failures) > 0) != exact {
			t.Fatalf("exact=%v, failures=%v", exact, failures)
		}
	}
	req, _ := multipartTestRequest(t)
	var failures []string
	g := gomega.NewGomega(func(message string, _ ...int) { failures = append(failures, message) })
	ghttp.NewGHTTPWithGomega(g).VerifyMultipartForm(url.Values{"label": {"second", "first"}}, map[string][]ghttp.MultipartFile{
		"upload": {{Filename: "wrong.txt", Content: []byte("wrong")}, {Filename: "two.txt", Content: []byte("contents of two.txt")}},
	}, ghttp.MultipartFormOptions{})(httptest.NewRecorder(), req)
	message := strings.Join(failures, "\n")
	for _, want := range []string{"filename", "content", "value"} {
		if !strings.Contains(message, want) {
			t.Errorf("missing %s failure in %s", want, message)
		}
	}
}
