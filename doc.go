/*
Package vox is a Go web framework for humans, heavily inspired by Koa http://koajs.com.

# Introduction

Vox is a web framework inspired by Koa, which aims to be a minimal and elegant library for web applications.

Installation

	$ go get -u github.com/aisk/vox

Basic Example

	package main

	import (
		"fmt"
		"time"

		"github.com/aisk/vox"
	)

	func main() {
		app := vox.New()

		// X-Response-Time
		app.Use(func(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
			start := time.Now()
			ctx.Next()
			duration := time.Now().Sub(start)
			res.Header.Set("X-Response-Time", fmt.Sprintf("%s", duration))
		})

		// logger
		app.Use(func(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
			ctx.Next()
			fmt.Printf("%s %s\n", req.Method, req.URL)
		})

		// router param
		app.Get("/hello/{name}", func(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
			res.Body = "Hello, " + req.Params["name"] + "!"
		})

		// response
		app.Get("/", func(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
			// get the query string
			name := req.URL.Query().Get("name")
			if name == "" {
				name = "World"
			}
			res.Body = "Hello, " + name + "!"
		})

		app.Run("localhost:3000")
	}
*/
package vox
