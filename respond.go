package vox

import (
	"encoding/json"
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
	// TODO: support io.Reader type
	default:
		body, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		res.Writer.Write(body)
	}

}
