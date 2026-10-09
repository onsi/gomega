package ghttp

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"slices"
	"sort"
	"strings"

	"github.com/onsi/gomega"
)

// MultipartFile describes one uploaded file, in its field's wire order.
// Header selects headers to check; additional headers are allowed.
type MultipartFile struct {
	Filename string
	Content  []byte
	Header   textproto.MIMEHeader
}

// MultipartFormOptions controls parsing and whether additional fields are allowed.
type MultipartFormOptions struct {
	// MaxMemory is passed to net/http.Request.ParseMultipartForm. Larger files spill to disk.
	MaxMemory int64
	// Exact rejects value and file fields not included in the expectations.
	Exact bool
}

// VerifyMultipartForm verifies multipart body values and uploaded files.
// Repeated values and files must match in order. Query parameters are ignored.
// The handler closes the body and removes temporary files before reporting failures.
// Already parsed forms are reused and their temporary files are also removed.
func (g GHTTPWithGomega) VerifyMultipartForm(values url.Values, files map[string][]MultipartFile, options MultipartFormOptions) http.HandlerFunc {
	return func(_ http.ResponseWriter, r *http.Request) {
		request := *r
		request.Form = url.Values{}
		r = &request
		var problems []error
		func() {
			defer func() {
				if r.MultipartForm != nil {
					problems = append(problems, removeMultipartFiles(r.MultipartForm)...)
				}
				if r.Body != nil {
					if err := r.Body.Close(); err != nil {
						problems = append(problems, fmt.Errorf("multipart body close: %w", err))
					}
				}
			}()
			if options.MaxMemory < 0 {
				problems = append(problems, fmt.Errorf("multipart configuration: negative memory threshold"))
				return
			}
			mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
			if err != nil || mediaType != "multipart/form-data" {
				problems = append(problems, fmt.Errorf("expected multipart/form-data content type"))
				return
			}
			if r.Body == nil {
				problems = append(problems, fmt.Errorf("multipart parse: missing body"))
				return
			}
			if err := r.ParseMultipartForm(options.MaxMemory); err != nil {
				problems = append(problems, fmt.Errorf("multipart parse: %w", err))
				return
			}
			problems = append(problems, verifyMultipartFiles(r.MultipartForm.File, files)...)
			problems = append(problems, verifyMultipartValues(r.MultipartForm.Value, values)...)
			if options.Exact {
				problems = append(problems, unexpectedMultipartKeys(r.MultipartForm.Value, values, "value")...)
				problems = append(problems, unexpectedMultipartKeys(r.MultipartForm.File, files, "file")...)
			}
		}()
		for _, problem := range problems {
			g.gomega.Expect(problem).ShouldNot(gomega.HaveOccurred())
		}
	}
}

// VerifyMultipartForm uses the default Gomega instance.
func VerifyMultipartForm(values url.Values, files map[string][]MultipartFile, options MultipartFormOptions) http.HandlerFunc {
	return NewGHTTPWithGomega(gomega.Default).VerifyMultipartForm(values, files, options)
}

func multipartKeys[V any](values map[string]V) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func removeMultipartFiles(form *multipart.Form) []error {
	var problems []error
	for _, key := range multipartKeys(form.File) {
		for index, header := range form.File[key] {
			if header == nil {
				continue
			}
			part := multipart.Form{File: map[string][]*multipart.FileHeader{key: {header}}}
			if err := part.RemoveAll(); err != nil {
				problems = append(problems, fmt.Errorf("multipart file %q occurrence %d cleanup: %w", key, index, err))
			}
		}
	}
	return problems
}

func unexpectedMultipartKeys[A, E any](actual map[string]A, expected map[string]E, kind string) []error {
	var problems []error
	for _, key := range multipartKeys(actual) {
		if _, found := expected[key]; !found {
			problems = append(problems, fmt.Errorf("multipart unexpected %s key %q", kind, key))
		}
	}
	return problems
}

func verifyMultipartValues(actual map[string][]string, expected url.Values) []error {
	var problems []error
	for _, key := range multipartKeys(expected) {
		got, want := actual[key], expected[key]
		if len(got) != len(want) {
			problems = append(problems, fmt.Errorf("multipart value %q count: got %d, expected %d", key, len(got), len(want)))
			continue
		}
		for index, value := range want {
			if got[index] != value {
				problems = append(problems, fmt.Errorf("multipart value %q occurrence %d: got %q, expected %q", key, index, got[index], value))
			}
		}
	}
	return problems
}

