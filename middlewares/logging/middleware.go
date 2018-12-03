package logging

import (
	"fmt"
	"time"

	"github.com/aisk/vox"
)

// Middleware is logging middleware for vox.
func Middleware(req *vox.Request, res *vox.Response) {
	req.Next()
	username := "-"
	if req.URL.User != nil {
		if name := req.URL.User.Username(); name != "" {
			username = name
		}
	}
	fmt.Printf("%s - %s [%s] \"%s %s %s\" %d\n", req.RemoteAddr, username, time.Now().Format("02/Jan/2006:15:04:05 -0700"), req.Method, req.URL.Path, req.Proto, res.Status)
}
