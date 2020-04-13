package vox

import (
	"encoding/json"
	"io"
)

func respond(req *Request, res *Response) {
	req.Next()
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
	default:
		body, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		res.Writer.Write(body)
	}
}
