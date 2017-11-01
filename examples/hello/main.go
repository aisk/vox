package main

import (
	"github.com/aisk/vox"
)

func main() {
	app := vox.New()
	app.Get("/hello", func(ctx *vox.Context) {
		ctx.Response.SetBody("Hello, World!")
	})
	app.Run("localhost:3000")
}
