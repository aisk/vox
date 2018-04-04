# VOX

[![GoDoc](https://godoc.org/github.com/aisk/vox?status.svg)](https://godoc.org/github.com/aisk/vox) [![Build Status](https://travis-ci.org/aisk/vox.svg?branch=master)](https://travis-ci.org/aisk/vox) [![Codecov](https://img.shields.io/codecov/c/github/aisk/vox.svg)](https://codecov.io/gh/aisk/vox)

A golang web framework for humans, inspired by [Koa](http://koajs.com) heavily.

![](https://i.v2ex.co/9MO3sMs4.jpeg)

---

## Example

```go
package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/aisk/vox"
)

func main() {
	app := vox.New()

	// x-response-time
	app.Use(func(req *vox.Request, res *vox.Response) {
		start := time.Now()
		req.Next()
		ms := time.Now().Sub(start).Seconds() / 1000
		res.Header.Set("X-Response-Time", fmt.Sprintf("%fms", ms))
	})

	// logger
	app.Use(func(req *vox.Request, res *vox.Response) {
		req.Next()
		fmt.Printf("%s %s\n", req.Method, req.URL)
	})

	// router param
	app.Get(regexp.MustCompile(`/hello/(?P<name>\w+)`), func(req *vox.Request, res *vox.Response) {
		res.Body = "Hello, " + req.Params["name"] + "!"
	})

	// response
	app.Get(regexp.MustCompile("/"), func(req *vox.Request, res *vox.Response) {
		// get the query string
		name := req.URL.Query().Get("name")
		if name == "" {
			name = "World"
		}
		res.Body = "Hello, " + name + "!"
	})

	app.Run(":3000")
}
```
