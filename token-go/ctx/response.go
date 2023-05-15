package ctx

type Response interface {
	Source() interface{}
	DeleteCookie(name string, path string, domain string)
	AddCookie(name string, value string, path string, domain string, timeout int64)
	SetHeader(name string, value string)
	AddHeader(name string, value string)
	SetServer(value string)
	Redirect(url string)
	Status(status int)
}
