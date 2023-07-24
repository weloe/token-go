package main

import (
	"encoding/json"
	"fmt"
	tokenGo "github.com/weloe/token-go"
	"github.com/weloe/token-go/config"
	"github.com/weloe/token-go/model"
	"github.com/weloe/token-go/sso"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var enforcer *tokenGo.Enforcer

var ssoEnforcer *sso.SsoEnforcer

func main() {
	var err error
	// use default adapter
	adapter := tokenGo.NewDefaultAdapter()
	enforcer, err = tokenGo.NewEnforcer(adapter)
	if err != nil {
		log.Fatal(err)
	}
	// enable logger
	enforcer.EnableLog()
	ssoOptions := &config.SsoOptions{
		AuthUrl:        "/sso/auth",
		IsSlo:          true,
		IsHttp:         true,
		SloUrl:         "/sso/signout",
		CheckTicketUrl: "/sso/checkTicket",
		ServerUrl:      "http://token-go-sso-server.com:9000",
		SendHttp: func(url string) (string, error) {
			response, err := http.Get(url)
			if err != nil {
				log.Printf("http.Get() failed: %v", err)
				return "", err
			}

			defer func(Body io.ReadCloser) {
				err = Body.Close()
				if err != nil {
					log.Printf("read response body failed: %v", err)
				}
			}(response.Body)

			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Printf("ioutil.ReadAll() failed: %v", err)
				return "", err
			}

			return string(body), nil
		},
	}
	signOptions := &config.SignOptions{
		SecretKey:    "kQwIOrYvnXmSDkwEiFngrKidMcdrgKor",
		IsCheckNonce: true,
	}
	ssoEnforcer, err = sso.NewSsoEnforcer(&sso.Options{
		SsoOptions:  ssoOptions,
		SignOptions: signOptions,
		Enforcer:    enforcer,
	})
	if err != nil {
		log.Fatalf("NewSsoEnforcer() failed: %v", err)
	}

	engine := new(Engine)
	log.Fatal(http.ListenAndServe(":9001", engine))
}

// Engine is the uni handler for all requests
type Engine struct{}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.String() == "/" {
		isLogin, err := enforcer.IsLogin(tokenGo.NewHttpContext(req, w))
		if err != nil {
			fmt.Fprintf(w, "enforcer.IsLogin() failed: %v", err)
			return
		}
		response := "<h2>token-go SSO-Client</h2> <p>isLogin = " + strconv.FormatBool(isLogin) + "</p> <p><a href=\"javascript:location.href='/sso/login?back='+ encodeURIComponent(location.href);\">login</a>  <a href='/sso/logout?back=self'>logout</a></p>"
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, response)
	} else if strings.HasPrefix(req.URL.String(), "/sso/") {
		res := ssoEnforcer.ClientDispatcher(tokenGo.NewHttpContext(req, w))

		result, ok := res.(*model.Result)
		if ok {
			bytes, err := json.Marshal(result)
			if err != nil {
				fmt.Fprintf(w, "json.Marshal() = %v", err)
				return
			}

			_, err = w.Write(bytes)
			if err != nil {
				fmt.Fprintf(w, "w.Write() = %v", err)
				return
			}
			return
		}
		html, ok := res.(string)
		if ok {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, "%s", html)
			return
		}

		bytes, err := json.Marshal(model.Ok())
		if err != nil {
			fmt.Fprintf(w, "json.Marshal() = %v", err)
			return
		}

		_, err = w.Write(bytes)
		if err != nil {
			fmt.Fprintf(w, "w.Write() = %v", err)
			return
		}

	} else {
		fmt.Fprintf(w, "not this api")
	}
}
