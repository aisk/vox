package vox

import (
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type unsetBody struct{}

var (
	bodyNotSet   = unsetBody{}
	statusNotSet = 0
)

var htmlReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	`"`, "&#34;",
	"'", "&#39;",
)

func htmlEscape(s string) string {
	return htmlReplacer.Replace(s)
}

// A BaseResponse object contains all the information which will written to current
// HTTP client.
type BaseResponse struct {
	request    *BaseRequest
	redirected bool

	// Don't write headers, status and body from BaseResponse struct to client. In
	// the case of you're using the go's origin http.Response.
	DontRespond bool
	// Writer is the raw http.ResponseWriter for current request. You should
	// assign the Body / Status / Header value instead of using this field.
	Writer http.ResponseWriter
	// Body is the container for HTTP response's body.
	Body any
	// The status code which will respond as the HTTP response's status code.
	// 200 will be used as the default value if not set.
	Status int
	// Headers which will be written to the response.
	Header http.Header
}

func (response *BaseResponse) setImplicitContentType() {
	if response.Header.Get("Content-Type") != "" {
		return
	}

	if response.Body == bodyNotSet {
		return
	}

	if response.Status == 204 || response.Status == 304 {
		return
	}

	switch response.Body.(type) {
	case []byte, string, io.Reader, error:
	default:
		response.Header.Set("Content-Type", mime.TypeByExtension(".json"))
	}
}

var parseURL = url.Parse

// Redirect request to another url.
func (response *BaseResponse) Redirect(url string, code int) {
	response.redirected = true
	request := response.request

	if u, err := parseURL(url); err == nil {
		if u.Scheme == "" && u.Host == "" {
			oldpath := request.URL.Path
			if oldpath == "" {
				oldpath = "/"
			}

			if url == "" || url[0] != '/' {
				olddir, _ := path.Split(oldpath)
				url = olddir + url
			}

			var query string
			if i := strings.Index(url, "?"); i != -1 {
				url, query = url[:i], url[i:]
			}

			trailing := strings.HasSuffix(url, "/")
			url = path.Clean(url)
			if trailing && !strings.HasSuffix(url, "/") {
				url += "/"
			}
			url += query
		}
	}

	response.Header.Set("Location", url)
	if request.Method == "GET" || request.Method == "HEAD" {
		response.Header.Set("Content-Type", "text/html; charset=utf-8")
	}
	response.Status = code
	response.Body = ""

	if request.Method == "GET" {
		response.Body = "<a href=\"" + htmlEscape(url) + "\">" + http.StatusText(code) + "</a>.\n"
	}
}

// SetCookie sets cookies on response.
func (response *BaseResponse) SetCookie(cookie *http.Cookie) {
	if v := cookie.String(); v != "" {
		response.Header.Add("Set-Cookie", v)
	}
}

func (response *BaseResponse) setImplicitBody() {
	if response.Body == bodyNotSet {
		response.Body = http.StatusText(response.Status)
	}
}

func (response *BaseResponse) setImplicitStatus() {
	if response.Status != statusNotSet {
		return
	}

	if response.Body == bodyNotSet {
		response.Status = 404
		return
	}

	if _, ok := response.Body.(error); ok {
		response.Status = 500
		return
	}

	response.Status = 200
}

func (response *BaseResponse) setImplicit() {
	response.setImplicitStatus()
	response.setImplicitContentType()
	response.setImplicitBody()
}

func createResponse(rw http.ResponseWriter) *BaseResponse {
	return &BaseResponse{
		Writer: rw,
		Body:   bodyNotSet,
		Status: statusNotSet,
		Header: rw.Header(),
	}
}

// Response contains a typed body and shared response metadata. On normal return,
// Body is committed even when it has its zero value, unless the status is 400 or
// above, where a zero Body is replaced by the status text. Redirect, DontRespond
// and a body assigned to BaseResponse.Body take precedence over Body.
type Response[T any] struct {
	*BaseResponse
	Body T
}

// HasBody reports whether a response body has been set, including a zero value.
func (response *BaseResponse) HasBody() bool {
	return response.Body != bodyNotSet
}
