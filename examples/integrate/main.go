package main

import (
	"io"
	"net/http"

	"github.com/aisk/vox"
)

func rawHandler(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "Hello from a raw handler")
}

func voxHandler(ctx *vox.Context, req *vox.Request, res *vox.Response) {
	res.Body = "Hello from a vox handler"
}

func main() {
	var app = vox.New()
	app.Get("/vox", voxHandler)
	http.HandleFunc("/raw", rawHandler)
	http.Handle("/", app)
	http.ListenAndServe("localhost:3000", nil)
}
