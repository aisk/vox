---
title: Response
nav_order: 6
---

# Response

Route handlers receive `*vox.Response[T]`. Set `Body` to a value of type `T`, `Status` to an HTTP status code, and `Header` to the response headers. The response is written after the handler and the middleware wrapping it return.

## Body

The body type determines how the value is written. Middleware receives `*vox.BaseResponse`, whose `Body` field accepts any value. Use `res.HasBody()` in middleware to check whether a body has been set.

If the value is a `[]byte`, `string`, `io.Reader` or `io.ReadCloser`, it will be written to the response body directly.

```go
func StringHandler(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
    res.Body = "Hello, World!"
}

func BytesHandler(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[[]byte]) {
    res.Body = []byte("Hello, World!")
}

func ReaderHandler(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[io.Reader]) {
    res.Body = strings.NewReader("Hello, World!")
}
```

If the value is an `error`, the error message will be written to the response body and the status code defaults to 500.

```go
func ErrorHandler(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[error]) {
    res.Body = errors.New("internal server error")
}
```

For any other type, it will be marshaled to JSON and the `Content-Type` header will be set to `application/json`, unless the handler has set one or the status is 204 or 304.

```go
func JSONHandler(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[map[string]string]) {
    res.Body = map[string]string{"foo": "bar"}
}
```

## Status

The `Status` field is an `int` type, which will be used as the HTTP response's status code.

A route commits its body on normal return, including its zero value, and defaults to 200. A matched route never falls back to 404 on its own, so a handler that leaves a nil body responds with `null`. An error body defaults to 500. `Response[vox.NoBody]` produces an empty body and defaults to 204. An unmatched request with no middleware response defaults to 404. HEAD, 204, and 304 responses omit the body.

```go
func StatusHandler(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
    res.Status = 201
    res.Body = "created"
}
```

## Errors

`T` describes the successful response. When the status is 400 or above and `Body` still has its zero value, the body is not committed and the status text is written instead.

```go
func UserHandler(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[User]) {
    user, ok := findUser(req.Params["id"])
    if !ok {
        res.Status = 404 // responds with "Not Found"
        return
    }
    res.Body = user
}
```

To respond with a body that does not fit `T`, such as an error message or a stream, assign it to `res.BaseResponse.Body`. It accepts any value and takes precedence over `Body`.

```go
func UserHandler(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[User]) {
    user, err := loadUser(req.Params["id"])
    if err != nil {
        res.Status = 502
        res.BaseResponse.Body = map[string]string{"error": err.Error()}
        return
    }
    res.Body = user
}
```

## Header

The `Header` field is an `http.Header` type, which will be written to the response.

```go
func HeaderHandler(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[vox.NoBody]) {
    res.Header.Set("X-Custom-Header", "foobar")
}
```

## Redirect

The `Redirect` method redirects the request to another URL with a given status code. Its generated body takes precedence over the `Body` field.

```go
func RedirectHandler(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[vox.NoBody]) {
    res.Redirect("/new-location", 302)
}
```

## SetCookie

The `SetCookie` method sets a cookie on the response.

```go
func CookieHandler(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[vox.NoBody]) {
    res.SetCookie(&http.Cookie{Name: "foo", Value: "bar"})
}
```

## DontRespond

If you want to use Go's native `http.ResponseWriter` to write the response, you can set the `DontRespond` field to `true`.

```go
func DontRespondHandler(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[vox.NoBody]) {
    res.DontRespond = true
    res.Writer.WriteHeader(200)
    res.Writer.Write([]byte("Hello, World!"))
}
```
