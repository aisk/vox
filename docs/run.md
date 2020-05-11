---
title: Run
nav_order: 3
---

# Run You Application
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Run

You can run your vox application by simply call the `Run` method:

```go
app := vox.New()
app.Run("localhost:3000")
```

Now your application has listened on port 3000 with localhost. `Run` accept whatever [net/http.ListenAndServe](https://golang.org/pkg/net/http/#ListenAndServe) takes in first parameter.

## Integrate with your exists HTTP server

If you have an http server in go already, you can integrate vox to it. This may help you migrate from and to vox.

Actually, `vox.Application` implemented the [net/http.Handler](https://golang.org/pkg/net/http/#Handler) interface. So you can pass a `vox.Application` instance to where [net/http.Handler](https://golang.org/pkg/net/http/#Handler) are accepted, like `http.Handle`:

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