func verifyMultipartFiles(actual map[string][]*multipart.FileHeader, expected map[string][]MultipartFile) []error {
	var problems []error
	for _, key := range multipartKeys(expected) {
		got, want := actual[key], expected[key]
		if len(got) != len(want) {
			problems = append(problems, fmt.Errorf("multipart file %q count: got %d, expected %d", key, len(got), len(want)))
			continue
		}
		for index, file := range want {
			location := fmt.Sprintf("multipart file %q occurrence %d", key, index)
			problems = append(problems, verifyMultipartFile(got[index], file, location)...)
		}
	}
	return problems
}

func verifyMultipartFile(actual *multipart.FileHeader, expected MultipartFile, location string) []error {
	if actual == nil {
		return []error{fmt.Errorf("%s: missing file header", location)}
	}
	problems := verifyMultipartFileMetadata(actual, expected, location)
	file, err := actual.Open()
	if err != nil {
		return append(problems, fmt.Errorf("%s open: %w", location, err))
	}
	contentErr := compareMultipartContent(file, expected.Content)
	closeErr := file.Close()
	if contentErr != nil {
		problems = append(problems, fmt.Errorf("%s content: %w", location, contentErr))
	}
	if closeErr != nil {
		problems = append(problems, fmt.Errorf("%s close: %w", location, closeErr))
	}
	return problems
}

func verifyMultipartFileMetadata(actual *multipart.FileHeader, expected MultipartFile, location string) []error {
	var problems []error
	if actual.Filename != expected.Filename {
		problems = append(problems, fmt.Errorf("%s filename: got %q, expected %q", location, actual.Filename, expected.Filename))
	}
	for _, key := range multipartKeys(expected.Header) {
		canonical := textproto.CanonicalMIMEHeaderKey(key)
		got := multipartHeaderValues(actual.Header, canonical)
		if !slices.Equal(got, expected.Header[key]) {
			problems = append(problems, fmt.Errorf("%s header %q: got %q, expected %q", location, canonical, got, expected.Header[key]))
		}
	}
	return problems
}

func multipartHeaderValues(header textproto.MIMEHeader, name string) []string {
	var values []string
	for _, key := range multipartKeys(header) {
		if strings.EqualFold(key, name) {
			values = append(values, header[key]...)
		}
	}
	return values
}

// Compare incrementally so a large unexpected upload does not require a second
// copy in memory. Continue to EOF to distinguish read errors from mismatches.
func compareMultipartContent(reader io.Reader, expected []byte) error {
	var buffer [32 * 1024]byte
	offset := 0
	difference := multipartContentDifference{offset: -1, expectedSize: len(expected), expectedByte: -1, actualByte: -1}
	for {
		n, err := reader.Read(buffer[:])
		if n > 0 {
			remaining := len(expected) - offset
			compare := min(n, remaining)
			if difference.offset < 0 {
				for index := 0; index < compare; index++ {
					if buffer[index] != expected[offset+index] {
						difference.offset = difference.actualSize + int64(index)
						difference.expectedByte = int(expected[offset+index])
						difference.actualByte = int(buffer[index])
						break
					}
				}
				if difference.offset < 0 && n > remaining {
					difference.offset = int64(len(expected))
					difference.actualByte = int(buffer[compare])
				}
			}
			offset += compare
			difference.actualSize += int64(n)
		}
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("read: %w", err)
			}
			break
		}
	}
	if difference.offset < 0 && offset < len(expected) {
		difference.offset = difference.actualSize
		difference.expectedByte = int(expected[offset])
	}
	if difference.offset >= 0 {
		return difference
	}
	return nil
}

type multipartContentDifference struct {
	offset       int64
	expectedSize int
	actualSize   int64
	expectedByte int
	actualByte   int
}

func (difference multipartContentDifference) Error() string {
	return fmt.Sprintf("offset=%d expected-byte=%s actual-byte=%s expected-length=%d actual-length=%d",
		difference.offset,
		multipartByteDescription(difference.expectedByte),
		multipartByteDescription(difference.actualByte),
		difference.expectedSize,
		difference.actualSize)
}

func multipartByteDescription(value int) string {
	if value < 0 {
		return "EOF"
	}
	return fmt.Sprintf("0x%02x", value)
}
