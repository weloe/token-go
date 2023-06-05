package ctx

type Request interface {
	Source() interface{}
	// Header get from request header
	Header(key string) string
	// PostForm get value from postForm
	PostForm(key string) string
	// Query https://example.org/?a=1&a=2&b=&=3&&&&" Query(a) return 1
	Query(key string) string
	// Path https://example.org/ex?a=1&a=2&b=&=3&&&& Path() return /ex
	Path() string
	// Url https://example.org/?a=1&a=2&b=&=3&&&& Url() return https://example.org/?a=1&a=2&b=&=3&&&&
	Url() string
	// UrlNoQuery return Url without query param
	UrlNoQuery() string
	// Method request method
	Method() string
	// Cookie get value from cookie
	Cookie(key string) string
}
