---
title: Static Files
parent: Bundled Middleware
nav_order: 1
---

# Static Files
{: .no_toc }

The `static` package serves files from a local directory under a URL prefix.

1. TOC
{:toc}

---

## Usage

```go
package main

import (
	"github.com/aisk/vox"
	"github.com/aisk/vox/middlewares/static"
)

func main() {
	app := vox.New()

	// GET /assets/logo.png -> ./public/logo.png
	app.Use(static.Middleware("/assets", "./public"))

	app.Run("localhost:3000")
}
```

`static.Middleware(prefix, root)` takes the URL prefix and the directory to serve. A relative `root` is resolved against the working directory of the process.

## How requests are resolved

The middleware calls `ctx.Next()` first and only serves a file when nothing else has responded. It handles a request when all of the following hold:

- The method is `GET` or `HEAD`.
- The path is the prefix itself or lies below it. With the prefix `/assets`, both `/assets` and `/assets/css/site.css` qualify, and `/assetsx` does not.
- No route or other middleware has set a status or a body.

In consequence, **routes take precedence over files**. If `/assets/logo.png` is both a route and a file, the route wins.

Files are served by Go's [`http.FileServer`](https://pkg.go.dev/net/http#FileServer), which brings its usual behavior:

- The content type is derived from the file extension, and conditional and range requests are supported.
- A request for a directory serves its `index.html`. Without one, a listing of the directory is returned.
- A request for a directory without a trailing slash, or for `index.html` itself, is redirected to the canonical directory URL.
- A file that does not exist yields `404 page not found`.

## Prefix

The prefix is normalized, so `"assets"`, `"/assets"` and `"/assets/"` are equivalent. An empty prefix or `"/"` serves the directory at the root of the site:

```go
app.Use(static.Middleware("/", "./public"))
```

Since routes take precedence, this can be combined with API routes. Requests that match a route are handled by it, and everything else is looked up in the directory.

## Security notes

- Every file below `root` is reachable, including dot files such as `.env` or `.git`. Point `root` at a directory that contains only public files.
- Directory listings are enabled for directories without an `index.html`.
- Paths cannot escape `root` with `..` segments.
