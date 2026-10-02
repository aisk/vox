# Response

Route handlers receive a `*vox.Response[T]` and middleware receives a `*vox.BaseResponse`. Handlers describe the response by assigning fields. Vox writes it to the client after the handler and the middleware around it have returned.

## The response types

```go
type BaseResponse struct {
	Body        any
	Status      int
	Header      http.Header
	Writer      http.ResponseWriter
	DontRespond bool
}

type Response[T any] struct {
	*BaseResponse
	Body T
}
```

In a route handler, `res.Body` has the route's output type `T`. In middleware, `res.Body` accepts any value. Middleware can call `res.HasBody()` to check whether a body has been set.

## Body

The type of the body determines how it is written.

| Body type | Written as | Default `Content-Type` | Default status |
|:--|:--|:--|:--|
| `string`, `[]byte` | the bytes as is | detected from the content | 200 |
| `io.Reader`, `io.ReadCloser` | a copy of the stream | detected from the content | 200 |
| `error` | the error message | detected from the content | 500 |
| `vox.NoBody` | nothing | none | 204 |
| anything else | JSON | `application/json` | 200 |

```go
func text(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
	res.Body = "Hello, World!"
}

func raw(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[[]byte]) {
	res.Body = []byte("Hello, World!")
}

func stream(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[io.Reader]) {
	res.Body = strings.NewReader("Hello, World!")
}

func fail(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[error]) {
	res.Body = errors.New("internal server error")
}

func data(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[map[string]string]) {
	res.Body = map[string]string{"foo": "bar"}
}

func empty(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[vox.NoBody]) {
	// responds with 204 and no body
}
```

An `io.ReadCloser` is closed after it has been copied.

A route commits its body when the handler returns, even when the body still has its zero value. A matched route never falls back to 404 on its own: a handler that leaves a `Response[*User]` body nil responds with `null`, and one that leaves a `Response[int]` untouched responds with `0`.

### Content type

For JSON bodies, Vox sets `Content-Type: application/json` unless the handler has set a content type or the status is 204 or 304.

For the other body types Vox sets no content type. Go's HTTP server then [detects one](https://pkg.go.dev/net/http#DetectContentType) from the first bytes of the body, which gives `text/plain; charset=utf-8` for ordinary text and `text/html; charset=utf-8` for markup. Set the header yourself whenever the type matters:

```go
func page(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
	res.Header.Set("Content-Type", "text/html; charset=utf-8")
	res.Body = "<h1>Hello, World!</h1>"
}
```

## Status

`Status` is the HTTP status code of the response.

```go
func create(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
	res.Status = 201
	res.Body = "created"
}
```

When the status is not set, it defaults as follows:

| Situation | Status |
|:--|:--|
| A route responded with a body | 200 |
| The body is an `error` | 500 |
| The route's output type is `vox.NoBody` | 204 |
| No route matched and no middleware responded | 404 |

Responses to `HEAD` requests and responses with status 204 or 304 are sent without a body.

## Header

`Header` is the `http.Header` of the response. Headers can be changed until the response is written, which includes the code after `ctx.Next()` in middleware.

```go
func download(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[[]byte]) {
	res.Header.Set("Content-Type", "text/csv")
	res.Header.Set("Content-Disposition", `attachment; filename="report.csv"`)
	res.Body = []byte("id,name\n1,Ada\n")
}
```

## Error responses

`T` describes the successful response. Error responses often look different, and there are two ways to send them.

When the status is 400 or above and `Body` still has its zero value, the body is not committed and the status text is written instead:

```go
func showUser(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[User]) {
	user, ok := findUser(req.Params["id"])
	if !ok {
		res.Status = 404 // responds with "Not Found"
		return
	}
	res.Body = user
}
```

To respond with a body that does not fit `T`, assign it to `res.BaseResponse.Body`. It accepts any value and takes precedence over `Body`:

```go
func showUser(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[User]) {
	user, err := loadUser(req.Params["id"])
	if err != nil {
		res.Status = 502
		res.BaseResponse.Body = map[string]string{"error": err.Error()}
		return
	}
	res.Body = user
}
```

[Error Handling](./errors.md) covers these patterns in more detail.

## Redirect

`Redirect` sends the client to another URL with the given status code. A relative URL is resolved against the path of the current request.

```go
func oldPage(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[vox.NoBody]) {
	res.Redirect("/new-location", 302)
}
```

For `GET` requests the response includes a small HTML body with a link to the target. A redirect takes precedence over the `Body` field.

## SetCookie

`SetCookie` adds a `Set-Cookie` header to the response.

```go
func login(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[vox.NoBody]) {
	res.SetCookie(&http.Cookie{
		Name:     "session",
		Value:    "opaque-token",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})
}
```

## Streaming

Assign an `io.Reader` to stream a response without loading it into memory. An `io.ReadCloser`, such as an `*os.File`, is closed for you.

```go
func export(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[io.ReadCloser]) {
	file, err := os.Open("export.csv")
	if err != nil {
		res.Status = 404
		return
	}
	res.Header.Set("Content-Type", "text/csv")
	res.Body = file
}
```

## Writing the response yourself

To use the underlying `http.ResponseWriter` directly, set `DontRespond` to `true`. Vox then writes nothing and ignores `Status` and `Body`. `res.Header` is the writer's own header map, so headers set through it are still sent.

```go
func custom(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[vox.NoBody]) {
	res.DontRespond = true
	res.Writer.WriteHeader(200)
	res.Writer.Write([]byte("Hello, World!"))
}
```

This is how existing `http.Handler` implementations can be called from a handler or middleware:

```go
func legacy(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[vox.NoBody]) {
	res.DontRespond = true
	legacyHandler.ServeHTTP(res.Writer, req.Request)
}
```

## Precedence

When a route handler sets more than one of these, the first in this list wins:

1. `DontRespond`
2. `Redirect`
3. `res.BaseResponse.Body`
4. `res.Body`
