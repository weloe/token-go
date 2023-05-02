package ctx

import (
	"fmt"
	"github.com/weloe/token-go/constant"
	"time"
)

var _ Response = (*DefaultRespImplement)(nil)

type DefaultRespImplement struct {
}

func (r *DefaultRespImplement) Source() interface{} {
	panic("implement me")
}

func (r *DefaultRespImplement) SetHeader(name string, value string) {
	panic("implement me")
}

func (r *DefaultRespImplement) AddHeader(name string, value string) {
	panic("implement me")
}

func (r *DefaultRespImplement) Redirect(url string) {
	panic("implement me")
}

func (r *DefaultRespImplement) Status(status int) {
	panic("implement me")
}

func (r *DefaultRespImplement) DeleteCookie(name string, path string, domain string) {
	r.AddCookie(name, "", path, domain, 0)
}

func (r *DefaultRespImplement) AddCookie(name string, value string, path string, domain string, timeout int64) {
	cookie := fmt.Sprintf("%s=%s; Path=%s; Domain=%s; Expires=%s",
		name,
		value,
		path,
		domain,
		time.Now().Add(time.Second*time.Duration(timeout)).Format(time.RFC1123),
	)
	r.AddHeader(constant.SetCookie, cookie)
}

func (r *DefaultRespImplement) SetServer(value string) {
	r.SetHeader("Server", value)
}
