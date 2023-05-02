package config

type tokenConfig struct {
	// TokenName prefix
	TokenStyle  string
	TokenPrefix string
	TokenName   string

	Timeout         int64
	ActivityTimeout int64
	// Data clean period
	DataRefreshPeriod int64
	// Auto refresh token
	AutoRenew bool

	// Allow multi login
	IsConcurrent bool
	IsShare      bool
	// If (IsConcurrent == true && IsShare == false), support MaxLoginCount
	MaxLoginCount int16

	// Read token method
	// Set to true to read token from these method before login.
	IsReadBody   bool
	IsReadHeader bool
	IsReadCookie bool

	// Write token to response header.
	// Set to true to write after login.
	IsWriteHeader bool

	TokenSessionCheckLogin bool

	JwtSecretKey string

	CurDomain string

	SameTokenTimeout int64

	CheckSameToken bool

	CookieConfig *cookieConfig
}

func DefaultTokenConfig() *tokenConfig {
	return &tokenConfig{
		TokenStyle:             "uuid",
		TokenPrefix:            "",
		TokenName:              "tokenGo",
		Timeout:                60 * 60 * 24 * 30,
		ActivityTimeout:        -1,
		DataRefreshPeriod:      30,
		AutoRenew:              true,
		IsConcurrent:           true,
		IsShare:                true,
		MaxLoginCount:          12,
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
