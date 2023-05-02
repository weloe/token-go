package config

type cookieConfig struct {
	Domain   string
	Path     string
	Secure   bool
	HttpOnly bool
	SameSite string
}

func DefaultCookieConfig() *cookieConfig {
	return &cookieConfig{
		Domain:   "",
		Path:     "",
		Secure:   false,
		HttpOnly: false,
		SameSite: "",
	}
}
