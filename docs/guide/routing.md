# Routing

A route binds an HTTP method and a path pattern to a handler.

## Registering routes

Each common HTTP method has a helper on the application:

```go
app.Get("/towels", listTowels)
app.Post("/towels", createTowel)
app.Put("/towels/{id}", replaceTowel)
app.Patch("/towels/{id}", updateTowel)
app.Delete("/towels/{id}", deleteTowel)
app.Head("/towels", headTowels)
app.Options("/towels", towelOptions)
app.Trace("/towels", traceTowels)
```

The helpers are shorthand for `Route`, which takes the method as its first argument. Use it for methods without a helper, including ones you make up:

```go
app.Route("PURGE", "/cache/{key}", purgeCache)
```

Method names are matched exactly, so write them in upper case.

### Matching any method

Pass `"*"` as the method to handle every method on a path:

```go
app.Route("*", "/health", func(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
	res.Body = "ok: " + req.Method
})
```

A route registered for a specific method takes precedence over a `"*"` route on the same path. This suits generic endpoints such as health probes. For business APIs, explicit methods are usually clearer.

### HEAD requests

A `HEAD` request is served by the matching `Get` route when no `Head` route is registered for the path. The handler runs as usual and Vox omits the response body.

## Path patterns

Paths use the pattern syntax of Go's [`http.ServeMux`](https://pkg.go.dev/net/http#ServeMux).

| Pattern | Matches | Does not match |
|:--|:--|:--|
| `/users` | `/users` | `/users/` |
| `/users/{id}` | `/users/42` | `/users/`, `/users/42/`, `/users/42/posts` |
| `/files/{path...}` | `/files/a/b.txt`, `/files/` | `/files` |
| `/static/` | `/static/`, `/static/css/site.css` | `/static` |
| `/{$}` | `/` | `/anything` |

- `{name}` matches exactly one path segment.
- `{name...}` matches the rest of the path, including slashes. It must be the last segment.
- A pattern ending in `/` matches every path below it. In particular `/` matches all paths, so once it is registered no request gets a 404.
- `{$}` at the end of a pattern matches only the trailing slash itself. Use `/{$}` for a home page.

## Path parameters

Wildcard values are available in `req.Params`, keyed by name:

```go
app.Get("/users/{id}/posts/{slug}", func(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
	res.Body = "user " + req.Params["id"] + ", post " + req.Params["slug"]
})

app.Get("/files/{path...}", func(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
	// GET /files/a/b.txt gives "a/b.txt"
	res.Body = req.Params["path"]
})
```

Values are always strings and are already percent-decoded. Convert them yourself, and treat a failed conversion as a client error:

```go
app.Get("/users/{id}", func(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[User]) {
	id, err := strconv.Atoi(req.Params["id"])
	if err != nil {
		res.Status = 400
		return
	}
	res.Body = User{ID: id}
})
```

Matching uses the decoded path, so an encoded slash (`%2F`) separates segments just like a literal one.

## Query strings

The query string is not part of route matching. Read it from the URL:

```go
app.Get("/search", func(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
	query := req.URL.Query()
	term := query.Get("q") // "" when absent
	tags := query["tag"]   // every value of a repeated key
	res.Body = fmt.Sprintf("q=%s tags=%v", term, tags)
})
```

## Precedence

When several patterns match a request, the most specific one wins, regardless of the order in which routes were registered:

```go
app.Get("/users/{id}", showUser)
app.Get("/users/me", showCurrentUser) // handles /users/me
```

A pattern is more specific than another when it matches a strict subset of the other's paths. Avoid pairs where neither is more specific, such as `/a/{x}/c` and `/a/b/{y}`, which both match `/a/b/c`. Vox does not reject them, and which one handles the overlap is an implementation detail.

## Unmatched requests

When no route matches, the response is `404 Not Found`, unless a middleware provides a response. See [Error Handling](./errors.md#not-found) to customize it.

Vox keeps matching simple, which makes it differ from `http.ServeMux` in a few places:

- A path that exists only for other methods yields 404, not 405.
- There is no redirect from `/static` to `/static/`, and paths are not cleaned. `/users//42` and `/users/42` are different paths.
- Patterns cannot include a host name.

## Registration errors

Registration panics when the pattern is invalid, for example an unclosed `{` or a path that does not start with `/`, and when the same method and path are registered twice. Routes are normally registered at startup, so these mistakes surface immediately.

## Organizing routes

Routes are plain method calls, so ordinary functions are enough to group them:

```go
func registerUserRoutes(app *vox.Application, store *UserStore) {
	app.Get("/users", store.list)
	app.Post("/users", store.create)
	app.Get("/users/{id}", store.show)
}

func main() {
	app := vox.New()
	registerUserRoutes(app, NewUserStore())
	app.Run("localhost:3000")
}
```

Handlers can be functions, closures or methods. [Route Handlers](./handlers.md#handlers-with-dependencies) shows how to give them access to shared dependencies.

Routes can be registered before or after `app.Use` calls. The router always runs after all middleware, as described in [Middleware](./middleware.md#execution-order).
