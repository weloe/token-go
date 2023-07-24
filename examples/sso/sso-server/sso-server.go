package main

import (
	"encoding/json"
	"fmt"
	tokenGo "github.com/weloe/token-go"
	"github.com/weloe/token-go/config"
	"github.com/weloe/token-go/ctx"
	"github.com/weloe/token-go/model"
	"github.com/weloe/token-go/sso"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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
		Mode:          "",
		TicketTimeout: 300,
		AllowUrl:      "*",
		IsSlo:         true,

		IsHttp:    true,
		ServerUrl: "http://token-go-sso-server.com:9000",
		NotLoginView: func() interface{} {
			msg := "not log in SSO-Server, please visit <a href='/sso/doLogin?name=tokengo&pwd=123456' target='_blank'> doLogin </a>"
			return msg
		},
		DoLoginHandle: func(name string, pwd string, ctx ctx.Context) (interface{}, error) {
			if name != "tokengo" {
				return "name error", nil
			}
			if pwd != "123456" {
				return "pwd error", nil
			}
			token, err := enforcer.Login("1001", ctx)
			if err != nil {
				return nil, err
			}
			return model.Ok().SetData(token), nil
		},
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
	log.Fatal(http.ListenAndServe(":9000", engine))
}

// Engine is the uni handler for all requests
type Engine struct{}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if strings.HasPrefix(req.URL.String(), "/sso/") {
		res := ssoEnforcer.ServerDisPatcher(tokenGo.NewHttpContext(req, w))

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
