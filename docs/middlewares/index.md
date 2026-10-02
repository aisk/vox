---
title: Bundled Middleware
nav_order: 4
---

# Bundled Middleware

Vox ships with a few middleware. Two of them are built into every application, and the others are optional packages in the same module.

| Middleware | Package | Purpose |
|:--|:--|:--|
| `logging` | built in | Prints an access log line for each request. |
| `respond` | built in | Writes the status, headers and body to the client. |
| [Static Files]({% link middlewares/static.md %}) | `github.com/aisk/vox/middlewares/static` | Serves files from a local directory. |
| [Pprof]({% link middlewares/pprof.md %}) | `github.com/aisk/vox/middlewares/pprof` | Exposes Go's runtime profiling endpoints. |

The built-in middleware is described in the [Middleware guide]({% link guide/middleware.md %}#built-in-middleware). To write your own, start with the same guide, and see [Recipes]({% link recipes.md %}) for common examples.
