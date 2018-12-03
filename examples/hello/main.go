package main

import (
	"fmt"
	"time"

	"github.com/aisk/vox"
	"github.com/aisk/vox/middlewares/logging"
)

func main() {
	app := vox.New()

	// logging
	app.Use(logging.Middleware)

	// custom middleware that add a x-response-time to the response header
	app.Use(func(req *vox.Request, res *vox.Response) {
		start := time.Now()
		req.Next()
		duration := time.Now().Sub(start)
		res.Header.Set("X-Response-Time", fmt.Sprintf("%s", duration))
	})

	// router param
	app.Get("/hello/{name}", func(req *vox.Request, res *vox.Response) {
		res.Body = "Hello, " + req.Params["name"] + "!"
	})

	// response
	app.Get("/", func(req *vox.Request, res *vox.Response) {
		// get the query string
		name := req.URL.Query().Get("name")
		if name == "" {
			name = "World"
		}
		res.Body = "Hello, " + name + "!"
	})

	app.Run("localhost:3000")
}
