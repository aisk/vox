---
title: Usage
nav_order: 3
has_toc: true
---

# Usage

## Quick review

```go
package main

import (
	"fmt"
	"time"

	"github.com/aisk/vox"
)

func main() {
	app := vox.New()

	// custom middleware that add a x-response-time to the response header
	app.Use(func(ctx *vox.Context, *vox.Request, res *vox.Response) {
		start := time.Now()
		ctx.Next()
		duration := time.Now().Sub(start)
		res.Header.Set("X-Response-Time", fmt.Sprintf("%s", duration))
	})

	// router param
	app.Get("/hello/{name}", func(ctx *vox.Context, *vox.Request, res *vox.Response) {
		res.Body = "Hello, " + req.Params["name"] + "!"
	})

	// response
	app.Get("/", func(ctx *vox.Context, *vox.Request, res *vox.Response) {
		// get the query string
		name := req.URL.Query().Get("name")
		if name == "" {
			name = "World"
		}
		res.Body = "Hello, " + name + "!"
	})

	app.Run("localhost:3000")
}
```

## Handle HTTP Methods

```go
package main

import (
	"github.com/aisk/vox"
)

func handler(ctx *vox.Context, *vox.Request, res *vox.Response) {
	// Get the current request's HTTP method and put it to the result page.
	res.Body = "HTTP Method is: " + req.Method
}

func main() {
	app := vox.New()

	app.Get("/", handler)
	app.Post("/", handler)
	app.Put("/", handler)
	app.Delete("/", handler)
	app.Head("/", handler)
	app.Options("/", handler)
	app.Trace("/", handler)

	// In some case you need handle custom HTTP method that not in the RFCs like FLY.
	app.Route("FLY", "/", handler)

	app.Run("localhost:3000")
}
```

## Get route parameters in URL path

```go
package main

import (
	"github.com/aisk/vox"
)

func hello(ctx *vox.Context, *vox.Request, res *vox.Response) {
	name := req.Params["name"]
	res.Body = "Hello, " + name + "!"
}

func main() {
	app := vox.New()
	app.Get("/hello/{name}", hello)
	app.Run("localhost:3000")
}
```

## Get querystring parameters in URL

```go
package main

import (
	"github.com/aisk/vox"
)

func hello(ctx *vox.Context, *vox.Request, res *vox.Response) {
	name := req.URL.Query().Get("name")
	res.Body = "Hello, " + name + "!"
}

func main() {
	app := vox.New()
	app.Get("/hello", hello)
	app.Run("localhost:3000")
}
```

## Set response data

```go
package main

import (
	"github.com/aisk/vox"
)

func towel(ctx *vox.Context, *vox.Request, res *vox.Response) {
	// Set the response body, it can be string or []byte or any thing that json.Marshal accepts.
	res.Body = "new towel is created!"
	// Set the response status code.
	res.Status = 201
	// Set the response header.
	res.Header.Set("Location", "/towels/42")
}

func main() {
	app := vox.New()
	app.Post("/towels", towel)
	app.Run("localhost:3000")
}
```

## Processing JSON request and send JSON response

```go
package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aisk/vox"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Towel struct {
	Color string `json:"color"`
	Size  string `json:"size"`
}

func towel(ctx *vox.Context, *vox.Request, res *vox.Response) {
	if !strings.HasPrefix(req.Header.Get("Content-Type"), "application/json") {
		res.Status = http.StatusUnsupportedMediaType // or just 415
		// Set the body with a map, vox will marshal it to JSON automatically for you.
		res.Body = map[string]interface{}{
			"code":    1,
			"message": "This is not a JSON request",
		}
		return
	}

	var t Towel
	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		res.Status = http.StatusUnprocessableEntity // or just 422
		res.Body = map[string]interface{}{
			"code": 1,
		}
	}

	// Set the body with a struct, vox will marshal it to JSON automatically for you.
	res.Body = t
	// Set the response status code.
	res.Status = 201
	// Set the response header.
	res.Header.Set("Location", "/towels/42")
}

func main() {
	app := vox.New()
	app.Post("/towels", towel)
	app.Run("localhost:3000")
}
```
