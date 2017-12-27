package vox

import (
	"net/http/httptest"
	"testing"
)

func TestNewResponse(t *testing.T) {
	w := httptest.NewRecorder()
	response := createResponse(w)
	response.setImplict()
	if response.Status != 404 {
		t.Fail()
	}

	w = httptest.NewRecorder()
	response = createResponse(w)
	response.Body = "plaintext"
	response.setImplict()
	if response.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Fail()
	}

	w = httptest.NewRecorder()
	response = createResponse(w)
	response.Body = `
	<!doctype html
	`
	response.setImplict()
	if response.Header.Get("Content-Type") != "text/html; charset=utf-8" {
		t.Fail()
	}

	w = httptest.NewRecorder()
	response = createResponse(w)
	response.Body = map[string]string{"foo": "bar"}
	response.setImplict()
	if response.Header.Get("Content-Type") != "application/json" {
		t.Fail()
	}
}
