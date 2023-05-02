package go_http_context

import (
	"context"
	"net/http"
)

type HttpReqStorage struct {
	source context.Context
}

func NewReqStorage(req *http.Request) *HttpReqStorage {
	return &HttpReqStorage{source: req.Context()}
}

func (r HttpReqStorage) Source() interface{} {
	return r.source
}

func (r HttpReqStorage) Get(key string) interface{} {
	return r.source.Value(key)
}

func (r HttpReqStorage) Set(key string, value string) {
	r.source = context.WithValue(r.source, key, value)
}

func (r HttpReqStorage) Delete(key string) {
	r.source = context.WithValue(r.source, key, nil)
}
