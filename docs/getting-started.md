# Getting Started

This page walks through a small application: a plain text route, a JSON route and a middleware.

## Requirements

Vox requires **Go 1.27 or later**. Routes are registered through generic methods, which older Go versions cannot compile.

## Installation

Create a module and add Vox to it:

```sh
mkdir hello && cd hello
go mod init example.com/hello
go get github.com/aisk/vox
```

## A first route

Put this in `main.go`:

```go
package main

import "github.com/aisk/vox"

func main() {
	app := vox.New()

	app.Get("/hello/{name}", func(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
		res.Body = "Hello, " + req.Params["name"] + "!"
	})

	app.Run("localhost:3000")
}
```

Run it and send a request:

```sh
$ go run .
$ curl localhost:3000/hello/gopher
Hello, gopher!
```

Three things happened here:

- `app.Get` registered a handler for `GET /hello/{name}`. The `{name}` segment is a path parameter, available as `req.Params["name"]`.
- The handler's signature says what it consumes and produces. `Request[vox.NoBody]` means the request body is not decoded, and `Response[string]` means the response body is a string.
- The handler only assigned `res.Body`. Vox wrote the status, headers and body after the handler returned.

Vox also printed an access log line for the request:

```
127.0.0.1:53412 - - [03/Oct/2026:10:00:00 +0000] "GET /hello/gopher HTTP/1.1" 200 14
```

## A JSON route

Declare the request and response bodies as Go types, and Vox handles JSON for you:

```go
type CreateUser struct {
	Name string `json:"name"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func createUser(ctx *vox.Context, req *vox.Request[CreateUser], res *vox.Response[User]) {
	res.Status = 201
	res.Body = User{ID: 1, Name: req.Body.Name}
}
```

Register it in `main` with `app.Post("/users", createUser)`, then try it:

```sh
$ curl -i localhost:3000/users -H 'Content-Type: application/json' -d '{"name":"Ada"}'
HTTP/1.1 201 Created
Content-Type: application/json

{"id":1,"name":"Ada"}
```

By the time `createUser` runs, `req.Body` is a decoded `CreateUser`. Requests that cannot be decoded never reach the handler:

```sh
$ curl -i localhost:3000/users -H 'Content-Type: application/json' -d '{"name":42}'
HTTP/1.1 400 Bad Request

invalid type for field "name"

$ curl -i localhost:3000/users -d 'name=Ada'
HTTP/1.1 415 Unsupported Media Type

content type must be application/json
```

See [Request](./guide/request.md) for the decoding rules and [Response](./guide/response.md) for how each body type is written.

## A middleware

Middleware wraps every route. This one measures how long the rest of the chain takes and reports it in a header:

```go
app.Use(func(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
	start := time.Now()
	ctx.Next()
	res.Header.Set("X-Response-Time", time.Since(start).String())
})
```

`ctx.Next()` runs the remaining middleware and the matched route. Code before it runs on the way in, and code after it runs on the way out. Middleware works with `BaseRequest` and `BaseResponse`, the untyped views of the request and response.

## Next steps

- [Routing](./guide/routing.md) covers HTTP methods, path patterns and precedence.
- [Route Handlers](./guide/handlers.md) explains how to choose the input and output types.
- [Middleware](./guide/middleware.md) explains the chain in detail.
- [Error Handling](./guide/errors.md) shows how to report failures.
- [Recipes](./recipes.md) has ready to use snippets for CORS, authentication, panic recovery and more.
