package vox

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestRoute(t *testing.T) {
	app := New()
	app.SetConfig("logging:disable", "true")
	app.Route("GET", "/test_route", func(ctx *Context, req *Request, res *Response) {
		res.Body = "Hello Vox!"
		res.Header.Set("foo", "bar")
	})
	r := httptest.NewRequest("GET", "http://test.com/test_route", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Result().StatusCode != 200 {
		t.Errorf("expect StatusCode 200, got %d\r\n", w.Result().StatusCode)
	}

	r = httptest.NewRequest("GET", "http://test.com/invalid_path", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Result().StatusCode != 404 {
		t.Errorf("expect StatusCode 404, got %d\r\n", w.Result().StatusCode)
	}
}

func TestMatchAnyMethod(t *testing.T) {
	app := New()
	app.SetConfig("logging:disable", "true")
	app.Route("*", "/test_route", func(ctx *Context, req *Request, res *Response) {
		res.Body = "matched!"
		res.Status = http.StatusFound
	})
	r := httptest.NewRequest("ANYMETHOD", "http://test.com/test_route", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Result().StatusCode != http.StatusFound {
		t.Errorf("expect StatusCode 302, got %d\r\n", w.Result().StatusCode)
	}
}

func TestRouteWithParams(t *testing.T) {
	app := New()
	app.SetConfig("logging:disable", "true")
	app.Route("GET", "/{first}/xxxxx/{second}", func(ctx *Context, req *Request, res *Response) {
		res.Body = "Hello Vox!"
		if req.Params["first"] != "foo" {
			t.Fail()
		}
		if req.Params["second"] != "bar" {
			t.Fail()
		}
	})
	r := httptest.NewRequest("GET", "http://test.com/foo/xxxxx/bar", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Result().StatusCode != 200 {
		t.Fail()
	}
}

func TestRouteShortcut(t *testing.T) {
	methods := []string{
		"Get",
		"Head",
		"Post",
		"Put",
		"Patch",
		"Delete",
		"Options",
		"Trace",
	}

	app := New()
	app.SetConfig("logging:disable", "true")
	for _, method := range methods {
		args := []reflect.Value{}
		args = append(args, reflect.ValueOf("/"))
		args = append(args, reflect.ValueOf(func(ctx *Context, req *Request, res *Response) {
			res.Body = method
		}))
		reflect.ValueOf(app).MethodByName(method).Call(args)
	}

	for _, method := range methods {
		r := httptest.NewRequest(strings.ToUpper(method), "http://test.com/", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, r)
		if w.Result().StatusCode != 200 && w.Body.String() != method {
			t.Fail()
		}
	}
}

func TestRouteFallthrough(t *testing.T) {
	app := New()
	app.SetConfig("logging:disable", "true")
	app.Get("/fallthrough", func(ctx *Context, req *Request, res *Response) {
	})
	app.Use(func(ctx *Context, req *Request, res *Response) {
		res.Body = "fallthrough"
	})
	r := httptest.NewRequest("Get", "http://test.com/", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Result().StatusCode != 200 && w.Body.String() != "fallthrough" {
		t.Fail()
	}
}

func TestRouteWithUnicodeParams(t *testing.T) {
	app := New()
	app.SetConfig("logging:disable", "true")
	app.Route("GET", "/{first}/xxxxx/{second}", func(ctx *Context, req *Request, res *Response) {
		res.Body = "Hello Vox!"
		if req.Params["first"] != "éèçà" {
			t.Fail()
		}
		if req.Params["second"] != "aa_1-.aspx" {
			t.Fail()
		}
	})
	r := httptest.NewRequest("GET", "http://test.com/éèçà/xxxxx/aa_1-.aspx", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Result().StatusCode != 200 {
		t.Fail()
	}
}
