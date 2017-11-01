package vox

import (
	"net/http/httptest"
	"testing"
)

func TestRoute(t *testing.T) {
	app := New()
	app.Route("GET", "/test_route", func(ctx *Context) {
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
