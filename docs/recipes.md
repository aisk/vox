---
title: Recipes
nav_order: 5
---

# Recipes
{: .no_toc }

Short, self-contained solutions to common tasks. Each one can be copied into an application and adapted.

1. TOC
{:toc}

---

## Health check

A route that answers every method, useful for load balancer probes:

```go
app.Route("*", "/healthz", func(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
	res.Body = "ok"
})
```

## Request ID

Reuse the ID sent by the client or a proxy, or generate one. The ID is echoed in the response and stored in the context for handlers and loggers.

```go
type requestIDKey struct{}

func requestID(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
	id := req.Header.Get("X-Request-ID")
	if id == "" {
		id = rand.Text()
	}
	res.Header.Set("X-Request-ID", id)
	ctx.Context = context.WithValue(ctx.Context, requestIDKey{}, id)
	ctx.Next()
}

// RequestID returns the ID of the current request.
func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}
```

`rand.Text` comes from `crypto/rand`.

## CORS

Allow a browser application on another origin to call your API. Preflight requests are answered by the middleware and never reach a route.

```go
func cors(origin string) vox.Handler {
	return func(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
		res.Header.Set("Access-Control-Allow-Origin", origin)
		res.Header.Add("Vary", "Origin")

		if req.Method == http.MethodOptions && req.Header.Get("Access-Control-Request-Method") != "" {
			res.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
			res.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			res.Header.Set("Access-Control-Max-Age", "86400")
			res.Status = http.StatusNoContent
			return
		}
		ctx.Next()
	}
}

app.Use(cors("https://app.example.com"))
```

## Bearer token authentication

Reject requests without a valid token, and make the authenticated user available to handlers.

```go
type userKey struct{}

func authenticate(lookup func(token string) (User, bool)) vox.Handler {
	return func(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
		token, ok := strings.CutPrefix(req.Header.Get("Authorization"), "Bearer ")
		user, found := lookup(token)
		if !ok || !found {
			res.Header.Set("WWW-Authenticate", "Bearer")
			res.Status = http.StatusUnauthorized
			return
		}
		ctx.Context = context.WithValue(ctx.Context, userKey{}, user)
		ctx.Next()
	}
}

func me(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[User]) {
	res.Body = ctx.Value(userKey{}).(User)
}
```

## Basic authentication

Protect an application with a single user name and password. The comparison is done in constant time.

```go
func basicAuth(username, password string) vox.Handler {
	return func(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
		u, p, ok := req.BasicAuth()
		userOK := subtle.ConstantTimeCompare([]byte(u), []byte(username)) == 1
		passOK := subtle.ConstantTimeCompare([]byte(p), []byte(password)) == 1
		if !ok || !userOK || !passOK {
			res.Header.Set("WWW-Authenticate", `Basic realm="restricted"`)
			res.Status = http.StatusUnauthorized
			return
		}
		ctx.Next()
	}
}
```

## Security headers

Headers set before `ctx.Next()` apply to every response, including errors and responses that are written directly, such as static files.

```go
func securityHeaders(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
	res.Header.Set("X-Content-Type-Options", "nosniff")
	res.Header.Set("X-Frame-Options", "DENY")
	res.Header.Set("Referrer-Policy", "no-referrer")
	ctx.Next()
}
```

## Request timeout

Give every request a deadline. Handlers that pass `ctx` to their database and HTTP calls are interrupted when it expires.

```go
func timeout(d time.Duration) vox.Handler {
	return func(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
		var cancel context.CancelFunc
		ctx.Context, cancel = context.WithTimeout(ctx.Context, d)
		defer cancel()
		ctx.Next()
	}
}

func report(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[Report]) {
	result, err := buildReport(ctx)
	if errors.Is(err, context.DeadlineExceeded) {
		res.Status = http.StatusGatewayTimeout
		return
	}
	if err != nil {
		res.Status = http.StatusInternalServerError
		return
	}
	res.Body = result
}
```

## File upload

Read a multipart upload in a `NoBody` handler and save it to disk. The size of the request is capped first, because the automatic limit only covers JSON bodies.

```go
func upload(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[map[string]any]) {
	req.Request.Body = http.MaxBytesReader(res.Writer, req.Request.Body, 10<<20)

	file, header, err := req.FormFile("file")
	if err != nil {
		res.Status = http.StatusBadRequest
		res.BaseResponse.Body = "a file field is required and must be at most 10 MiB"
		return
	}
	defer file.Close()

	// Never trust the client's file name as a path.
	dst, err := os.Create(filepath.Join("uploads", filepath.Base(header.Filename)))
	if err != nil {
		res.Status = http.StatusInternalServerError
		return
	}
	defer dst.Close()

	size, err := io.Copy(dst, file)
	if err != nil {
		res.Status = http.StatusInternalServerError
		return
	}
	res.Status = http.StatusCreated
	res.Body = map[string]any{"name": header.Filename, "size": size}
}
```

## File download

Stream a file as an attachment. The file is closed after it has been sent.

```go
func download(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[io.ReadCloser]) {
	file, err := os.Open("reports/2026.csv")
	if err != nil {
		res.Status = http.StatusNotFound
		return
	}
	res.Header.Set("Content-Type", "text/csv")
	res.Header.Set("Content-Disposition", `attachment; filename="2026.csv"`)
	res.Body = file
}
```

To serve a whole directory, use the [Static Files]({% link middlewares/static.md %}) middleware.

## Structured access log

The built-in access log has a fixed format. To log in your own format, turn it off and wrap the application at the `net/http` level, where the final status is visible:

```go
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}

func accessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", recorder.status,
			"duration", time.Since(start),
		)
	})
}

func main() {
	app := vox.New()
	app.SetConfig("logging:disable", "true")
	// register middleware and routes

	http.ListenAndServe("localhost:3000", accessLog(app))
}
```

A Vox middleware is not the right place for this, because the default status is only filled in after all your middleware has returned.

## More

These tasks are covered in the guide:

- [Recovering from panics]({% link guide/errors.md %}#panics)
- [A custom 404 response]({% link guide/errors.md %}#not-found)
- [Formatting all errors as JSON]({% link guide/errors.md %}#formatting-errors-in-one-place)
- [Graceful shutdown]({% link guide/running.md %}#graceful-shutdown)
- [Mounting an application under a prefix]({% link guide/running.md %}#mounting-under-a-prefix)
