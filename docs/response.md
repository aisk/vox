---
title: Response
nav_order: 6
---

# Response

Vox's `Response` object is built on top of go's native [`net/http.ResponseWriter`](https://golang.org/pkg/net/http/#ResponseWriter).

A `vox.Response` contains all the information which will be written to the HTTP client.

## Body

The `Body` field is an `interface{}` type, so you can assign any type of value to it.

If the value is a `[]byte`, `string`, `io.Reader` or `io.ReadCloser`, it will be written to the response body directly.

```go
func StringHandler(ctx *vox.Context, req *vox.Request, res *vox.Response) {
    res.Body = "Hello, World!"
}

func BytesHandler(ctx *vox.Context, req *vox.Request, res *vox.Response) {
    res.Body = []byte("Hello, World!")
}

func ReaderHandler(ctx *vox.Context, req *vox.Request, res *vox.Response) {
    res.Body = strings.NewReader("Hello, World!")
}
```

If the value is an `error`, the error message will be written to the response body and the status code will be set to 500.

```go
func ErrorHandler(ctx *vox.Context, req *vox.Request, res *vox.Response) {
    res.Body = errors.New("internal server error")
}
```

For any other type, it will be marshaled to JSON and the `Content-Type` header will be set to `application/json`.

```go
func JSONHandler(ctx *vox.Context, req *vox.Request, res *vox.Response) {
    res.Body = map[string]string{"foo": "bar"}
}
```

## Status

The `Status` field is an `int` type, which will be used as the HTTP response's status code.

If not set, `200` will be used as the default value. If the `Body` is not set, `404` will be used. If the `Body` is an `error`, `500` will be used.

```go
func StatusHandler(ctx *vox.Context, req *vox.Request, res *vox.Response) {
    res.Status = 201
    res.Body = "created"
}
```

## Header

The `Header` field is an `http.Header` type, which will be written to the response.

```go
func HeaderHandler(ctx *vox.Context, req *vox.Request, res *vox.Response) {
    res.Header.Set("X-Custom-Header", "foobar")
}
```

## Redirect

The `Redirect` method redirects the request to another URL with a given status code.

```go
func RedirectHandler(ctx *vox.Context, req *vox.Request, res *vox.Response) {
    res.Redirect("/new-location", 302)
}
```

## SetCookie

The `SetCookie` method sets a cookie on the response.

```go
func CookieHandler(ctx *vox.Context, req *vox.Request, res *vox.Response) {
    res.SetCookie(&http.Cookie{Name: "foo", Value: "bar"})
}
```

## DontRespond

If you want to use the go's native `http.ResponseWriter` to write the response, you can set the `DontRespond` field to `true`.

```go
func DontRespondHandler(ctx *vox.Context, req *vox.Request, res *vox.Response) {
    res.DontRespond = true
    res.Writer.WriteHeader(200)
    res.Writer.Write([]byte("Hello, World!"))
}
```
