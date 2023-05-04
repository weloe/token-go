package persist

import "github.com/weloe/token-go/model"

// Watcher event watcher
type Watcher interface {
	// Login called after login
	Login(loginType string, id interface{}, tokenValue string, loginModel *model.Login)
	// Logout called after logout
	Logout(loginType string, id interface{}, tokenValue string)
	// Kickout called when being kicked out of the server
	Kickout(loginType string, id interface{}, tokenValue string)
	// Replace called when Someone else has taken over your account
	Replace(loginType string, id interface{}, tokenValue string)
	// Ban called when account banned
	Ban(loginType string, id interface{}, service string)
	// UnBan called when account has been unbanned.
	UnBan(loginType string, id interface{}, service string)
	// RefreshToken called when renew token timeout
	RefreshToken(tokenValue string, id interface{}, timeout int64)
}
