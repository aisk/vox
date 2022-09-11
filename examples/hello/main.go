package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/aisk/vox"
)

func main() {
	app := vox.New()

	// custom middleware that add a x-response-time to the response header
	app.Use(func(ctx *vox.Context, req *vox.Request, res *vox.Response) {
		start := time.Now()
		ctx.Next()
		duration := time.Now().Sub(start)
		res.Header.Set("X-Response-Time", fmt.Sprintf("%s", duration))
	})

	// router param
	app.Get("/hello/{name}", func(ctx *vox.Context, req *vox.Request, res *vox.Response) {
		res.Body = "Hello, " + req.Params["name"] + "!"
	})

	// response
	app.Get("/", func(ctx *vox.Context, req *vox.Request, res *vox.Response) {
		// get the query string
		name := req.URL.Query().Get("name")
		if name == "" {
			name = "World"
		}
		res.Body = "Hello, " + name + "!"
	})

	// error as body
	app.Get("/error", func(ctx *vox.Context, req *vox.Request, res *vox.Response) {
		res.Body = errors.New("Error!")
	})

	app.Run("localhost:3000")
}
