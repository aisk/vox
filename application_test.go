package vox

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestEmptyApplication(t *testing.T) {
	app := New()
	r := httptest.NewRequest("GET", "http://test.com/", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Result().StatusCode != 404 {
		t.Fail()
	}
}

func TestBasicApplication(t *testing.T) {
	app := New()
	app.Use(func(ctx *Context) {
		if ctx.Request.Method != "GET" {
			t.Fail()
		}
		if ctx.Request.URL.String() != "http://test.com/" {
			t.Fail()
		}
		ctx.Response.SetStatus(200)
		ctx.Response.SetBody("Hello Vox!")
		ctx.Response.Header.Set("foo", "bar")
	})
	r := httptest.NewRequest("GET", "http://test.com/", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Result().StatusCode != 200 {
		t.Fail()
	}
	if w.Result().Header.Get("foo") != "bar" {
		t.Fail()
	}
	body, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		t.Fail()
	}
	if string(body) != "Hello Vox!" {
		t.Fail()
	}
}
