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
	app.Use(func(req *Request, res *Response) {
		if req.Method != "GET" {
			t.Fail()
		}
		if req.URL.String() != "http://test.com/" {
			t.Fail()
		}
		res.Status = 200
		res.Body = "Hello Vox!"
		res.Header.Set("foo", "bar")
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
