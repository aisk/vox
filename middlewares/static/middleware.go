package static

import (
	"net/http"
	"path"
	"strings"

	"github.com/aisk/vox"
)

// Middleware serves static files from root under prefix.
//
// For example: Middleware("/assets", "./public")
//
//	GET /assets/logo.png -> ./public/logo.png
func Middleware(prefix string, root string) vox.Handler {
	prefix = normalizePrefix(prefix)

	var fileServer http.Handler = http.FileServer(http.Dir(root))
	if prefix != "/" {
		fileServer = http.StripPrefix(prefix, fileServer)
	}

	return func(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
		ctx.Next()

		if req.Method != http.MethodGet && req.Method != http.MethodHead {
			return
		}
		if !matchPrefix(req.URL.Path, prefix) {
			return
		}
		if !isResponseUnwritten(res) {
			return
		}

		res.DontRespond = true
		fileServer.ServeHTTP(res.Writer, req.Request)
	}
}

func normalizePrefix(prefix string) string {
	if prefix == "" {
		return "/"
	}
	if prefix[0] != '/' {
		prefix = "/" + prefix
	}
	cleaned := path.Clean(prefix)
	if cleaned == "." {
		return "/"
	}
	return cleaned
}

func matchPrefix(requestPath string, prefix string) bool {
	if prefix == "/" {
		return true
	}
	return requestPath == prefix || strings.HasPrefix(requestPath, prefix+"/")
}

func isResponseUnwritten(res *vox.BaseResponse) bool {
	if res.DontRespond {
		return false
	}
	if res.Status != 0 {
		return false
	}
	return !res.HasBody()
}
