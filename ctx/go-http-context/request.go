package go_http_context

import (
	"github.com/weloe/token-go/ctx"
	"net/http"
)

var _ ctx.Request = (*HttpRequest)(nil)

type HttpRequest struct {
	source *http.Request
}

func NewHttpRequest(r *http.Request) *HttpRequest {
	return &HttpRequest{r}
}

func (d *HttpRequest) Source() interface{} {
	return d.source
}

func (d *HttpRequest) Header(key string) string {
	return d.source.Header.Get(key)
}

func (d *HttpRequest) PostForm(key string) string {
	return d.source.PostFormValue(key)
}

func (d *HttpRequest) Query(key string) string {
	return d.source.URL.Query().Get(key)
}

func (d *HttpRequest) Path() string {
	return d.source.URL.Path
}

func (d *HttpRequest) Url() string {
	return d.source.URL.String()
}

func (d *HttpRequest) Method() string {
	return d.source.Method
}

func (d *HttpRequest) Cookie(key string) string {
	cookie, err := d.source.Cookie(key)
	if err != nil {
		return ""
	}
	return cookie.Value
}
