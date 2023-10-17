package main

import (
	"fmt"
	tokenGo "github.com/weloe/token-go"
	"github.com/weloe/token-go/model"
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

	http.HandleFunc("/qrcode/create", create)
	http.HandleFunc("/qrcode/scanned", scanned)
	http.HandleFunc("/qrcode/confirmAuth", confirmAuth)
	http.HandleFunc("/qrcode/cancelAuth", cancelAuth)
	http.HandleFunc("/qrcode/getState", getState)

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func create(w http.ResponseWriter, request *http.Request) {
	// you should implement generate QR code method, returns QRCodeId to CreateQRCodeState
	// called generate QR code, returns QRCodeId to CreateQRCodeState
	//
	QRCodeId := "generatedQRCodeId"
	err := enforcer.CreateQRCodeState(QRCodeId, 50000)
	if err != nil {
		fmt.Fprintf(w, "CreateQRCodeState() failed: %v", err)
		return
	}
	fmt.Fprintf(w, "QRCodeId = %v", QRCodeId)
}

func scanned(w http.ResponseWriter, req *http.Request) {
	loginId, err := enforcer.GetLoginId(tokenGo.NewHttpContext(req, w))
	if err != nil {
		fmt.Fprintf(w, "GetLoginId() failed: %v", err)
		return
	}
	QRCodeId := req.URL.Query().Get("QRCodeId")
	tempToken, err := enforcer.Scanned(QRCodeId, loginId)
	if err != nil {
		fmt.Fprintf(w, "Scanned() failed: %v", err)
		return
	}
	fmt.Fprintf(w, "tempToken = %v", tempToken)
}
func getState(w http.ResponseWriter, req *http.Request) {
	QRCodeId := req.URL.Query().Get("QRCodeId")
	state := enforcer.GetQRCodeState(QRCodeId)
	if state == model.ConfirmAuth {
		qrCode := enforcer.GetQRCode(QRCodeId)
		if qrCode == nil {
			fmt.Fprintf(w, "login error. state = %v, code is nil", state)
			return
		}
		loginId := qrCode.LoginId
		token, err := enforcer.LoginById(loginId)
		if err != nil {
			fmt.Fprintf(w, "Login error: %s\n", err)
		}
		fmt.Fprintf(w, "%v login success. state = %v, token = %v", loginId, state, token)
		return
	} else if state == model.CancelAuth {
		fmt.Fprintf(w, "QRCodeId be cancelled: %v", QRCodeId)
		return
	}
	fmt.Fprintf(w, "state = %v", state)
}

func cancelAuth(w http.ResponseWriter, req *http.Request) {
	tempToken := req.URL.Query().Get("tempToken")
	err := enforcer.CancelAuth(tempToken)
	if err != nil {
		fmt.Fprintf(w, "CancelAuth() failed: %v", err)
		return
	}
	fmt.Fprint(w, "ConfirmAuth() success")
}

func confirmAuth(w http.ResponseWriter, req *http.Request) {
	tempToken := req.URL.Query().Get("tempToken")
	err := enforcer.ConfirmAuth(tempToken)
	if err != nil {
		fmt.Fprintf(w, "ConfirmAuth() failed: %v", err)
		return
	}
	fmt.Fprint(w, "ConfirmAuth() success")
}
