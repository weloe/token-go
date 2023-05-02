package go_http_context

import (
	"github.com/weloe/token-go/ctx"
	"reflect"
)

var _ ctx.Context = (*HttpContext)(nil)

type HttpContext struct {
	req        ctx.Request
	response   ctx.Response
	reqStorage ctx.ReqStorage
}

func (h *HttpContext) IsValidContext() bool {
	return h.req != nil && !reflect.DeepEqual(h.req, &HttpRequest{})
}

func (h *HttpContext) Request() ctx.Request {
	return h.req
}

func (h *HttpContext) ReqStorage() ctx.ReqStorage {
	return h.reqStorage
}

func (h *HttpContext) Response() ctx.Response {
	return h.response
}

func (h *HttpContext) MatchPath(pattern string, path string) bool {
	return true
}
