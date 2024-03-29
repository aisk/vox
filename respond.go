package vox

import (
	"encoding/json"
	"io"
)

func respond(ctx *Context, req *Request, res *Response) {
	ctx.Next()
	if res.DontRespond {
		return
	}

	res.setImplicit()

	res.Writer.WriteHeader(res.Status)

	switch v := res.Body.(type) {
	case []byte:
		res.Writer.Write(v)
	case string:
		res.Writer.Write([]byte(v))
	case io.ReadCloser:
		_, err := io.Copy(res.Writer, v)
		if err != nil {
			panic(err)
		} else {
			if err = v.Close(); err != nil {
				panic(err)
			}
		}
	case io.Reader:
		_, err := io.Copy(res.Writer, v)
		if err != nil {
			panic(err)
		}
	case error:
		_, err := res.Writer.Write([]byte(v.Error()))
		if err != nil {
			panic(err)
		}
	default:
		body, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		res.Writer.Write(body)
	}
}
