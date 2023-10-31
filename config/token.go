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

	// If you enable DoubleToken, returns double token - token and refreshToken, when token is timeout, you can use refreshToken to refreshToken to auto login.
	DoubleToken         bool
	RefreshTokenName    string
	RefreshTokenTimeout int64

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

func (t *TokenConfig) InitConfig() {
	if t.TokenStyle == "" {
		t.TokenStyle = "uuid"
	}
	if t.TokenName == "" {
		t.TokenName = constant.TokenName
	}
	if t.Timeout == 0 {
		t.Timeout = 60 * 60 * 24 * 30
	}
	if t.DeviceMaxLoginCount == 0 {
		t.DeviceMaxLoginCount = 12
	}
	if t.DoubleToken {
		if t.RefreshTokenName == "" {
			t.RefreshTokenName = constant.RefreshToken
		}
		if t.RefreshTokenTimeout == 0 {
			t.RefreshTokenTimeout = t.Timeout * 2
		}
	}
	if t.MaxLoginCount == 0 {
		t.MaxLoginCount = 12
	}
	if t.CookieConfig == nil {
		t.CookieConfig = DefaultCookieConfig()
	}

}

func DefaultTokenConfig() *TokenConfig {
	return &TokenConfig{
		TokenStyle:             "uuid",
		TokenPrefix:            "",
		TokenName:              constant.TokenName,
		Timeout:                60 * 60 * 24 * 30,
		DoubleToken:            false,
		ActivityTimeout:        -1,
		DataRefreshPeriod:      30,
		AutoRenew:              false,
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
