package vox

import "net/http"

type Cookie struct {
	rq *http.Request
	rw http.ResponseWriter
}

func (*Cookie) Set() {}

func (*Cookie) Get() {}

func (*Cookie) Add() {}

func (*Cookie) Del() {}

func createCookies(rq *http.Request, rw http.ResponseWriter) *Cookie {
	return &Cookie{
		rq: rq,
		rw: rw,
	}
}
