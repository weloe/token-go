package ctx

type Context interface {
	Request() Request
	Response() Response
	ReqStorage() ReqStorage
	MatchPath(pattern string, path string) bool
	IsValidContext() bool
}
