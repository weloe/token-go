package model

type Login struct {
	Device          string
	IsLastingCookie bool
	Timeout         int64
	JwtData         map[string]interface{}
	Token           string
	IsWriteHeader   bool
}

func DefaultLoginModel() *Login {
	return &Login{
		Device:          "default",
		IsLastingCookie: true,
		Timeout:         60 * 60 * 24 * 30,
		JwtData:         nil,
		Token:           "",
		IsWriteHeader:   true,
	}
}
