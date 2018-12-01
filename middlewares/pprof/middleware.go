package pprof

import (
	"net/http/pprof"
	"strings"

	"github.com/aisk/vox"
)

// Middleware is vox's pprof middleware.
func Middleware(req *vox.Request, res *vox.Response) {
	if strings.HasPrefix(req.URL.Path, "/debug/pprof/cmdline") {
		res.DontRespond = true
		pprof.Cmdline(res.Writer, req.Request)
		return
	}
	if strings.HasPrefix(req.URL.Path, "/debug/pprof/profile") {
		res.DontRespond = true
		pprof.Profile(res.Writer, req.Request)
		return
	}
	if strings.HasPrefix(req.URL.Path, "/debug/pprof/symbol") {
		res.DontRespond = true
		pprof.Symbol(res.Writer, req.Request)
		return
	}
	if strings.HasPrefix(req.URL.Path, "/debug/pprof/trace") {
		res.DontRespond = true
		pprof.Trace(res.Writer, req.Request)
		return
	}
	if strings.HasPrefix(req.URL.Path, "/debug/pprof") {
		res.DontRespond = true
		pprof.Index(res.Writer, req.Request)
		return
	}
	req.Next()
}
