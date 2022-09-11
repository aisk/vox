package vox

import (
	"fmt"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	written int
}

func (writer *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := writer.ResponseWriter.Write(b)
	writer.written += n
	return n, err
}

func logging(ctx *Context, req *Request, res *Response) {
	if ctx.App.GetConfig("logging:disable") != "" {
		ctx.Next()
		return
	}

	writer := &loggingResponseWriter{res.Writer, 0}
	res.Writer = writer
	ctx.Next()
	username := "-"
	if req.URL.User != nil {
		if name := req.URL.User.Username(); name != "" {
			username = name
		}
	}
	fmt.Printf("%s - %s [%s] \"%s %s %s\" %d %d\n", req.RemoteAddr, username, time.Now().Format("02/Jan/2006:15:04:05 -0700"), req.Method, req.URL.Path, req.Proto, res.Status, writer.written)
}
