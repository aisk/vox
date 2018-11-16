package vox

import (
	"mime"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	explictSetStatus = -1
	explictSetBody   = struct{}{}

	bodyMatcher = regexp.MustCompile("^\\s*<")
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

func hexEscapeNonASCII(s string) string {
	newLen := 0
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			newLen += 3
		} else {
			newLen++
		}
	}
	if newLen == len(s) {
		return s
	}
	b := make([]byte, 0, newLen)
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			b = append(b, '%')
			b = strconv.AppendInt(b, int64(s[i]), 16)
		} else {
			b = append(b, s[i])
		}
	}
	return string(b)
}

// A Response object contains all the information which will written to current
// HTTP client.
type Response struct {
	request *Request
	// Writer is the raw http.ResponseWriter for current request. You should
	// assign the Body / Status / Header value instead of using this field.
	Writer http.ResponseWriter
	// Body is the container for HTTP response's body.
	Body interface{}
	// The status code which will respond as the HTTP Response's status code.
	// 200 will be used as the default value if not set.
	Status int
	// Headers which will be written to the response.
	Header http.Header
}

func (response *Response) setImplictStatus() {
	if response.Status != explictSetStatus {
		// response's status is set by user.
		return
	}
	if response.Body == explictSetBody {
		// response's body is not set, set it to 404.
		response.Status = 404
		return
	}
	// response's status is not set by user, give it a default value.
	response.Status = 200
}

func (response *Response) setImplictContentType() {
	if response.Header.Get("Content-Type") != "" {
		return
	}

	if response.Body == explictSetBody {
		return
	}

	switch v := response.Body.(type) {
	case []byte:
		if bodyMatcher.Match(v) {
			response.Header.Set("Content-Type", mime.TypeByExtension(".html"))
		} else {
			response.Header.Set("Content-Type", mime.TypeByExtension(".text"))
		}
	case string:
		if bodyMatcher.MatchString(v) {
			response.Header.Set("Content-Type", mime.TypeByExtension(".html"))
		} else {
			response.Header.Set("Content-Type", mime.TypeByExtension(".text"))
		}
	// case map[string]interface{}, map[string]string:
	default:
		response.Header.Set("Content-Type", mime.TypeByExtension(".json"))
	}
}

var parseURL = url.Parse

// Redirect request to another url.
func (response *Response) Redirect(url string, code int) {
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

	if request.Method == "GET" {
		response.Body = "<a href=\"" + htmlEscape(url) + "\">" + http.StatusText(code) + "</a>.\n"
	}
}

func (response *Response) setImplictBody() {
	if response.Body == explictSetBody {
		response.Body = ""
	}
}

func (response *Response) setImplict() {
	response.setImplictStatus()
	response.setImplictContentType()
	response.setImplictBody()
}

func createResponse(rw http.ResponseWriter) *Response {
	return &Response{
		Writer: rw,
		Body:   explictSetBody,
		Status: explictSetStatus,
		Header: rw.Header(),
	}
}
