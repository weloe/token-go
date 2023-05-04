package go_http_context

import (
	"encoding/json"
	"fmt"
	"github.com/weloe/token-go/ctx"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func NewTestRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest("GET", "https://baidu.com/api/", strings.NewReader(""))
	// cookie
	cookie := &http.Cookie{Name: "myCookie", Value: "cookieValue"}
	req.AddCookie(cookie)

	// header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "My Custom User Agent")

	// POST form
	form := url.Values{}
	form.Add("key1", "value1")
	form.Add("key2", "value2")
	req.PostForm = form

	// add query
	q := req.URL.Query()
	q.Add("query1", "value1")
	q.Add("query2", "value2")
	req.URL.RawQuery = q.Encode()
	if err != nil {
		t.Errorf("new request error: %v", err)
	}
	return req
}

func NewTestHttpRequest(t *testing.T) *HttpRequest {
	request := NewTestRequest(t)

	httpRequest := NewHttpRequest(request)
	return httpRequest
}

func NewTestHttpReqStore(t *testing.T) *HttpReqStorage {
	request := NewTestRequest(t)
	httpReqStorage := NewReqStorage(nil)
	if httpReqStorage != nil {
		t.Errorf("NewReqStorage() failed: value = %v", httpReqStorage)
	}
	httpReqStorage = NewReqStorage(request)
	return httpReqStorage
}

func TestHttpContext_IsValidContext(t *testing.T) {
	type fields struct {
		req        ctx.Request
		response   ctx.Response
		reqStorage ctx.ReqStorage
	}
	request := NewTestRequest(t)

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
	httpReqStorage := NewTestHttpReqStore(t)
	httpReqStorage.Set("s", "k")
	httpReqStorage.Delete("s")
	get := httpReqStorage.Get("s")

	if get != nil {
		t.Errorf("Delete error")
	}
}

func TestHttpReqStorage_Get(t *testing.T) {
	httpReqStorage := NewTestHttpReqStore(t)
	httpReqStorage.Set("s", "k")
	k := fmt.Sprintf("%v", httpReqStorage.Get("s"))
	if k != "k" {
		t.Errorf("get method error,Get() = %s want 'k'", k)
	}
}

func TestHttpRequest_Cookie(t *testing.T) {
	request := NewTestHttpRequest(t)
	cookie := request.Cookie("err")
	if cookie != "" {
		t.Errorf("Cookie() = %v,want ' '", cookie)
	}
	cookie = request.Cookie("myCookie")
	if cookie != "cookieValue" {
		t.Errorf("Cookie() = %v,want cookieValue", cookie)
	}
}

func TestHttpRequest_Header(t *testing.T) {
	request := NewTestHttpRequest(t)
	hV := request.Header("err")
	if hV != "" {
		t.Errorf("Header() = %v,want ' '", hV)
	}
	hV = request.Header("Content-Type")
	if hV != "application/x-www-form-urlencoded" {
		t.Errorf("Header() = %v,want application/x-www-form-urlencoded", hV)
	}
}

func TestHttpRequest_Method(t *testing.T) {
	request := NewTestHttpRequest(t)
	if request.Method() != "GET" {
		t.Errorf("Method() = %s,want GET", request.Method())
	}
}

func TestHttpRequest_Path(t *testing.T) {
	request := NewTestHttpRequest(t)
	t.Log(request.Path())
	if request.Path() != "/api/" {
		t.Errorf("Path() = %s  want /api/", request.Path())
	}
}

func TestHttpRequest_PostForm(t *testing.T) {
	request := NewTestHttpRequest(t)
	if request.PostForm("key") != "" {
		t.Errorf("PostForm() = %s want ' '", request.PostForm("key"))
	}

	if request.PostForm("key1") != "value1" {
		t.Errorf("PostForm() = %s want ' '", request.PostForm("key1"))
	}

	if request.PostForm("key2") != "value2" {
		t.Errorf("PostForm() = %s want ' '", request.PostForm("key2"))
	}
}

func TestHttpRequest_Query(t *testing.T) {
	request := NewTestHttpRequest(t)
	if request.Query("key") != "" {
		t.Errorf("Query() = %s want ' '", request.Query("key"))
	}

	if request.Query("query1") != "value1" {
		t.Errorf("Query() = %s want ' '", request.Query("query1"))
	}

	if request.Query("query2") != "value2" {
		t.Errorf("Query() = %s want ' '", request.Query("query2"))
	}
}

func TestHttpRequest_Source(t *testing.T) {
	request := NewTestHttpRequest(t)
	t.Log(request.Source())
}

func TestHttpRequest_Url(t *testing.T) {
	s := NewTestHttpRequest(t).Url()
	if s != "https://baidu.com/api/?query1=value1&query2=value2" {
		t.Errorf("Url() = %s ,want https://baidu.com/api/?query1=value1&query2=value2", s)
	}
}

func TestHttpResponse_SetHeader(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	resp := NewResponse(req, w)
	resp.SetHeader("Content-Type", "application/json")

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("SetHeader() failed: expected Content-Type header to be 'application/json', got '%s'", ct)
	}
}

