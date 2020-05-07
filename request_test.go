package vox

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestJSONWithInvalidContentHeader(t *testing.T) {
	app := New()
	app.Use(func(ctx *Context, req *Request, res *Response) {
		data := make(map[string]interface{})
		if err := req.JSON(&data); err != nil {
			res.Body = "error"
		} else {
			res.Status = 200
			res.Body = data
		}
	})

	r := httptest.NewRequest("POST", "http://test.com/", strings.NewReader(`{"foo": "bar"}`))
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Result().StatusCode != 406 {
		t.Fail()
	}
	body, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		t.Fail()
	}
	if string(body) != "error" {
		t.Fail()
	}
}

func TestRequestJSONWithInvalidBody(t *testing.T) {
	app := New()
	app.Use(func(ctx *Context, req *Request, res *Response) {
		data := make(map[string]interface{})
		if err := req.JSON(&data); err != nil {
			res.Body = "error"
		} else {
			res.Status = 200
			res.Body = data
		}
	})

	r := httptest.NewRequest("POST", "http://test.com/", strings.NewReader(`INVALID JSON!`))
	r.Header.Set("content-type", "application/json; charset=utf-8")
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Result().StatusCode != 406 {
		t.Fail()
	}
	body, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		t.Fail()
	}
	if string(body) != "error" {
		t.Fail()
	}
}

func TestRequestJSON(t *testing.T) {
	app := New()
	app.Use(func(ctx *Context, req *Request, res *Response) {
		data := make(map[string]interface{})
		if err := req.JSON(&data); err != nil {
			res.Body = "error"
		} else {
			res.Status = 200
			res.Body = data
		}
	})

	r := httptest.NewRequest("POST", "http://test.com/", strings.NewReader(`{"foo": "bar"}`))
	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Result().StatusCode != 200 {
		t.Fail()
	}
	body, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		t.Fail()
	}
	if string(body) != `{"foo":"bar"}` {
		t.Fail()
	}
}
