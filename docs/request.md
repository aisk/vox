---
title: Request
nav_order: 5
---

# Request

Vox's `Request` object is built on top of Go's native [`net/http.Request`](https://golang.org/pkg/net/http/#Request).

Actually, `vox.Request` [embeds](https://golang.org/doc/effective_go.html#embedding) a [`net/http.Request`](https://golang.org/pkg/net/http/#Request). So you can access any public field or method from [`net/http.Request`](https://golang.org/pkg/net/http/#Request) directly on a `vox.Request`.

For example, you can access a request's HTTP header like this:

```go
func ExampleHandler(ctx *vox.Context, req *vox.Request, res *vox.Response) {
    fmt.Println("secret from request header: ", req.Header.Get("X-Secret"))
}
```

Additionally, `vox.Request` has some extra fields and methods that [`net/http.Request`](https://golang.org/pkg/net/http/#Request) does not provide.

For example, Vox has a `JSON` method to decode a JSON request body into Go values, with additional logic to validate the `Content-Type` header. If `Content-Type` does not start with "application/json", or a decode error occurs, this function returns an error and sets the response status code to 406.

```go
func PostJSONHandler(ctx *vox.Context, req *vox.Request, res *vox.Response) {
    body := map[string]string{}
    if err := req.JSON(&body); err != nil {
        return  // You do not need to set the response's status code, as vox has already set it.
    }
    // ...
}
```
