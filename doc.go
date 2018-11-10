/*
Package vox is a golang web framework for humans, inspired by Koa http://koajs.com heavily.

Introduction

Vox is a web framework inspired by Koa, which aims to be a minimal and elegent library for web applications.

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

		// x-response-time
		app.Use(func(req *vox.Request, res *vox.Response) {
			start := time.Now()
			req.Next()
			duration := time.Now().Sub(start)
			res.Header.Set("X-Response-Time", fmt.Sprintf("%s", duration))
		})

		// logger
		app.Use(func(req *vox.Request, res *vox.Response) {
			req.Next()
			fmt.Printf("%s %s\n", req.Method, req.URL)
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

		app.Run(":3000")
	}


*/
package vox
