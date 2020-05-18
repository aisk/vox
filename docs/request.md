---
title: Request
nav_order: 5
---

# Request

Vox's `Request` object is built on top of go's native [`net/http.Request`](https://golang.org/pkg/net/http/#Request).

Actually, a `vox.Request` is [embedding](https://golang.org/doc/effective_go.html#embedding) a [`net/http.Request`](https://golang.org/pkg/net/http/#Request) in it's struct definitation. So you can access any [`net/http.Request`](https://golang.org/pkg/net/http/#Request)'s public fields or methods from a `vox.Request`.

For example, you can access a reequest's HTTP header like this:

```go
func ExampleHandler(ctx *vox.Context, req *vox.Request, res *vox.Respomse) {
    fmt.Println("secret from request header: ", req.Header.Get("X-Secret"))
}
```

Additionally, `vox.Request` have some extra fields / methods which [`net/http.Request`](https://golang.org/pkg/net/http/#Request) dose not provides.

For example, vox has a `JSON` method to decode JSON request body to go values, with additional functionality to check the content-type header from the request. If the content-type header does not start with "application/json" or decode errors, this function will return an error and set the response status code to 406.

```go
func PostJSONHandler(ctx *vox.Context, req *vox.Request, res *vox.Respomse) {
    body := map[string]string{}
    if err := req.JSON(&body); err != nil {
        return  // You do not need to set the response's status code for vox have set it.
    }
    ...
}
```