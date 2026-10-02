---
title: Request
nav_order: 5
---

# Request

Route handlers receive `*vox.Request[T]`. Its `Body` field has type `T`, decoded from JSON before the handler runs. Headers, URL, route `Params`, and other HTTP metadata are available through the embedded `*vox.BaseRequest`.

```go
type CreateUser struct {
    Name string `json:"name"`
}

func createUser(ctx *vox.Context, req *vox.Request[CreateUser], res *vox.Response[string]) {
    res.Body = "Hello, " + req.Body.Name
}
```

`application/json` and `application/*+json` media types are accepted, including parameters such as `charset=utf-8`. Unsupported or invalid content types yield
415. A body over the size limit yields 413. Empty, `null`, malformed, type-incompatible JSON or extra data after the first JSON
value yield 400. The route handler is not called on decode failure.

Decoding follows `encoding/json`: unknown fields are allowed and missing fields retain zero values. Business validation must be applied separately.

## Body size limit

JSON request bodies are limited to 1 MiB by default. Set `request:max-body-size` to a number of bytes to change the limit, or to `0` to disable it. An invalid value panics when it is set.

```go
app.SetConfig("request:max-body-size", "4194304") // 4 MiB
```

The limit applies to JSON decoding only. A `Request[vox.NoBody]` handler reads the original stream without a limit and can wrap it with `http.MaxBytesReader` itself.

## Decode errors

A rejected body is responded to with a short plain text message, such as `malformed JSON at offset 12` or `invalid type for field "name"`. The response body seen by middleware is a `*vox.DecodeError`, whose `Err` field holds the underlying error. Middleware can use it to log the failure or to replace the response.

```go
app.Use(func(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
    if err, ok := res.Body.(*vox.DecodeError); ok {
        log.Printf("decode %s: %v", req.URL.Path, err.Err)
        res.Body = map[string]string{"error": err.Message}
    }
    ctx.Next()
})
```

## Manual decoding

Use `Request[vox.NoBody]` to skip automatic decoding, regardless of HTTP method. The original stream remains accessible as `req.Request.Body`; the original HTTP request is `req.Request`. This supports uploads and manual decoding.

Middleware receives `*vox.BaseRequest`, whose `Body` is the original stream. Its `JSON(&value)` helper is what automatic decoding calls, so the same rules apply. On failure it sets the response status and body and returns the `*vox.DecodeError`. Do not call it on an already decoded request: automatic decoding has consumed the stream.
