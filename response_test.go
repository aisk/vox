package vox

import (
	"net/http/httptest"
	"testing"
)

func TestNewResponse(t *testing.T) {
	w := httptest.NewRecorder()
	response := createResponse(w)

	if response.Status() != 404 {
		t.Fail()
	}

	response.SetBody(nil)
	if response.Status() != 204 {
		t.Fail()
	}

	response.SetBody("plaintext")
	if response.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Fail()
	}

	response.SetBody(`
	<!doctype html
	`)
	if response.Header.Get("Content-Type") != "text/html; charset=utf-8" {
		t.Fail()
	}

	response.SetBody(map[string]string{"foo": "bar"})
	if response.Header.Get("Content-Type") != "application/json" {
		t.Fail()
	}
}
