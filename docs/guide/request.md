# Request

Route handlers receive a `*vox.Request[T]` and middleware receives a `*vox.BaseRequest`. Both give access to the underlying `*http.Request`.

## The request types

```go
type BaseRequest struct {
	*http.Request
	Params map[string]string
}

type Request[T any] struct {
	*BaseRequest
	Body T
}
```

`BaseRequest` embeds the standard [`*http.Request`](https://pkg.go.dev/net/http#Request) and adds the route parameters. `Request[T]` embeds `BaseRequest` and adds the decoded body. Thanks to the embedding, every field and method of `http.Request` is available directly on `req`.

The one name that differs between the two is `Body`:

| | In a route handler (`*Request[T]`) | In middleware (`*BaseRequest`) |
|:--|:--|:--|
| `req.Body` | the decoded value of type `T` | the raw body stream |
| `req.Request.Body` | the raw body stream | the raw body stream |

## Reading request data

```go
func inspect(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[map[string]string]) {
	res.Body = map[string]string{
		"method": req.Method,
		"path":   req.URL.Path,
		"id":     req.Params["id"],            // path parameter
		"page":   req.URL.Query().Get("page"), // query string
		"agent":  req.Header.Get("User-Agent"),
		"remote": req.RemoteAddr,
	}
}
```

Path parameters are described in [Routing](./routing.md#path-parameters). Cookies are read with the standard method:

```go
cookie, err := req.Cookie("session")
if err != nil {
	// the cookie is not present
}
```

## JSON bodies

When the input type of a route is anything other than `vox.NoBody`, Vox decodes the request body as JSON into `req.Body` before the handler runs.

```go
type CreateUser struct {
	Name string `json:"name"`
}

func createUser(ctx *vox.Context, req *vox.Request[CreateUser], res *vox.Response[string]) {
	res.Body = "Hello, " + req.Body.Name
}
```

A body that cannot be decoded is rejected, and the route handler is not called:

| Condition | Status | Response body |
|:--|:--|:--|
| `Content-Type` missing, invalid, or not JSON | 415 | `content type must be application/json` |
| Body larger than the size limit | 413 | `request body is too large` |
| Empty body | 400 | `request body is empty` |
| Body is `null` | 400 | `request body must not be null` |
| Malformed JSON, or extra data after the first value | 400 | `malformed JSON at offset 12` |
| A field has the wrong type | 400 | `invalid type for field "name"` |
| The top-level value has the wrong type | 400 | `invalid type for request body` |

`application/json` and any `application/*+json` media type are accepted, including parameters such as `charset=utf-8`.

Decoding follows [`encoding/json`](https://pkg.go.dev/encoding/json): unknown fields are ignored and missing fields keep their zero values. Vox checks that the body is well formed, not that it makes sense. Validate required fields and value ranges in the handler:

```go
func createUser(ctx *vox.Context, req *vox.Request[CreateUser], res *vox.Response[User]) {
	if req.Body.Name == "" {
		res.Status = 422
		res.BaseResponse.Body = "name is required"
		return
	}
	res.Status = 201
	res.Body = User{ID: 1, Name: req.Body.Name}
}
```

### Body size limit

JSON request bodies are limited to 1 MiB by default. Set `request:max-body-size` to a number of bytes to change the limit, or to `0` to disable it. An invalid value panics when it is set.

```go
app.SetConfig("request:max-body-size", "4194304") // 4 MiB
```

The limit applies to JSON decoding only. A `Request[vox.NoBody]` handler reads the original stream without a limit and can wrap it with `http.MaxBytesReader` itself.

### Decode errors

A rejected body is answered with one of the short plain text messages listed above. Middleware sees the failure as a `*vox.DecodeError` in `res.Body` after `ctx.Next()` returns:

```go
type DecodeError struct {
	Status  int    // the HTTP status the failure maps to
	Message string // safe to send to the client
	Err     error  // the underlying error, if any
}
```

Use it to log the failure or to replace the response with your own format:

```go
app.Use(func(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
	ctx.Next()
	if err, ok := res.Body.(*vox.DecodeError); ok {
		log.Printf("decode %s: %v", req.URL.Path, err.Err)
		res.Body = map[string]string{"error": err.Message}
	}
})
```

## Reading the body yourself

Use `Request[vox.NoBody]` to skip automatic decoding, whatever the HTTP method. The original stream is `req.Request.Body`.

```go
func upload(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
	body := http.MaxBytesReader(res.Writer, req.Request.Body, 10<<20)
	data, err := io.ReadAll(body)
	if err != nil {
		res.Status = 413
		return
	}
	res.Body = fmt.Sprintf("received %d bytes", len(data))
}
```

### Forms and uploads

The form helpers of `http.Request` work in `NoBody` handlers:

```go
func login(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
	username := req.FormValue("username")
	res.Body = "Welcome, " + username
}

func avatar(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
	file, header, err := req.FormFile("avatar")
	if err != nil {
		res.Status = 400
		res.Body = "avatar is required"
		return
	}
	defer file.Close()
	res.Body = "received " + header.Filename
}
```

### Decoding JSON in middleware

`BaseRequest` has a `JSON` method, which is what automatic decoding calls. It applies the same content type, size and syntax rules. On failure it sets the response status and body and returns the `*vox.DecodeError`, so the caller only needs to return:

```go
func audit(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
	var event AuditEvent
	if err := req.JSON(&event); err != nil {
		return
	}
	// ...
}
```

> [!WARNING]
> A request body can be read only once. Do not call `JSON` for a request whose route also decodes the body, and do not call it on a request that was already decoded.
