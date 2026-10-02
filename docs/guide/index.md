# Guide

The guide explains each part of Vox in depth. If you are new to Vox, read [Getting Started](../getting-started.md) first, then the first three pages below. The remaining pages can be read in any order.

| Page | What it covers |
|:--|:--|
| [Routing](./routing.md) | HTTP methods, path patterns, parameters and precedence |
| [Route Handlers](./handlers.md) | The handler signature and how to choose input and output types |
| [Middleware](./middleware.md) | How the chain works and how to write your own middleware |
| [Context](./context.md) | Passing values, deadlines and cancellation |
| [Request](./request.md) | Reading request data, JSON decoding rules, forms and uploads |
| [Response](./response.md) | Body types, status codes, headers, redirects and streaming |
| [Error Handling](./errors.md) | Error responses, decode errors, custom 404 pages and panics |
| [Configuration](./configuration.md) | Built-in settings and your own |
| [Running](./running.md) | Serving, graceful shutdown and integration with `net/http` |
| [Testing](./testing.md) | Testing routes and middleware with `httptest` |
