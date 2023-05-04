package go_http_context

import (
	"context"
	"github.com/weloe/token-go/ctx"
	"net/http"
)

type HttpReqStorage struct {
	source context.Context
}

func NewReqStorage(req *http.Request) *HttpReqStorage {
	if req == nil {
		return nil
	}
	return &HttpReqStorage{source: req.Context()}
}

func (r *HttpReqStorage) Source() interface{} {
	return r.source
}

func (r *HttpReqStorage) Get(key ctx.StorageKey) interface{} {
	return r.source.Value(key)
}

func (r *HttpReqStorage) Set(key ctx.StorageKey, value string) {
	r.source = context.WithValue(r.source, key, value)
}

func (r *HttpReqStorage) Delete(key ctx.StorageKey) {
	r.source = context.WithValue(r.source, key, nil)
}
