---
title: Run
nav_order: 3
---

# Run Your Application
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Run

You can run your vox application by simply calling the `Run` method:

```go
app := vox.New()
app.Run("localhost:3000")
```

Now your application is listening on port 3000. The `Run` method accepts the same arguments as [`net/http.ListenAndServe`](https://golang.org/pkg/net/http/#ListenAndServe).

## Integrate with an existing HTTP server

If you already have an HTTP server in Go, you can integrate vox with it. This can help you migrate to and from vox.

Actually, `vox.Application` implements the [`net/http.Handler`](https://golang.org/pkg/net/http/#Handler) interface. So you can pass a `vox.Application` instance to any function that accepts a [`net/http.Handler`](https://golang.org/pkg/net/http/#Handler), like `http.Handle`:

```go
func rawHandler(w http.ResponseWriter, _ *http.Request) {
    io.WriteString(w, "Hello from a raw handler")
}

func voxHandler(ctx *vox.Context, req *vox.Request, res *vox.Response) {
    res.Body = "Hello from a vox handler"
}

func main() {
    var app = vox.New()
    app.Get("/vox", voxHandler)
    http.HandleFunc("/raw", rawHandler)
    http.Handle("/", app)
    http.ListenAndServe("localhost:3000", nil)
}
```