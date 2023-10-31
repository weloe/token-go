package model

type Login struct {
	Device              string
	IsLastingCookie     bool
	Timeout             int64
	JwtData             map[string]interface{}
	Token               string
	RefreshToken        string
	RefreshTokenTimeout int64
}

func DefaultLoginModel() *Login {
	return &Login{
		Device:          "default-device",
		IsLastingCookie: true,
		Timeout:         60 * 60 * 24 * 30,
		JwtData:         nil,
		Token:           "",
	}
}

func CreateLoginModelByDevice(device string) *Login {
	return &Login{
		Device:              device,
		IsLastingCookie:     true,
		Timeout:             60 * 60 * 24 * 30,
		JwtData:             nil,
		Token:               "",
		RefreshTokenTimeout: 60 * 60 * 24 * 30 * 2,
	}
}
