package vox

type Next func() *Response

type MiddleWare func(*Request, Next) *Response