func TestHttpResponse_AddHeader(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	resp := NewResponse(req, w)
	resp.AddHeader("Content-Type", "application/json")
	resp.AddHeader("Cache-Control", "no-cache")

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("AddHeader() failed: expected Content-Type header to be 'application/json', got '%s'", ct)
	}
	if cc := w.Header().Get("Cache-Control"); cc != "no-cache" {
		t.Errorf("AddHeader() failed: expected Cache-Control header to be 'no-cache', got '%s'", cc)
	}
}

func TestHttpResponse_Redirect(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	resp := NewResponse(req, w)
	resp.Redirect("/new-page")

	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("Redirect() failed: expected status code %d, got %d", http.StatusTemporaryRedirect, w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/new-page" {
		t.Errorf("Redirect() failed: expected Location header to be '/new-page', got '%s'", loc)
	}
}

func TestHttpResponse_Status(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	resp := NewResponse(req, w)
	resp.Status(http.StatusOK)

	if w.Code != http.StatusOK {
		t.Errorf("Status() failed: expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHttpResponse_JSON(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	resp := NewResponse(req, w)

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	person := &Person{Name: "John Doe", Age: 30}
	resp.JSON(http.StatusOK, person)

	if w.Code != http.StatusOK {
		t.Errorf("JSON() failed: expected status code %d, got %d", http.StatusOK, w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("JSON() failed: expected Content-Type header to be 'application/json', got '%s'", ct)
	}

	var p Person
	decoder := json.NewDecoder(w.Body)
	err := decoder.Decode(&p)
	if err != nil {
		t.Fatalf("Error decoding response body: %s", err)
	}
	if p.Name != "John Doe" {
		t.Errorf("JSON() failed: expected person name to be 'John Doe', got '%s'", p.Name)
	}
	if p.Age != 30 {
		t.Errorf("JSON() failed: expected person age to be '30', got '%d'", p.Age)
	}
}

func TestHttpResponse_HTML(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	resp := NewResponse(req, w)

	htmlCode := "<html><body><h1>Hello World!</h1></body></html>"
	err := resp.HTML(http.StatusOK, htmlCode)
	if err != nil {
		t.Fatalf("Error writing HTML code to response: %s", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("HTML() failed: expected status code %d, got %d", http.StatusOK, w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/html" {
		t.Errorf("HTML() failed: expected Content-Type header to be 'text/html', got '%s'", ct)
	}
	if hc := w.Body.String(); hc != htmlCode {
		t.Errorf("HTML() failed: expected HTML code to be '%s'", hc)
	}
}

func TestSetCookieHandler(t *testing.T) {
	// Create a new request that simulates a call to the /setcookie endpoint
	req, err := http.NewRequest("GET", "/setcookie", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()

	// Call SetCookieHandler to set a new cookie
	// Create a new cookie
	cookie := &http.Cookie{
		Name:    "my_cookie",
		Value:   "1234567890",
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
		Domain:  "localhost",
	}
	NewResponse(req, rr).AddCookie(cookie.Name, cookie.Value, cookie.Path, cookie.Domain, 2)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check that the cookie was set correctly
	cookies := rr.Header().Values("Set-Cookie")
	if len(cookies) != 1 {
		t.Errorf("handler returned wrong number of Set-Cookie headers: got %v want %v",
			len(cookies), 1)
	}
	if got, want := cookies[0], "my_cookie=1234567890; Path=/; Domain=localhost; Expires="; !containsString(got, want) {
		t.Errorf("handler returned unexpected Set-Cookie header: got %v want %v",
			got, want)
	}
}

func TestDeleteCookieHandler(t *testing.T) {
	// Create a new request that simulates a call to the /deletecookie endpoint
	req, err := http.NewRequest("GET", "/deletecookie", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()

	// Call DeleteCookieHandler to delete an existing cookie
	// Create a new cookie with the same name and domain as the cookie to be deleted
	cookie := &http.Cookie{
		Name:   "my_cookie",
		Value:  "",
		Path:   "/",
		Domain: "localhost",
		MaxAge: -1,
	}

	// Set the cookie in the response header with MaxAge = -1 to delete it
	NewResponse(req, rr).DeleteCookie(cookie.Name, cookie.Path, cookie.Domain)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check that the cookie was deleted correctly
	cookies := rr.Header().Values("Set-Cookie")
	if len(cookies) != 1 {
		t.Errorf("handler returned wrong number of Set-Cookie headers: got %v want %v",
			len(cookies), 1)
	}
	if got, want := cookies[0], "my_cookie=; Path=/; Domain=localhost; Max-Age=0"; !containsString(got, want) {
		t.Errorf("handler returned unexpected Set-Cookie header: got %v want %v",
			got, want)
	}
}

func containsString(s string, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}

func TestNewHttpContext(t *testing.T) {
	context := NewHttpContext(nil, nil)
	if context == nil {
		t.Errorf("NewHttpContext() failed: value = %v", context)
	}
	request := context.Request()
	response := context.Response()
	storage := context.ReqStorage()
	if request == nil || response == nil || storage == nil {
		t.Errorf("HttpContext failed ")
	}
}
