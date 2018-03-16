package vox

import (
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestRoute(t *testing.T) {
	app := New()
	app.Route("GET", regexp.MustCompile("/test_route"), func(ctx *Context) {
		ctx.Response.Body = "Hello Vox!"
		ctx.Response.Header.Set("foo", "bar")
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

func TestRouteWithParams(t *testing.T) {
	app := New()
	app.Route("GET", regexp.MustCompile(`/(?P<first>\w+)/\w+/(?P<second>\w+)`), func(ctx *Context) {
		ctx.Response.Body = "Hello Vox!"
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

func TestRouteShortcut(t *testing.T) {
	methods := []string{"Get", "Post", "Put", "Delete", "Option"}

	app := New()
	for _, method := range methods {
		args := []reflect.Value{}
		args = append(args, reflect.ValueOf(regexp.MustCompile("/")))
		args = append(args, reflect.ValueOf(func(ctx *Context) {
			ctx.Response.Body = method
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
	app.Get(regexp.MustCompile("/fallthrough"), func(ctx *Context, next func()) {
		println("called!")
		next()
	})
	app.Use(func(ctx *Context) {
		ctx.Response.Body = "fallthrough"
	})
	r := httptest.NewRequest("Get", "http://test.com/", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Result().StatusCode != 200 && w.Body.String() != "fallthrough" {
		t.Fail()
	}
}