package go_http_context

import (
	"encoding/json"
	"github.com/weloe/token-go/ctx"
	"net/http"
)

var _ ctx.Response = (*HttpResponse)(nil)

type HttpResponse struct {
	*ctx.DefaultRespImplement
	req    *http.Request
	writer http.ResponseWriter
}

func NewResponse(req *http.Request, writer http.ResponseWriter) *HttpResponse {
	return &HttpResponse{
		DefaultRespImplement: &ctx.DefaultRespImplement{},
		req:                  req,
		writer:               writer,
	}
}

func (r *HttpResponse) Source() interface{} {
	return r.writer
}

func (r *HttpResponse) SetHeader(name string, value string) {
	r.writer.Header().Set(name, value)
}

func (r *HttpResponse) AddHeader(name string, value string) {
	r.writer.Header().Add(name, value)
}

func (r *HttpResponse) Redirect(url string) {
	http.Redirect(r.writer, r.req, url, http.StatusTemporaryRedirect)
}

func (r *HttpResponse) Status(status int) {
	r.writer.WriteHeader(status)
}

// JSON response json data
func (r *HttpResponse) JSON(code int, obj interface{}) {
	r.SetHeader("Content-Type", "application/json")
	r.Status(code)

	encoder := json.NewEncoder(r.writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(r.writer, err.Error(), 500)
	}
}

// HTML response .html
func (r *HttpResponse) HTML(code int, html string) error {
	r.SetHeader("Content-Type", "text/html")
	r.Status(code)
	_, err := r.writer.Write([]byte(html))
	if err != nil {
		return err
	}
	return nil
}
