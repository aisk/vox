package vox

import (
	"context"
	"errors"
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
	app.Get("/matched", func(ctx *Context, req *Request[NoBody], res *Response[string]) {
		res.Body = "matched"
	})
	app.Use(func(ctx *Context, req *BaseRequest, res *BaseResponse) {
		ctx.Next()
		if !res.HasBody() {
			res.Body = "fallthrough"
		}
	})
	for path, want := range map[string]string{"/matched": "matched", "/unmatched": "fallthrough"} {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, httptest.NewRequest("GET", "http://test.com"+path, nil))
		if w.Result().StatusCode != 200 || w.Body.String() != want {
			t.Errorf("%s: status=%d body=%q", path, w.Code, w.Body.String())
		}
	}
}

func TestMiddlewareWrapsRoute(t *testing.T) {
	app := New()
	app.SetConfig("logging:disable", "true")
	var order []string
	middleware := func(name string) Handler {
		return func(ctx *Context, req *BaseRequest, res *BaseResponse) {
			order = append(order, name+" before")
			ctx.Next()
			order = append(order, name+" after")
		}
	}
	// Routes run last regardless of where they are registered.
	app.Use(middleware("a"))
	app.Get("/", func(ctx *Context, req *Request[NoBody], res *Response[string]) {
		order = append(order, "route")
		res.Body = "ok"
	})
	app.Use(middleware("b"))
	w := httptest.NewRecorder()
	app.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	want := "a before, b before, route, b after, a after"
	if got := strings.Join(order, ", "); got != want || w.Body.String() != "ok" {
		t.Fatalf("order=%q body=%q", got, w.Body.String())
	}
}

