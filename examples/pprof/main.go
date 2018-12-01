package main

import (
	"github.com/aisk/vox"
	"github.com/aisk/vox/middlewares/pprof"
)

func main() {
	app := vox.New()
	app.Use(pprof.Middleware)

	app.Get("/", func(req *vox.Request, res *vox.Response) {
		res.Body = "Hello, World!"
	})

	app.Run("localhost:3000")
}
