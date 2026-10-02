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
415. Empty, malformed, type-incompatible JSON or extra data after the first JSON
value yield 400. The route handler is not called on decode failure.

Decoding follows `encoding/json`: unknown fields are allowed, missing fields retain zero values, and `null` can produce nil pointers. Business validation and request size limits must be applied separately.

Use `Request[vox.NoBody]` to skip automatic decoding, regardless of HTTP method. The original stream remains accessible as `req.Request.Body`; the original HTTP request is `req.Request`. This supports uploads and manual decoding.

Middleware receives `*vox.BaseRequest`, whose `Body` is the original stream. Its `JSON(&value)` helper decodes JSON and sets status 406 if the content type or body is invalid. Do not call it on an already decoded request: automatic decoding has consumed the stream.
