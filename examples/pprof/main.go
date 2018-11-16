package main

import (
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"

	"github.com/aisk/vox"
)

func main() {
	app := vox.New()

	app.Get("/", func(req *vox.Request, res *vox.Response) {
		res.Body = "Hello, World!"
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/", app)
	http.ListenAndServe(":3000", mux)
}
