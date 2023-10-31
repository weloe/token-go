package model

import "fmt"

type RefreshRes struct {
	Token        string
	RefreshToken string
}

type Refresh struct {
	IsLastingCookie     bool
	Token               string
	Timeout             int64
	JwtData             map[string]interface{}
	RefreshToken        string
	RefreshTokenTimeout int64
}

func DefaultRefresh() *Refresh {
	return &Refresh{
		IsLastingCookie: true,
		Timeout:         60 * 60 * 24 * 30,
		JwtData:         nil,
		Token:           "",
	}
}

func (r *RefreshRes) String() string {
	return fmt.Sprintf("Token: %s, RefreshToken: %s", r.Token, r.RefreshToken)
}
