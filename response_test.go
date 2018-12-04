package vox

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewResponse(t *testing.T) {
	w := httptest.NewRecorder()
	response := createResponse(w)
	response.setImplicit()
	if response.Status != 404 {
		t.Fail()
	}

	w = httptest.NewRecorder()
	response = createResponse(w)
	response.Body = map[string]string{"foo": "bar"}
	response.setImplicit()
	if response.Header.Get("Content-Type") != "application/json" {
		t.Fail()
	}
}

func TestRedirect(t *testing.T) {
	app := New()
	app.Use(func(req *Request, res *Response) {
		res.Redirect("/new_location", 302)
	})
	r := httptest.NewRequest("GET", "http://test.com/", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Result().StatusCode != 302 {
		t.Fatal()
	}
	if w.HeaderMap.Get("Location") != "/new_location" {
		t.Fatal()
	}
}

func TestSetCookie(t *testing.T) {
	app := New()
	app.Use(func(req *Request, res *Response) {
		res.SetCookie(&http.Cookie{Name: "foo", Value: "bar"})
	})
	r := httptest.NewRequest("GET", "http://test.com/", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.HeaderMap.Get("Set-Cookie") != "foo=bar" {
		t.Fatal()
	}
}
