package vox

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRoute(t *testing.T) {
	app := New()
	app.SetConfig("logging:disable", "true")
	app.Route("GET", "/test_route", func(ctx *Context, req *Request[NoBody], res *Response[any]) {
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
	app.Route("*", "/test_route", func(ctx *Context, req *Request[NoBody], res *Response[any]) {
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
	app.Route("GET", "/{first}/xxxxx/{second}", func(ctx *Context, req *Request[NoBody], res *Response[any]) {
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
	app := New()
	app.SetConfig("logging:disable", "true")
	methods := []struct {
		name     string
		register func(string, RouteHandler[NoBody, string])
	}{
		{"GET", app.Get[NoBody, string]},
		{"HEAD", app.Head[NoBody, string]},
		{"POST", app.Post[NoBody, string]},
		{"PUT", app.Put[NoBody, string]},
		{"PATCH", app.Patch[NoBody, string]},
		{"DELETE", app.Delete[NoBody, string]},
		{"OPTIONS", app.Options[NoBody, string]},
		{"TRACE", app.Trace[NoBody, string]},
	}
	for _, method := range methods {
		method.register("/", func(ctx *Context, req *Request[NoBody], res *Response[string]) {
			res.Body = method.name
		})
	}
	for _, method := range methods {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, httptest.NewRequest(method.name, "/", nil))
		if w.Code != 200 || (method.name != "HEAD" && w.Body.String() != method.name) {
			t.Fatalf("%s: status=%d body=%q", method.name, w.Code, w.Body.String())
		}
	}
}

func TestRouteFallthrough(t *testing.T) {
	app := New()
	app.SetConfig("logging:disable", "true")
	app.Get("/fallthrough", func(ctx *Context, req *Request[NoBody], res *Response[any]) {
	})
	app.Use(func(ctx *Context, req *BaseRequest, res *BaseResponse) {
		res.Body = "fallthrough"
	})
	r := httptest.NewRequest("GET", "http://test.com/fallthrough", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Result().StatusCode != 200 || w.Body.String() != "fallthrough" {
		t.Fail()
	}
}

func TestRouteWithUnicodeParams(t *testing.T) {
	app := New()
	app.SetConfig("logging:disable", "true")
	app.Route("GET", "/{first}/xxxxx/{second}", func(ctx *Context, req *Request[NoBody], res *Response[any]) {
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

type createUser struct {
	Name string `json:"name"`
}
type userBody struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestRouteBodyTypes(t *testing.T) {
	app := New()
	app.SetConfig("logging:disable", "true")
	app.Post("/users/{id}", func(ctx *Context, req *Request[createUser], res *Response[userBody]) {
		if req.Params["id"] != "42" || req.Header.Get("X-Test") != "yes" || ctx.App != app {
			t.Fatal("request metadata or context lost")
		}
		res.Status = http.StatusCreated
		res.Header.Set("X-Route", "users")
		res.SetCookie(&http.Cookie{Name: "session", Value: "value"})
		res.Body = userBody{42, req.Body.Name}
	})
	app.Get[NoBody, string]("/health", func(_ *Context, _ *Request[NoBody], res *Response[string]) {
		res.Body = "ok"
	})
	app.Use(func(ctx *Context, req *BaseRequest, res *BaseResponse) {
		if req.URL.Path == "/users/42" {
			body, ok := res.Body.(userBody)
			if !ok || body.Name != "Ada" {
				t.Fatalf("body not committed before middleware: %#v", res.Body)
			}
		}
		res.Header.Set("X-Middleware", "yes")
		ctx.Next()
	})
	r := httptest.NewRequest("POST", "/users/42", strings.NewReader(`{"name":"Ada"}`))
	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	r.Header.Set("X-Test", "yes")
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Code != 201 || w.Body.String() != `{"id":42,"name":"Ada"}` || w.Header().Get("Content-Type") != "application/json" || w.Header().Get("X-Route") != "users" || w.Header().Get("X-Middleware") != "yes" || len(w.Result().Cookies()) != 1 {
		t.Fatalf("unexpected response: %d %s %v", w.Code, w.Body.String(), w.Header())
	}
	w = httptest.NewRecorder()
	app.ServeHTTP(w, httptest.NewRequest("GET", "/health", nil))
	if w.Code != 200 || w.Body.String() != "ok" {
		t.Fatalf("health: %d %q", w.Code, w.Body.String())
	}
}

func TestRouteDecode(t *testing.T) {
	cases := []struct {
		name, contentType, body string
		status                  int
	}{
		{"valid", "application/json", `{"name":"Ada"}`, 200},
		{"suffix", "application/problem+json", `{"name":"Ada"}`, 200},
		{"whitespace", "application/json", "{\"name\":\"Ada\"}\n ", 200},
		{"missing type", "", `{}`, 415},
		{"wrong type", "text/plain", `{}`, 415},
		{"invalid type", "application/json-invalid", `{}`, 415},
		{"bad parameter", "application/json; charset", `{}`, 415},
		{"empty", "application/json", "", 400},
		{"syntax", "application/json", `{`, 400},
		{"field type", "application/json", `{"name":42}`, 400},
		{"second value", "application/json", `{} {}`, 400},
		{"trailing garbage", "application/json", `{} nope`, 400},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := New()
			app.SetConfig("logging:disable", "true")
			called := false
			app.Post("/", func(_ *Context, req *Request[createUser], res *Response[string]) {
				called = true
				res.Body = req.Body.Name
			})
			r := httptest.NewRequest("POST", "/", strings.NewReader(tc.body))
			r.Header.Set("Content-Type", tc.contentType)
			w := httptest.NewRecorder()
			app.ServeHTTP(w, r)
			if w.Code != tc.status || called != (tc.status == 200) {
				t.Fatalf("status=%d called=%v", w.Code, called)
			}
		})
	}
}

func TestRouteResponseSemantics(t *testing.T) {
	app := New()
	app.SetConfig("logging:disable", "true")
	app.Get("/zero", func(_ *Context, _ *Request[NoBody], _ *Response[int]) {})
	app.Get("/struct", func(_ *Context, _ *Request[NoBody], _ *Response[struct{}]) {})
	app.Get("/struct-pointer", func(_ *Context, _ *Request[NoBody], res *Response[*struct{}]) { res.Body = &struct{}{} })
	app.Get("/nil", func(_ *Context, _ *Request[NoBody], _ *Response[*userBody]) {})
	app.Get("/empty", func(_ *Context, _ *Request[NoBody], _ *Response[NoBody]) {})
	app.Get("/explicit-empty", func(_ *Context, _ *Request[NoBody], res *Response[NoBody]) { res.Status = 202 })
	app.Get("/redirect", func(_ *Context, _ *Request[NoBody], res *Response[userBody]) { res.Redirect("/target", 302) })
	app.Post("/redirect", func(_ *Context, _ *Request[NoBody], res *Response[userBody]) { res.Redirect("/target", 303) })
	app.Get("/raw", func(_ *Context, _ *Request[NoBody], res *Response[userBody]) {
		res.DontRespond = true
		res.Writer.WriteHeader(202)
		io.WriteString(res.Writer, "raw")
	})
	app.Head("/head", func(_ *Context, _ *Request[NoBody], res *Response[string]) { res.Body = "hidden" })
	app.Get("/204", func(_ *Context, _ *Request[NoBody], res *Response[userBody]) { res.Status = 204 })
	for _, tc := range []struct {
		method, path string
		status       int
		body         string
	}{
		{"GET", "/struct-pointer", 200, "{}"}, {"GET", "/zero", 200, "0"}, {"GET", "/struct", 200, "{}"}, {"GET", "/nil", 200, "null"},
		{"GET", "/empty", 204, ""}, {"GET", "/explicit-empty", 202, ""},
		{"GET", "/redirect", 302, "<a href=\"/target\">Found</a>.\n"}, {"POST", "/redirect", 303, ""},
		{"GET", "/raw", 202, "raw"}, {"HEAD", "/head", 200, ""}, {"GET", "/204", 204, ""},
		{"GET", "/missing", 404, "Not Found"},
	} {
		t.Run(tc.method+tc.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			app.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, nil))
			if w.Code != tc.status || w.Body.String() != tc.body {
				t.Fatalf("status=%d body=%q", w.Code, w.Body.String())
			}
			if tc.path == "/redirect" && w.Header().Get("Location") != "/target" {
				t.Fatal("redirect header missing")
			}
		})
	}
}

func TestRouteRawRequestAndPointerInput(t *testing.T) {
	app := New()
	app.SetConfig("logging:disable", "true")
	app.Post("/raw", func(_ *Context, req *Request[NoBody], res *Response[string]) {
		b, err := io.ReadAll(req.Request.Body)
		if err != nil {
			t.Fatal(err)
		}
		res.Body = string(b)
	})
	app.Post("/pointer", func(_ *Context, req *Request[*createUser], res *Response[string]) { res.Body = req.Body.Name })
	for _, tc := range []struct{ path, body, want string }{{"/raw", "not json", "not json"}, {"/pointer", `{"name":"Ada"}`, "Ada"}} {
		r := httptest.NewRequest("POST", tc.path, strings.NewReader(tc.body))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.ServeHTTP(w, r)
		if w.Code != 200 || w.Body.String() != tc.want {
			t.Fatalf("%s: %d %q", tc.path, w.Code, w.Body.String())
		}
	}
}

func TestRouteCannotAdvanceMiddleware(t *testing.T) {
	app := New()
	app.SetConfig("logging:disable", "true")
	app.Get("/", func(ctx *Context, _ *Request[NoBody], _ *Response[string]) { ctx.Next() })
	w := httptest.NewRecorder()
	defer func() {
		if recover() == nil {
			t.Error("expected route Next to panic")
		}
		if w.Flushed || w.Body.Len() != 0 {
			t.Error("response written before handler completed")
		}
	}()
	app.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
}

func TestRouteSharesContext(t *testing.T) {
	type key struct{}
	app := New()
	app.SetConfig("logging:disable", "true")
	var routeContext *Context
	app.Get("/", func(ctx *Context, _ *Request[NoBody], res *Response[string]) {
		routeContext = ctx
		ctx.Context = context.WithValue(ctx.Context, key{}, "value")
		res.Body = "ok"
	})
	calls := 0
	app.Use(func(ctx *Context, _ *BaseRequest, _ *BaseResponse) {
		calls++
		if ctx != routeContext || ctx.Value(key{}) != "value" {
			t.Fatal("route context was not shared")
		}
		ctx.Next()
	})
	w := httptest.NewRecorder()
	app.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	if calls != 1 || w.Body.String() != "ok" {
		t.Fatalf("calls=%d body=%q", calls, w.Body.String())
	}
}
