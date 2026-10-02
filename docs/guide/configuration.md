---
title: Configuration
parent: Guide
nav_order: 8
---

# Configuration
{: .no_toc }

An application carries a small set of string settings. Vox uses them for its own options, and you can use them for yours.

1. TOC
{:toc}

---

## Setting and reading values

```go
app := vox.New()
app.SetConfig("request:max-body-size", "4194304")

limit := app.GetConfig("request:max-body-size") // "4194304"
```

Keys and values are strings. `GetConfig` returns an empty string for a key that was never set.

Configure the application before it starts serving. Settings are not synchronized, so changing them while requests are being handled is a data race.

## Built-in settings

| Key | Default | Meaning |
|:--|:--|:--|
| `request:max-body-size` | `1048576` | The maximum size in bytes of a JSON request body. Larger bodies are rejected with 413. `0` disables the limit. |
| `logging:disable` | unset | Any non-empty value turns off the built-in access log. |

`request:max-body-size` must be a non-negative integer. `SetConfig` panics on any other value, so a typo is caught at startup. The limit applies to automatic JSON decoding and to `BaseRequest.JSON`. See [Request]({% link guide/request.md %}#body-size-limit).

`logging:disable` only checks whether the value is empty, so `"false"` disables the log as well. Leave the key unset to keep logging on.

```go
app.SetConfig("logging:disable", "true")
```

## Your own settings

Any other key is stored as is. Handlers and middleware reach the application through the context:

```go
app.SetConfig("app:version", "1.4.2")

app.Get("/version", func(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
	res.Body = ctx.App.GetConfig("app:version")
})
```

Prefix your keys, as in `app:version`, to keep them apart from the built-in ones.

Settings are a convenience for simple string values. For structured configuration or dependencies such as a database handle, pass them to your handlers directly, as shown in [Route Handlers]({% link guide/handlers.md %}#handlers-with-dependencies).

## Reading from the environment

Vox does not read environment variables or files. Load the values yourself and pass them on:

```go
app := vox.New()
if size := os.Getenv("MAX_BODY_SIZE"); size != "" {
	app.SetConfig("request:max-body-size", size)
}
if os.Getenv("ACCESS_LOG") == "off" {
	app.SetConfig("logging:disable", "true")
}
```
