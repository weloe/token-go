package main

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
