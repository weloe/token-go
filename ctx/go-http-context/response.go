package go_http_context

import (
	"encoding/json"
	"github.com/weloe/token-go/constant"
	"github.com/weloe/token-go/ctx"
	"net/http"
	"time"
)

var _ ctx.Response = (*HttpResponse)(nil)

type HttpResponse struct {
	req    *http.Request
	writer http.ResponseWriter
}

func NewResponse(req *http.Request, writer http.ResponseWriter) *HttpResponse {
	return &HttpResponse{
		req:    req,
		writer: writer,
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

func (r *HttpResponse) DeleteCookie(name string, path string, domain string) {
	cookie := http.Cookie{
		Name:   name,
		Value:  "",
		Path:   path,
		Domain: domain,
		MaxAge: -1,
	}
	r.AddHeader(constant.SetCookie, cookie.String())
}

func (r *HttpResponse) AddCookie(name string, value string, path string, domain string, timeout int64) {
	var expiration time.Time
	if timeout == -1 {
		expiration = time.Unix(0, 0)
	} else {
		expiration = time.Now().Add(time.Second * time.Duration(timeout))
	}
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Path:    path,
		Domain:  domain,
		Expires: expiration,
	}
	r.AddHeader(constant.SetCookie, cookie.String())
}

func (r *HttpResponse) SetServer(value string) {
	r.SetHeader("Server", value)
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
