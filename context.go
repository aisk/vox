package vox

import "context"

// Context is vox's handler context, which implemented the context.Context interface
// by wrapping the context from current requests'. You should use it when you need
// the standard context.Context.
type Context struct {
	context.Context
	// Next will call the next handler / middleware to processing request.
	// It's the middleware's responsibility to call the Next function (or not).
	Next func()
}
