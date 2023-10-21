package config

import "github.com/weloe/token-go/constant"

type TokenConfig struct {
	// TokenStyle
	// uuid | uuid-simple | random-string32 | random-string64 | random-string128
	TokenStyle string
	// TokenName prefix
	TokenPrefix string
	TokenName   string

	Timeout int64
	// If last operate time < ActivityTimeout, token expired
	ActivityTimeout int64
	// Data clean period
	DataRefreshPeriod int64
	// Auto refresh token
	AutoRenew bool

	// Allow multi login
	IsConcurrent bool
	// Multi login share same token
	IsShare bool
	// If (IsConcurrent == true && IsShare == false), support MaxLoginCount
	// If IsConcurrent == -1, do not need to check loginCount
	MaxLoginCount int16
	// Maximum number of logins per device
	DeviceMaxLoginCount int16

	// Read token method
	// Set to true to read token from these method before login.
	IsReadBody   bool
	IsReadHeader bool
	// If IsReadCookie is set to true, a cookie will be set after successful login
	IsReadCookie bool

	// Write token to response header.
	// Set to true to write after login.
	IsWriteHeader bool

	TokenSessionCheckLogin bool

	JwtSecretKey string

	CurDomain string

	SameTokenTimeout int64

	CheckSameToken bool

	CookieConfig *CookieConfig
}

func DefaultTokenConfig() *TokenConfig {
	return &TokenConfig{
		TokenStyle:             "uuid",
		TokenPrefix:            "",
		TokenName:              constant.TokenName,
		Timeout:                60 * 60 * 24 * 30,
		ActivityTimeout:        -1,
		DataRefreshPeriod:      30,
		AutoRenew:              true,
		IsConcurrent:           true,
		IsShare:                true,
		MaxLoginCount:          12,
		DeviceMaxLoginCount:    12,
		IsReadBody:             true,
		IsReadHeader:           true,
		IsReadCookie:           true,
		IsWriteHeader:          false,
		TokenSessionCheckLogin: true,
		JwtSecretKey:           "",
		CurDomain:              "",
		SameTokenTimeout:       60 * 60 * 24,
		CheckSameToken:         false,
		CookieConfig:           DefaultCookieConfig(),
	}
}
