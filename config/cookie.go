package config

type CookieConfig struct {
	Domain   string
	Path     string
	Secure   bool
	HttpOnly bool
	SameSite string
}

func DefaultCookieConfig() *CookieConfig {
	return &CookieConfig{
		Domain:   "",
		Path:     "",
		Secure:   false,
		HttpOnly: false,
		SameSite: "",
	}
}
