---
title: Route Handlers
parent: Guide
nav_order: 2
---

# Route Handlers
{: .no_toc }

A route handler is a function that receives a typed request and fills in a typed response.

1. TOC
{:toc}

---

## Signature

```go
func(ctx *vox.Context, req *vox.Request[In], res *vox.Response[Out])
```

- `ctx` is the [Context]({% link guide/context.md %}) of the current request.
- `req` is the [Request]({% link guide/request.md %}). `req.Body` has type `In`.
- `res` is the [Response]({% link guide/response.md %}). `res.Body` has type `Out`.

A handler has no return value. It reports its result by assigning to `res` and returning.

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

`In` and `Out` are inferred when the handler is registered, so `app.Post("/users", createUser)` needs no type arguments.

## Choosing the input type

| `In` | Behavior |
|:--|:--|
| `vox.NoBody` | The request body is not read. Use it for routes without a body and for bodies you read yourself. |
| A struct, map, slice or other JSON compatible type | The body is decoded from JSON before the handler runs. |
| A pointer to one of those | Same, and `req.Body` is never nil when the handler runs. |

{: .warning }
Any input type other than `vox.NoBody` makes a JSON body mandatory. A request without one is rejected with 415 or 400, whatever its method. Routes that take no body, which includes most `GET` and `DELETE` routes, should use `Request[vox.NoBody]`.

`vox.NoBody` is also the way to accept forms, uploads and other non JSON bodies. See [Request]({% link guide/request.md %}#forms-and-uploads).

## Choosing the output type

| `Out` | Response |
|:--|:--|
| `string`, `[]byte` | Written as is. |
| `io.Reader`, `io.ReadCloser` | Streamed to the client. A closer is closed afterwards. |
| `error` | The error message, with status 500 unless you set one. |
| `vox.NoBody` | An empty body, with status 204 unless you set one. |
| Anything else | Encoded as JSON with `Content-Type: application/json`. |
| `any` | Decided at run time by the value you assign, following the rows above. |

The output type describes the successful response. Failures do not have to fit it, as [Error Handling]({% link guide/errors.md %}) explains.

## Lifecycle

For a request that matches a route, Vox does the following:

1. Middleware runs up to its `ctx.Next()` call, in order.
2. The router matches the path and fills `req.Params`.
3. Unless `In` is `vox.NoBody`, the body is decoded. On failure the handler is skipped and an error response is prepared.
4. The handler runs.
5. `res.Body` is committed to the response, so middleware can see it.
6. Middleware resumes after `ctx.Next()`, in reverse order.
7. The status, headers and body are written to the client.

Because nothing is written until step 7, a handler can set the status, headers and body in any order, and middleware can still change all three afterwards.

## Handlers with dependencies

Handlers often need a database handle or another shared dependency. Closures and methods both work.

With a closure:

```go
func showUser(db *sql.DB) vox.RouteHandler[vox.NoBody, User] {
	return func(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[User]) {
		var user User
		err := db.QueryRowContext(ctx, "SELECT id, name FROM users WHERE id = ?", req.Params["id"]).Scan(&user.ID, &user.Name)
		if err != nil {
			res.Status = 404
			return
		}
		res.Body = user
	}
}

app.Get("/users/{id}", showUser(db))
```

With methods on a struct:

```go
type UserHandlers struct {
	DB *sql.DB
}

func (h *UserHandlers) Show(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[User]) {
	// use h.DB
}

users := &UserHandlers{DB: db}
app.Get("/users/{id}", users.Show)
```

Values that change per request, such as the authenticated user, belong in the [Context]({% link guide/context.md %}#passing-values).

## Do not call Next

`ctx.Next()` is for middleware. A route handler is the end of the chain and returns when it is done. Calling `ctx.Next()` from a route handler panics.

## Explicit type arguments

The type arguments can also be written out. This is required when you refer to a registration method without calling it, for example to store it in a variable:

```go
register := app.Get[vox.NoBody, string]
register("/health", healthHandler)
register("/ready", readyHandler)
```
