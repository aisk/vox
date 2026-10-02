package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/aisk/vox"
)

type CreateUser struct {
	Name string `json:"name"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	app := vox.New()

	// custom middleware that add a x-response-time to the response header
	app.Use(func(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
		start := time.Now()
		ctx.Next()
		duration := time.Now().Sub(start)
		res.Header.Set("X-Response-Time", fmt.Sprintf("%s", duration))
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
		res.Body = `
		<!doctype html>
			<head>
				<title>xxx</title>
			</head>
			<body>
				<p>Hello, ` + name + `!</p>
			</body>
		</html>
		`
	})

	// error as body
	app.Get("/error", func(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[error]) {
		res.Body = errors.New("Error!")
	})

	app.Post("/users", func(ctx *vox.Context, req *vox.Request[CreateUser], res *vox.Response[User]) {
		res.Status = 201
		res.Body = User{ID: 1, Name: req.Body.Name}
	})

	app.Run("[::]:3000")
}
