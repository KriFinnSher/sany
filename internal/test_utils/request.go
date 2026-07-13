package test_utils

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MultipartRequest builds a multipart request with one optional file field.
func MultipartRequest(t testing.TB, method, path, key, name string, data []byte) *http.Request {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if key != "" {
		part, err := writer.CreateFormFile(key, name)
		if err != nil {
			t.Fatalf("create multipart field: %v", err)
		}
		if _, err := part.Write(data); err != nil {
			t.Fatalf("write multipart field: %v", err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	request := httptest.NewRequest(method, path, &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	return request
}
