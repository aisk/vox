package vox

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
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
	app.Use(func(ctx *Context, req *Request, res *Response) {
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
	app.Use(func(ctx *Context, req *Request, res *Response) {
		res.SetCookie(&http.Cookie{Name: "foo", Value: "bar"})
	})
	r := httptest.NewRequest("GET", "http://test.com/", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.HeaderMap.Get("Set-Cookie") != "foo=bar" {
		t.Fatal()
	}
}

func TestResponseReader(t *testing.T) {
	app := New()
	app.Use(func(ctx *Context, req *Request, res *Response) {
		res.Body = strings.NewReader("Hello io.Reader!")
	})
	r := httptest.NewRequest("GET", "http://test.com/", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	body, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		t.Fail()
	}
	if string(body) != "Hello io.Reader!" {
		t.Fail()
	}
}

type MockReadCloser struct {
	io.Reader
	closed bool
}

func (rc *MockReadCloser) Close() error {
	rc.closed = true
	return nil
}

var _ io.ReadCloser = &MockReadCloser{}

func TestResponseReadCloser(t *testing.T) {
	app := New()
	rc := &MockReadCloser{strings.NewReader("Hello io.Reader!"), false}
	app.Use(func(ctx *Context, req *Request, res *Response) {
		res.Body = rc
	})
	r := httptest.NewRequest("GET", "http://test.com/", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	body, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		t.Fail()
	}
	if string(body) != "Hello io.Reader!" {
		t.Fail()
	}
	if !rc.closed {
		t.Fail()
	}
}
