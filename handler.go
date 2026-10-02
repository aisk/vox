package vox

// Handler is a middleware function.
type Handler func(*Context, *BaseRequest, *BaseResponse)

// RouteHandler handles a request with input and output body types.
type RouteHandler[In, Out any] func(*Context, *Request[In], *Response[Out])
