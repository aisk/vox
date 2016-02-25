package vox

import (
	"testing"
)

func TestNewResponse(t *testing.T) {
	var res *Response

	res = NewResponse("xxx")
	if string(res.Content) != "xxx" {
		t.Error()
	}
	if res.StatusCode != 200 {
		t.Error()
	}

	res = NewResponse("xxx", 123)
	if string(res.Content) != "xxx" {
		t.Error()
	}
	if res.StatusCode != 123 {
		t.Error()
	}

	res = NewResponse("xxx", 123, map[string]string{"foo": "bar"})
	if string(res.Content) != "xxx" {
		t.Error()
	}
	if res.StatusCode != 123 {
		t.Error()
	}
	if res.Header["foo"] != "bar" {
		t.Error()
	}

	res = NewResponse("xxx", map[string]string{"foo": "bar"}, 123)
	if string(res.Content) != "xxx" {
		t.Error()
	}
	if res.StatusCode != 123 {
		t.Error()
	}
	if res.Header["foo"] != "bar" {
		t.Error()
	}
}
