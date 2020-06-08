package request

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"

	"github.com/valyala/fasthttp"
)

func NewRequest(url string, fs map[string]*Part) (*fasthttp.Request, error) {
	var body bytes.Buffer
	var req fasthttp.Request

	w := multipart.NewWriter(&body)
	for k, v := range fs {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, k, k))
		h.Set("Content-Type", v.Typ)
		p, _ := w.CreatePart(h)
		io.Copy(p, bytes.NewReader(v.Data))
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	req.SetRequestURI(url)
	req.SetBody(body.Bytes())
	req.Header.SetMethod("POST")
	req.Header.Add("Content-Type", w.FormDataContentType())
	return &req, nil
}
