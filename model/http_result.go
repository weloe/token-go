package model

const (
	SUCCESS = 1
	ERROR   = 0
)

// Result wrap the http request result.
type Result struct {
	Code int
	Msg  string
	Data interface{}
}

func Ok() *Result {
	return &Result{
		Code: SUCCESS,
		Msg:  "success",
		Data: nil,
	}
}

func Error() *Result {
	return &Result{
		Code: -1,
		Msg:  "error",
		Data: nil,
	}
}

func (r *Result) SetData(data interface{}) *Result {
	r.Data = data
	return r
}

func (r *Result) SetMsg(msg string) *Result {
	r.Msg = msg
	return r
}
