package go_http_context

import (
	"context"
	"github.com/weloe/token-go/ctx"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestHttpContext_IsValidContext(t *testing.T) {
	type fields struct {
		req        ctx.Request
		response   ctx.Response
		reqStorage ctx.ReqStorage
	}
	request, err := http.NewRequest("GET", "https://baidu.com", strings.NewReader(""))
	if err != nil {

	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "1",
			fields: fields{
				req:        nil,
				response:   nil,
				reqStorage: nil,
			},
			want: false,
		},
		{
			name: "2",
			fields: fields{
				req:        &HttpRequest{},
				response:   nil,
				reqStorage: nil,
			},
			want: false,
		},
		{
			name: "3",
			fields: fields{
				req:        &HttpRequest{source: &http.Request{}},
				response:   nil,
				reqStorage: nil,
			},
			want: true,
		},
		{
			name: "4",
			fields: fields{
				req:        &HttpRequest{source: request},
				response:   nil,
				reqStorage: nil,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HttpContext{
				req:        tt.fields.req,
				response:   tt.fields.response,
				reqStorage: tt.fields.reqStorage,
			}
			if got := h.IsValidContext(); got != tt.want {
				t.Errorf("IsValidContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpContext_MatchPath(t *testing.T) {
	type fields struct {
		req        ctx.Request
		response   ctx.Response
		reqStorage ctx.ReqStorage
	}
	type args struct {
		pattern string
		path    string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HttpContext{
				req:        tt.fields.req,
				response:   tt.fields.response,
				reqStorage: tt.fields.reqStorage,
			}
			if got := h.MatchPath(tt.args.pattern, tt.args.path); got != tt.want {
				t.Errorf("MatchPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpReqStorage_Delete(t *testing.T) {
	type fields struct {
		source context.Context
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := HttpReqStorage{
				source: tt.fields.source,
			}
			r.Delete(tt.args.key)
		})
	}
}

func TestHttpReqStorage_Get(t *testing.T) {
	type fields struct {
		source context.Context
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := HttpReqStorage{
				source: tt.fields.source,
			}
			if got := r.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpReqStorage_Set(t *testing.T) {
	type fields struct {
		source context.Context
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := HttpReqStorage{
				source: tt.fields.source,
			}
			r.Set(tt.args.key, tt.args.value)
		})
	}
}

func TestHttpRequest_Cookie(t *testing.T) {
	type fields struct {
		source *http.Request
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &HttpRequest{
				source: tt.fields.source,
			}
			if got := d.Cookie(tt.args.key); got != tt.want {
				t.Errorf("Cookie() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpRequest_Header(t *testing.T) {
	type fields struct {
		source *http.Request
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &HttpRequest{
				source: tt.fields.source,
			}
			if got := d.Header(tt.args.key); got != tt.want {
				t.Errorf("Header() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpRequest_Method(t *testing.T) {
	type fields struct {
		source *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &HttpRequest{
				source: tt.fields.source,
			}
			if got := d.Method(); got != tt.want {
				t.Errorf("Method() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpRequest_Path(t *testing.T) {
	type fields struct {
		source *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &HttpRequest{
				source: tt.fields.source,
			}
			if got := d.Path(); got != tt.want {
				t.Errorf("Path() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpRequest_PostForm(t *testing.T) {
	type fields struct {
		source *http.Request
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &HttpRequest{
				source: tt.fields.source,
			}
			if got := d.PostForm(tt.args.key); got != tt.want {
				t.Errorf("PostForm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpRequest_Query(t *testing.T) {
	type fields struct {
		source *http.Request
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &HttpRequest{
				source: tt.fields.source,
			}
			if got := d.Query(tt.args.key); got != tt.want {
				t.Errorf("Query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpRequest_Source(t *testing.T) {
	type fields struct {
		source *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &HttpRequest{
				source: tt.fields.source,
			}
			if got := d.Source(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Source() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpRequest_Url(t *testing.T) {
	type fields struct {
		source *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &HttpRequest{
				source: tt.fields.source,
			}
			if got := d.Url(); got != tt.want {
				t.Errorf("Url() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpResponse_AddHeader(t *testing.T) {
	type fields struct {
		DefaultRespImplement *ctx.DefaultRespImplement
		req                  *http.Request
		writer               http.ResponseWriter
	}
	type args struct {
		name  string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &HttpResponse{
				DefaultRespImplement: tt.fields.DefaultRespImplement,
				req:                  tt.fields.req,
				writer:               tt.fields.writer,
			}
			r.AddHeader(tt.args.name, tt.args.value)
		})
	}
}

func TestHttpResponse_HTML(t *testing.T) {
	type fields struct {
		DefaultRespImplement *ctx.DefaultRespImplement
		req                  *http.Request
		writer               http.ResponseWriter
	}
	type args struct {
		code int
		html string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &HttpResponse{
				DefaultRespImplement: tt.fields.DefaultRespImplement,
				req:                  tt.fields.req,
				writer:               tt.fields.writer,
			}
			if err := r.HTML(tt.args.code, tt.args.html); (err != nil) != tt.wantErr {
				t.Errorf("HTML() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHttpResponse_JSON(t *testing.T) {
	type fields struct {
		DefaultRespImplement *ctx.DefaultRespImplement
		req                  *http.Request
		writer               http.ResponseWriter
	}
	type args struct {
		code int
		obj  interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &HttpResponse{
				DefaultRespImplement: tt.fields.DefaultRespImplement,
				req:                  tt.fields.req,
				writer:               tt.fields.writer,
			}
			r.JSON(tt.args.code, tt.args.obj)
		})
	}
}

func TestHttpResponse_Redirect(t *testing.T) {
	type fields struct {
		DefaultRespImplement *ctx.DefaultRespImplement
		req                  *http.Request
		writer               http.ResponseWriter
	}
	type args struct {
		url string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &HttpResponse{
				DefaultRespImplement: tt.fields.DefaultRespImplement,
				req:                  tt.fields.req,
				writer:               tt.fields.writer,
			}
			r.Redirect(tt.args.url)
		})
	}
}

func TestHttpResponse_SetHeader(t *testing.T) {
	type fields struct {
		DefaultRespImplement *ctx.DefaultRespImplement
		req                  *http.Request
		writer               http.ResponseWriter
	}
	type args struct {
		name  string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &HttpResponse{
				DefaultRespImplement: tt.fields.DefaultRespImplement,
				req:                  tt.fields.req,
				writer:               tt.fields.writer,
			}
			r.SetHeader(tt.args.name, tt.args.value)
		})
	}
}

func TestHttpResponse_Source(t *testing.T) {
	type fields struct {
		DefaultRespImplement *ctx.DefaultRespImplement
		req                  *http.Request
		writer               http.ResponseWriter
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &HttpResponse{
				DefaultRespImplement: tt.fields.DefaultRespImplement,
				req:                  tt.fields.req,
				writer:               tt.fields.writer,
			}
			if got := r.Source(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Source() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpResponse_Status(t *testing.T) {
	type fields struct {
		DefaultRespImplement *ctx.DefaultRespImplement
		req                  *http.Request
		writer               http.ResponseWriter
	}
	type args struct {
		status int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &HttpResponse{
				DefaultRespImplement: tt.fields.DefaultRespImplement,
				req:                  tt.fields.req,
				writer:               tt.fields.writer,
			}
			r.Status(tt.args.status)
		})
	}
}

func TestNewHttpRequest(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want *HttpRequest
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHttpRequest(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHttpRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewReqStorage(t *testing.T) {
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name string
		args args
		want *HttpReqStorage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewReqStorage(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReqStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewResponse(t *testing.T) {
	type args struct {
		req    *http.Request
		writer http.ResponseWriter
	}
	tests := []struct {
		name string
		args args
		want *HttpResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResponse(tt.args.req, tt.args.writer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
