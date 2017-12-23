package vox

import (
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestRoute(t *testing.T) {
	app := New()
	app.Route("GET", regexp.MustCompile("/test_route"), func(ctx *Context) {
		ctx.Response.SetStatus(200)
		ctx.Response.SetBody("Hello Vox!")
		ctx.Response.Header.Set("foo", "bar")
	})
	r := httptest.NewRequest("GET", "http://test.com/test_route", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Result().StatusCode != 200 {
		t.Fail()
	}

	r = httptest.NewRequest("GET", "http://test.com/invalid_path", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Result().StatusCode != 404 {
		t.Fail()
	}
}

func TestRouteWithParams(t *testing.T) {
	app := New()
	app.Route("GET", regexp.MustCompile(`/(?P<first>\w+)/\w+/(?P<second>\w+)`), func(ctx *Context) {
		ctx.Response.SetStatus(200)
		ctx.Response.SetBody("Hello Vox!")
		if ctx.Request.Params["first"] != "foo" {
			t.Fail()
		}
		if ctx.Request.Params["second"] != "bar" {
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
