# Testing

An application is an `http.Handler`, so the standard [`net/http/httptest`](https://pkg.go.dev/net/http/httptest) package is all you need to test it.

## Structuring for tests

Build the application in a function, separate from `main`, so that tests can create one without starting a server:

```go
func newApp() *vox.Application {
	app := vox.New()
	app.Get("/hello/{name}", hello)
	app.Post("/users", createUser)
	return app
}

func main() {
	newApp().Run("localhost:3000")
}
```

## Testing a route

Create a request, record the response, and check it. No network is involved.

```go
func TestHello(t *testing.T) {
	app := newApp()
	app.SetConfig("logging:disable", "true")

	w := httptest.NewRecorder()
	app.ServeHTTP(w, httptest.NewRequest("GET", "/hello/gopher", nil))

	if w.Code != 200 {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if got := w.Body.String(); got != "Hello, gopher!" {
		t.Fatalf("body = %q", got)
	}
}
```

Setting `logging:disable` keeps the access log out of the test output.

## Testing JSON routes

Remember the `Content-Type` header. Without it the request is rejected with 415 before it reaches the handler.

```go
func TestCreateUser(t *testing.T) {
	app := newApp()
	app.SetConfig("logging:disable", "true")

	r := httptest.NewRequest("POST", "/users", strings.NewReader(`{"name":"Ada"}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)

	if w.Code != 201 {
		t.Fatalf("status = %d, want 201", w.Code)
	}
	var user User
	if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
		t.Fatal(err)
	}
	if user.Name != "Ada" {
		t.Fatalf("name = %q", user.Name)
	}
}
```

## Testing middleware

Test a middleware by installing it in a small application with a route that makes its effect visible:

```go
func TestRequireToken(t *testing.T) {
	app := vox.New()
	app.SetConfig("logging:disable", "true")
	app.Use(requireToken)
	app.Get("/", func(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
		res.Body = "secret"
	})

	w := httptest.NewRecorder()
	app.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	if w.Code != 403 {
		t.Fatalf("without token: status = %d, want 403", w.Code)
	}

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-API-Token", "a-secret")
	w = httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Code != 200 || w.Body.String() != "secret" {
		t.Fatalf("with token: status = %d, body = %q", w.Code, w.Body.String())
	}
}
```

## Testing over a real connection

`httptest.NewRecorder` records what the handler writes. A few things are only added by Go's HTTP server, for example the detected `Content-Type` of a plain text response. To test those, or to exercise a real client, start a test server:

```go
func TestOverHTTP(t *testing.T) {
	app := newApp()
	app.SetConfig("logging:disable", "true")

	server := httptest.NewServer(app)
	defer server.Close()

	resp, err := http.Get(server.URL + "/hello/gopher")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if ct := resp.Header.Get("Content-Type"); ct != "text/plain; charset=utf-8" {
		t.Fatalf("content type = %q", ct)
	}
}
```
