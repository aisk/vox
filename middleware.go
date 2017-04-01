package vox

type Next func()

type MiddleWare func(*Request, Next) *Response
