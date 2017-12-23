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
	app.Use(func(ctx *vox.Context, next func()) {
		start := time.Now()
		next()
		ms := time.Now().Sub(start).Seconds() / 1000
		ctx.Response.Header.Set("X-Response-Time", fmt.Sprintf("%fms", ms))
	})

	// logger
	app.Use(func(ctx *vox.Context, next func()) {
		next()
		fmt.Printf("%s %s\n", ctx.Request.Method, ctx.Request.URL)
	})

	// router param
	app.Get(regexp.MustCompile(`/hello/(?P<name>\w+)`), func(ctx *vox.Context) {
		ctx.Response.SetBody("Hello, " + ctx.Request.Params["name"] + "!")
	})

	// response
	app.Get(regexp.MustCompile("/"), func(ctx *vox.Context) {
		ctx.Response.SetBody("Hello, World!")
	})

	app.Run(":3000")
}