func TestMiddlewareGuardsRoute(t *testing.T) {
	app := New()
	app.SetConfig("logging:disable", "true")
	called := false
	app.Get("/", func(ctx *Context, req *Request[NoBody], res *Response[string]) {
		called = true
		res.Body = "secret"
	})
	app.Use(func(ctx *Context, req *BaseRequest, res *BaseResponse) {
		if req.Header.Get("X-API-Token") != "a-secret" {
			res.Status = http.StatusForbidden
			res.Body = "denied"
			return
		}
		ctx.Next()
	})
	w := httptest.NewRecorder()
	app.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	if called || w.Code != http.StatusForbidden || w.Body.String() != "denied" {
		t.Fatalf("called=%v status=%d body=%q", called, w.Code, w.Body.String())
	}
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-API-Token", "a-secret")
	w = httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if !called || w.Code != 200 || w.Body.String() != "secret" {
		t.Fatalf("called=%v status=%d body=%q", called, w.Code, w.Body.String())
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
		ctx.Next()
		if req.URL.Path == "/users/42" {
			body, ok := res.Body.(userBody)
			if !ok || body.Name != "Ada" {
				t.Fatalf("body not committed when the route returned: %#v", res.Body)
			}
		}
		res.Header.Set("X-Middleware", "yes")
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
	const unsupported = "content type must be application/json"
	cases := []struct {
		name, contentType, body string
		status                  int
		want                    string
	}{
		{"valid", "application/json", `{"name":"Ada"}`, 200, "Ada"},
		{"suffix", "application/problem+json", `{"name":"Ada"}`, 200, "Ada"},
		{"whitespace", "application/json", "{\"name\":\"Ada\"}\n ", 200, "Ada"},
		{"missing type", "", `{}`, 415, unsupported},
		{"wrong type", "text/plain", `{}`, 415, unsupported},
		{"invalid type", "application/json-invalid", `{}`, 415, unsupported},
		{"bad parameter", "application/json; charset", `{}`, 415, unsupported},
		{"empty", "application/json", "", 400, "request body is empty"},
		{"blank", "application/json", " \n", 400, "request body is empty"},
		{"null", "application/json", " null\n", 400, "request body must not be null"},
		{"syntax", "application/json", `{`, 400, "malformed JSON at offset 1"},
		{"field type", "application/json", `{"name":42}`, 400, `invalid type for field "name"`},
		{"body type", "application/json", `42`, 400, "invalid type for request body"},
		{"second value", "application/json", `{} {}`, 400, "malformed JSON at offset 4"},
		{"trailing garbage", "application/json", `{} nope`, 400, "malformed JSON at offset 4"},
		{"too large", "application/json", `{"name":"` + strings.Repeat("a", 64) + `"}`, 413, "request body is too large"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := New()
			app.SetConfig("logging:disable", "true")
			app.SetConfig("request:max-body-size", "64")
			called := false
			app.Post("/", func(_ *Context, req *Request[createUser], res *Response[string]) {
				called = true
				res.Body = req.Body.Name
			})
			var decodeErr *DecodeError
			app.Use(func(ctx *Context, _ *BaseRequest, res *BaseResponse) {
				ctx.Next()
				decodeErr, _ = res.Body.(*DecodeError)
			})
			r := httptest.NewRequest("POST", "/", strings.NewReader(tc.body))
			r.Header.Set("Content-Type", tc.contentType)
			w := httptest.NewRecorder()
			app.ServeHTTP(w, r)
			if w.Code != tc.status || w.Body.String() != tc.want || called != (tc.status == 200) {
				t.Fatalf("status=%d body=%q called=%v", w.Code, w.Body.String(), called)
			}
			if (decodeErr != nil) != (tc.status != 200) {
				t.Fatalf("decode error not exposed to middleware: %v", decodeErr)
			}
			if w.Header().Get("Content-Type") == "application/json" {
				t.Fatal("text response labeled as JSON")
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
	app.Get("/error-status", func(_ *Context, _ *Request[NoBody], res *Response[userBody]) { res.Status = 404 })
	app.Get("/error-message", func(_ *Context, _ *Request[NoBody], res *Response[string]) {
		res.Status = 400
		res.Body = "invalid id"
	})
	app.Get("/error-typed", func(_ *Context, _ *Request[NoBody], res *Response[userBody]) {
		res.Status = 409
		res.Body = userBody{1, "Ada"}
	})
	app.Get("/error-empty", func(_ *Context, _ *Request[NoBody], res *Response[NoBody]) { res.Status = 404 })
	app.Get("/base-error", func(_ *Context, _ *Request[NoBody], res *Response[userBody]) {
		res.Status = 500
		res.BaseResponse.Body = errors.New("boom")
	})
	app.Get("/base-body", func(_ *Context, _ *Request[NoBody], res *Response[userBody]) {
		res.Body = userBody{1, "Ada"}
		res.BaseResponse.Body = "base"
	})
	app.Get("/base-stream", func(_ *Context, _ *Request[NoBody], res *Response[NoBody]) {
		res.BaseResponse.Body = strings.NewReader("stream")
	})
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
		{"GET", "/error-status", 404, "Not Found"}, {"GET", "/error-message", 400, "invalid id"},
		{"GET", "/error-typed", 409, `{"id":1,"name":"Ada"}`}, {"GET", "/error-empty", 404, ""},
		{"GET", "/base-error", 500, "boom"}, {"GET", "/base-body", 200, "base"}, {"GET", "/base-stream", 200, "stream"},
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
			if (tc.path == "/204" || tc.path == "/base-error" || tc.path == "/base-stream") && w.Header().Get("Content-Type") == "application/json" {
				t.Fatal("non-JSON response labeled as JSON")
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
	for _, tc := range []struct {
		path, body string
		status     int
		want       string
	}{
		{"/raw", "not json", 200, "not json"},
		{"/raw", strings.Repeat("a", 2<<20), 200, strings.Repeat("a", 2<<20)},
		{"/pointer", `{"name":"Ada"}`, 200, "Ada"},
		{"/pointer", `null`, 400, "request body must not be null"},
	} {
		r := httptest.NewRequest("POST", tc.path, strings.NewReader(tc.body))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.ServeHTTP(w, r)
		if w.Code != tc.status || w.Body.String() != tc.want {
			t.Fatalf("%s: %d %.40q", tc.path, w.Code, w.Body.String())
		}
	}
}

func TestRouteMaxBodySize(t *testing.T) {
	body := `{"name":"` + strings.Repeat("a", 1<<20) + `"}`
	for _, tc := range []struct {
		name, limit string
		status      int
	}{{"default", "", 413}, {"raised", "2097152", 200}, {"disabled", "0", 200}} {
		t.Run(tc.name, func(t *testing.T) {
			app := New()
			app.SetConfig("logging:disable", "true")
			if tc.limit != "" {
				app.SetConfig("request:max-body-size", tc.limit)
			}
			app.Post("/", func(_ *Context, req *Request[createUser], res *Response[int]) { res.Body = len(req.Body.Name) })
			r := httptest.NewRequest("POST", "/", strings.NewReader(body))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			app.ServeHTTP(w, r)
			if w.Code != tc.status {
				t.Fatalf("status=%d", w.Code)
			}
		})
	}
	for _, value := range []string{"-1", "1MB", ""} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("expected invalid max body size %q to panic", value)
				}
			}()
			New().SetConfig("request:max-body-size", value)
		}()
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
		ctx.Next()
		if ctx != routeContext || ctx.Value(key{}) != "value" {
			t.Fatal("route context was not shared")
		}
	})
	w := httptest.NewRecorder()
	app.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	if calls != 1 || w.Body.String() != "ok" {
		t.Fatalf("calls=%d body=%q", calls, w.Body.String())
	}
}
