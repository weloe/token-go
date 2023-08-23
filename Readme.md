# Token-Go

This library focuses on solving login authentication problems, such as: login, multi-account login, shared token, logout, kickout, banned, second auth, SSO ...

## Installation

```
go get github.com/weloe/token-go
```

## Simple Example

```go
import (
	"fmt"
	tokenGo "github.com/weloe/token-go"
	"log"
	"net/http"
)

var enforcer *tokenGo.Enforcer

func main() {
	var err error
	// use default adapter
	adapter := tokenGo.NewDefaultAdapter()
	enforcer, err = tokenGo.NewEnforcer(adapter)
	// enable logger
	enforcer.EnableLog()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/user/login", Login)
	http.HandleFunc("/user/logout", Logout)
	http.HandleFunc("/user/isLogin", IsLogin)
	http.HandleFunc("/user/kickout", Kickout)

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func Login(w http.ResponseWriter, req *http.Request) {
	token, err := enforcer.Login("1", tokenGo.NewHttpContext(req, w))
	if err != nil {
		fmt.Fprintf(w, "Login error: %s\n", err)
	}
	fmt.Fprintf(w, "token: %s\n", token)
}

func Logout(w http.ResponseWriter, req *http.Request) {
	err := enforcer.Logout(tokenGo.NewHttpContext(req, w))
	if err != nil {
		fmt.Fprintf(w, "Logout error: %s\n", err)
	}
	fmt.Fprintf(w, "logout success")
}

func IsLogin(w http.ResponseWriter, req *http.Request) {
	login, err := enforcer.IsLogin(tokenGo.NewHttpContext(req, w))
	if err != nil {
		fmt.Fprintf(w, "IsLogin() = %v: %v", login, err)
	}
	fmt.Fprintf(w, "IsLogin() = %v", login)
}

func Kickout(w http.ResponseWriter, req *http.Request) {
	err := enforcer.Kickout(req.URL.Query().Get("id"), "")
	if err != nil {
		fmt.Fprintf(w, "error: %s\n", err)
	}
	fmt.Fprintf(w, "logout success")
}
```

## Custom TokenConfig

The same user can only log in once:  `IsConcurrent = false && IsShare = false`

The same user logs in multiple times and shares a token:  `IsConcurrent = true && IsShare = false`

Multiple logins of the same user to multiple tokens:  `IsConcurrent = true && IsShare = true`

```go
import (
	"fmt"
	tokenGo "github.com/weloe/token-go"
	"github.com/weloe/token-go/config"
	"log"
	"net/http"
)

var enforcer *tokenGo.Enforcer

func main() {
	var err error
	// use default adapter
	adapter := tokenGo.NewDefaultAdapter()
	tokenConfig := &config.TokenConfig{
		TokenName:     "uuid",
		Timeout:       60,
		IsReadCookie:  true,
		IsReadHeader:  true,
		IsReadBody:    false,
		IsConcurrent:  true,
		IsShare:       true,
		MaxLoginCount: -1,
	}
	enforcer, err = tokenGo.NewEnforcer(adapter, tokenConfig)
}
```
You can also configure it using a yml or ini file like this

[token-go/token_conf.ini at master · weloe/token-go · GitHub](https://github.com/weloe/token-go/blob/master/examples/token_conf.ini)

[token-go/token_conf.yaml at master · weloe/token-go · GitHub](https://github.com/weloe/token-go/blob/master/examples/token_conf.yaml)

Then use `enforcer, err = tokenGo.NewEnforcer(adapter, filepath)`  to init.

## Authorization

A simple permission verification method is also provided
```go
type ACL interface {
	GetPermission(id string) []string
}
```
```go
type RBAC interface {
	GetRole(id string) []string
}
```
Implement either of these two interfaces and call `enforcer.SetAuth(model)`
After that, you can use these two APIs for permission verification

``` go
// implement RBAC
CheckRole(ctx ctx.Context, role string) error
// implement ACL
CheckPermission(ctx ctx.Context, permission string) error
```
### example

```go
type Auth struct {
}

func (m *Auth) GetRole(id string) []string {
	var arr = make([]string, 2)
	arr[1] = "user"
	return arr
}
func (m *Auth) GetPermission(id string) []string {
	var arr = make([]string, 2)
	arr[1] = "user::get"
	return arr
}


func main() {
	var err error
	// use default adapter
	adapter := tokenGo.NewDefaultAdapter()
	enforcer, err = tokenGo.NewEnforcer(adapter)
	// set auth
	enforcer.SetAuth(&Auth{})
	// enable logger
	enforcer.EnableLog()
	if err != nil {
		log.Fatal(err)
	}
	
	http.HandleFunc("/user/check", CheckAuth)
	
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func CheckAuth(w http.ResponseWriter, req *http.Request) {
	ctx := tokenGo.NewHttpContext(req, w)
	err := enforcer.CheckRole(ctx, "user")
	if err != nil {
		fmt.Fprintf(w, "CheckRole() error: %s\n", err)
		return
	}
	err = enforcer.CheckPermission(ctx, "user::get")
	if err != nil {
		fmt.Fprintf(w, "CheckPermission() error: %s\n", err)
		return
	}
	fmt.Fprintf(w, "you have authorization")
}
```
## SSO
SSO-Server examples: https://github.com/weloe/token-go/blob/master/examples/sso/sso-server/sso-server.go

SSO-Client examples: https://github.com/weloe/token-go/blob/master/examples/sso/sso-client-3/sso-client.go


## Api

[token_go package - github.com/weloe/token-go - Go Packages](https://pkg.go.dev/github.com/weloe/token-go#section-documentation)
