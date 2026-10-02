---
title: Pprof
parent: Bundled Middleware
nav_order: 2
redirect_from:
  - /docs/pprof
---

# Pprof
{: .no_toc }

The `pprof` package exposes the profiling endpoints of Go's [`net/http/pprof`](https://pkg.go.dev/net/http/pprof) through a Vox middleware.

1. TOC
{:toc}

---

## Usage

```go
package main

import (
	"github.com/aisk/vox"
	"github.com/aisk/vox/middlewares/pprof"
)

func main() {
	app := vox.New()
	app.Use(pprof.Middleware)

	app.Get("/", func(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
		res.Body = "Hello, World!"
	})

	app.Run("localhost:3000")
}
```

## Endpoints

The middleware handles every path that starts with `/debug/pprof`:

| Path | Content |
|:--|:--|
| `/debug/pprof/` | An index page listing the available profiles |
| `/debug/pprof/profile` | A CPU profile, 30 seconds by default (`?seconds=N`) |
| `/debug/pprof/heap` | A sample of live heap allocations |
| `/debug/pprof/allocs` | A sample of all past allocations |
| `/debug/pprof/goroutine` | Stack traces of all goroutines |
| `/debug/pprof/block` | Stack traces that led to blocking |
| `/debug/pprof/mutex` | Stack traces of contended mutex holders |
| `/debug/pprof/threadcreate` | Stack traces that created OS threads |
| `/debug/pprof/trace` | An execution trace (`?seconds=N`) |
| `/debug/pprof/cmdline` | The command line of the process |
| `/debug/pprof/symbol` | Symbol lookup for program counters |

Use them with `go tool pprof`:

```sh
# CPU profile over 30 seconds
go tool pprof http://localhost:3000/debug/pprof/profile

# Heap profile
go tool pprof http://localhost:3000/debug/pprof/heap

# Execution trace over 5 seconds
curl -o trace.out 'http://localhost:3000/debug/pprof/trace?seconds=5'
go tool trace trace.out
```

## How it works

For pprof paths, the middleware sets `res.DontRespond = true` and lets the standard pprof handlers write to the raw `http.ResponseWriter`. It does not call `ctx.Next()` for them, so routes and later middleware are skipped. All other requests pass through untouched.

## Restricting access

Profiling endpoints expose internals of your program, and a CPU profile or trace costs resources while it runs. In production, do not leave them open to the public.

Middleware registered before `pprof.Middleware` wraps it, so a guard placed first can reject requests before they reach the profiler:

```go
func guardPprof(token string) vox.Handler {
	return func(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
		if strings.HasPrefix(req.URL.Path, "/debug/pprof") {
			given := req.Header.Get("X-Debug-Token")
			if token == "" || subtle.ConstantTimeCompare([]byte(given), []byte(token)) != 1 {
				res.Status = 403
				return
			}
		}
		ctx.Next()
	}
}

app.Use(guardPprof(os.Getenv("DEBUG_TOKEN")))
app.Use(pprof.Middleware)
```

Pass the token along when profiling:

```sh
curl -H "X-Debug-Token: $DEBUG_TOKEN" -o heap.out http://localhost:3000/debug/pprof/heap
go tool pprof heap.out
```

Another common approach is to serve pprof from a second application bound to a private address:

```go
debug := vox.New()
debug.Use(pprof.Middleware)
go debug.Run("127.0.0.1:6060")
```
