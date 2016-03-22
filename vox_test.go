package vox

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEmpty(t *testing.T) {
	app := New()
	ts := httptest.NewServer(app)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	defer res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 404 {
		t.Error("wrong status code")
	}
}

func TestBasicMiddlewareResponse(t *testing.T) {
	app := New()
	app.Use(func(req *Request, next Next) *Response {
		return NewResponse("foo")
	})
	ts := httptest.NewServer(app)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	defer res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200 {
		t.Error("wrong status code")
	}
	content, _ := ioutil.ReadAll(res.Body)
	if string(content) != "foo" {
		t.Error("wrong response body")
	}
}

func TestMiddlewaresOrder(t *testing.T) {
	var order []int

	app := New()
	app.Use(func(req *Request, next Next) *Response {
		order = append(order, 0)
		res := next()
		order = append(order, 2)
		return res
	})
	app.Use(func(req *Request, next Next) *Response {
		order = append(order, 1)
		return NewResponse("foo")
	})
	app.Use(func(req *Request, next Next) *Response {
		// this should not be called because 'next' function is not called in prev middleware
		order = append(order, 3)
		return NewResponse("baz")
	})
	ts := httptest.NewServer(app)
	defer ts.Close()

	res, _ := http.Get(ts.URL)
	defer res.Body.Close()

	if len(order) != 3 || order[0] != 0 || order[1] != 1 || order[2] != 2 {
		t.Error("wrong order")
	}
}
