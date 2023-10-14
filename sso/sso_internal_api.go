package sso

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/weloe/token-go/constant"
	"github.com/weloe/token-go/ctx"
	"github.com/weloe/token-go/model"
	"github.com/weloe/token-go/util"
	"log"
	"strconv"
	"strings"
	"time"
)

/**
=========internal api
*/

// checkRequest check request param timestamp,nonce,sign.
func (s *SsoEnforcer) checkRequest(request ctx.Request) error {
	timestamp := request.Query(s.paramName.TimeStamp)
	nonce := request.Query(s.paramName.Nonce)
	sign := request.Query(s.paramName.Sign)
	err := s.checkTimeStamp(timestamp)
	if err != nil {
		return err
	}
	if s.signConfig.IsCheckNonce {
		err = s.checkNonce(nonce)
		if err != nil {
			return err
		}
	}
	err = s.checkSign(timestamp, nonce, sign)

	if err != nil {
		return err
	}
	return nil
}

// CreateTicket create ticket by account-id.
func (s *SsoEnforcer) CreateTicket(loginId string, client string) (string, error) {
	// create random string ticket
	ticket, err := util.GenerateRandomString64()
	if err != nil {
		return "", err
	}
	// save ticket-id+client
	err = s.saveTicket(ticket, loginId, client)
	if err != nil {
		return "", err
	}
	// save id-ticket
	err = s.saveTicketIndex(ticket, loginId)
	if err != nil {
		return "", err
	}

	return ticket, nil
}

// GetLoginId get loginId by ticket.
func (s *SsoEnforcer) GetLoginId(ticket string) string {
	if ticket == "" {
		return ""
	}
	loginId := s.getLoginIdByTicket(ticket)
	if loginId != "" && strings.Contains(loginId, ",") {
		split := strings.Split(loginId, ",")
		loginId = split[0]
	}
	return loginId
}

func (s *SsoEnforcer) getLoginIdByTicket(ticket string) string {
	loginId := s.enforcer.GetAdapter().GetStr(s.spliceTicketSaveKey(ticket))
	return loginId
}

// GetTicket get ticket by loginId.
func (s *SsoEnforcer) GetTicket(loginId string) string {
	if loginId == "" {
		return ""
	}
	return s.enforcer.GetAdapter().GetStr(s.spliceTicketIndexKey(loginId))
}

// CheckTicket use config.Client to check ticket,return loginId.
func (s *SsoEnforcer) CheckTicket(ticket string) (string, error) {
	return s.CheckTicketByClient(ticket, s.config.Client)
}

// CheckTicketByClient check ticket by pointing client,return loginId.
func (s *SsoEnforcer) CheckTicketByClient(ticket string, client string) (string, error) {
	id := s.getLoginIdByTicket(ticket)
	if id == "" {
		return "", nil
	}

	// get client from id
	var ticketClient string
	if strings.Contains(id, ",") {
		split := strings.Split(id, ",")
		id = split[0]
		ticketClient = split[1]
	}

	if client != "" && client != ticketClient {
		return "", fmt.Errorf("the ticket does not belong to the client, client: %v, ticket: %v", client, ticket)
	}
	err := s.deleteTicket(ticket)
	if err != nil {
		return "", err
	}
	err = s.deleteTicketIndex(id)
	if err != nil {
		return "", err
	}

	return id, nil
}

// CheckRedirectUrl check redirectUrl.
func (s *SsoEnforcer) CheckRedirectUrl(url string) error {
	if !util.IsValidUrl(url) {
		return fmt.Errorf("invalid redirect url: %v", url)
	}
	index := strings.Index(url, "?")
	if index != -1 {
		url = url[0:index]
	}
	allowUrls := strings.Split(s.GetAllowUrl(), ",")

	if !util.HasUrl(allowUrls, url) {
		return fmt.Errorf("illegal redirect url: %v", url)
	}
	return nil
}

// RegisterSloCallbackUrl register the URL of the single logout callback for the account id.
func (s *SsoEnforcer) RegisterSloCallbackUrl(loginId string, sloCallbackUrl string) error {
	if loginId == "" || sloCallbackUrl == "" {
		return nil
	}
	session := s.enforcer.GetSession(loginId)
	// splice session id
	sessionId := s.enforcer.GetTokenConfig().TokenName + ":" + s.enforcer.GetType() + ":session:" + loginId
	if session == nil {
		session = model.NewSession(sessionId, "account-session", loginId)
	}
	value := session.Get(constant.SLO_CALLBACK_SET_KEY)

	var v []string
	if value != nil {
		sv, ok := value.([]string)
		if !ok {
			return errors.New("session SLO_CALLBACK_SET_KEY_ data convert into []string failed")
		}
		v = sv
	}
	v = util.AppendStr(v, sloCallbackUrl)

	session.Set(sessionId, v)
	// update session
	err := s.enforcer.UpdateSession(loginId, session)
	if err != nil {
		return err
	}

	return nil
}

// ssoSignOutById single sign-out of the specified account.
// Use loginId to get single sign-out urls from session.
func (s *SsoEnforcer) ssoSignOutById(loginId string) error {
	// if loginId is not logged, return error
	session := s.enforcer.GetSession(loginId)
	if session == nil {
		return errors.New("this loginId is not logged in")
	}
	value := session.Get(constant.SLO_CALLBACK_SET_KEY)
	if value == nil {
		return nil
	}
	urls, ok := value.([]string)
	if !ok {
		return errors.New("convert into []string failed")
	}
	// range urls to make client logout
	for _, url := range urls {
		// join url
		newUrl, err := s.joinLoginIdAndSign(url, loginId)
		if err != nil {
			return err
		}
		// sent http
		_, err = s.config.SendHttp(newUrl)
		if err != nil {
			return err
		}
	}

	// server logout
	return s.enforcer.LogoutById(loginId)
}

// joinLoginIdAndSign splice the loginId to the url, and stitching parameters such as sign.
func (s *SsoEnforcer) joinLoginIdAndSign(url string, id string) (string, error) {
	nonce, err := util.GenerateRandomString32()
	if err != nil {
		return "", err
	}

	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	sign, err := s.createSign(timestamp, nonce)
	if err != nil {
		return "", err
	}
	str := url + "?" + s.paramName.LoginId + "=" + id + "&" + s.paramName.TimeStamp + "=" + timestamp + "&" + s.paramName.Nonce + "=" + nonce + "&" + s.paramName.Sign + "=" + sign
	return str, nil
}

// buildServerAuthUrl SSO-Client build SSO-Server single sign-on url.
func (s *SsoEnforcer) buildServerAuthUrl(clientLoginUrl string, back string) (string, error) {
	if clientLoginUrl == "" {
		return "", errors.New("arg[0] clientLoginUrl can not be nil")
	}

	// get server auth url
	authUrl := s.config.SpliceAuthUrl()

	client := s.config.Client

	if client != "" {
		authUrl = util.AddQuery(authUrl, s.paramName.Client, client)
	}

	// splice back url
	if back != "" {
		back = util.Encode(back)
		clientLoginUrl = util.AddQuery(clientLoginUrl, s.paramName.Back, back)
	}

	return util.AddQuery(authUrl, s.paramName.Redirect, clientLoginUrl), nil
}

// buildRedirectUrl the server gives the redirectUrl of the ticket to the client.
// Check redirect url, delete old ticket of loginId and create new ticket, then return url with new ticket.
func (s *SsoEnforcer) buildRedirectUrl(loginId string, client string, redirect string) (string, error) {
	// check redirect url
	err := s.CheckRedirectUrl(redirect)
	if err != nil {
		return "", err
	}
	// delete old ticket
	err = s.deleteTicket(s.GetTicket(loginId))
	if err != nil {
		return "", err
	}
	// create new ticket
	ticket, err := s.CreateTicket(loginId, client)
	if err != nil {
		return "", err
	}

	// return redirect + "?" + s.paramName.Ticket + "=" + ticket, nil
	return util.AddQuery(s.encodeBackParam(redirect), s.paramName.Ticket, ticket), nil
}

// encodeBackParam find back param from url, and encode back param.
func (s *SsoEnforcer) encodeBackParam(url string) string {
	// get back location
	index := strings.Index(url, "?"+s.paramName.Back+"=")
	if index == -1 {
		index = strings.Index(url, "&"+s.paramName.Back+"=")
		if index == -1 {
			return url
		}
	}

	// encode
	length := len(s.paramName.Back) + 2
	back := url[index+length:]
	back = util.Encode(back)

	// update back
	url = url[:index+length] + back
	return url
}

// buildCheckTicketUrl build to check ticket.
func (s *SsoEnforcer) buildCheckTicketUrl(ticket string, ssoLogoutCallUrl string) (string, error) {
	if ticket == "" {
		return "", errors.New("buildCheckTicketUrl() ticket can not be nil")
	}
	checkTicketUrl := s.config.SpliceCheckTicketUrl()
	client := s.config.Client
	paramMap := make(map[string]string)
	if client != "" {
		paramMap[s.paramName.Client] = client
	}
	paramMap[s.paramName.Ticket] = ticket
	if ssoLogoutCallUrl != "" {
		paramMap[s.paramName.SsoLogoutCall] = ssoLogoutCallUrl
	}

	return util.AddQueryMap(checkTicketUrl, paramMap), nil
}

// buildSloUrl build single-logout url.
func (s *SsoEnforcer) buildSloUrl(loginId string) (string, error) {
	sloUrl := s.config.SpliceSloUrl()
	url, err := s.joinLoginIdAndSign(sloUrl, loginId)
	if err != nil {
		return "", err
	}
	return url, nil
}

// buildGetDataUrl build getData url with sign,timestamp,nonce.
func (s *SsoEnforcer) buildGetDataUrl(paramMap map[string]string) (string, error) {
	getDataUrl := s.config.SpliceGetDataUrl()
	return s.buildCustomPathUrl(getDataUrl, paramMap)
}

// buildCustomPathUrl add paramMap to path.
func (s *SsoEnforcer) buildCustomPathUrl(path string, paramMap map[string]string) (string, error) {
	u := path

	if !strings.HasPrefix(u, "http") {
		serverUrl := s.config.ServerUrl
		if serverUrl == "" {
			return "", errors.New("please set sso serverUrl")
		}
		u = util.SpliceUrl(serverUrl, path)
	}
	// sign map
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	paramMap[s.paramName.TimeStamp] = timestamp
	nonce, err := util.GenerateRandomString32()
	if err != nil {
		return "", err
	}
	paramMap[s.paramName.Nonce] = nonce
	// create sign
	sign, err := s.createSign(timestamp, nonce)
	if err != nil {
		return "", err
	}
	paramMap[s.paramName.Sign] = sign
	finalUrl := util.AddQueryMap(u, paramMap)
	return finalUrl, nil
}

// request send http request and use json.Unmarshal to converted to *model.Result.
func (s *SsoEnforcer) request(url string) (*model.Result, error) {
	resp, err := s.config.SendHttp(url)

	log.Printf("http request response: %s", resp)
	if err != nil {
		return nil, err
	}
	result := &model.Result{}
	err = json.Unmarshal([]byte(resp), result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *SsoEnforcer) GetAllowUrl() string {
	return s.config.AllowUrl
}

// saveTicket save ticket-id+client.
func (s *SsoEnforcer) saveTicket(ticket string, loginId string, client string) error {
	value := loginId
	if client != "" {
		value += "," + client
	}
	ticketTimeout := s.config.TicketTimeout
	return s.enforcer.GetAdapter().SetStr(s.spliceTicketSaveKey(ticket), value, ticketTimeout)
}

// delete ticket - id,client
func (s *SsoEnforcer) deleteTicket(ticket string) error {
	if ticket == "" {
		return nil
	}
	return s.enforcer.GetAdapter().DeleteStr(s.spliceTicketSaveKey(ticket))
}

// spliceTicketSaveKey splice ticket-id,client key.
func (s *SsoEnforcer) spliceTicketSaveKey(ticket string) string {
	return s.enforcer.GetTokenConfig().TokenName + ":ticket:" + ticket
}

// saveTicketIndex save id-ticket.
func (s *SsoEnforcer) saveTicketIndex(ticket string, id string) error {
	ticketTimeout := s.config.TicketTimeout
	return s.enforcer.GetAdapter().SetStr(s.spliceTicketIndexKey(id), ticket, ticketTimeout)
}

func (s *SsoEnforcer) deleteTicketIndex(id string) error {
	if id == "" {
		return nil
	}
	return s.enforcer.GetAdapter().DeleteStr(s.spliceTicketIndexKey(id))
}

// spliceTicketIndexKey splice id-ticket key.
func (s *SsoEnforcer) spliceTicketIndexKey(id string) string {
	return s.enforcer.GetTokenConfig().TokenName + ":id-ticket:" + id
}

// checkTimeStamp determine whether the gap between the timestamp and the current timestamp is within the allowable range.
func (s *SsoEnforcer) checkTimeStamp(timestamp string) error {
	parseInt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return err
	}
	if !s.isValidTimeStamp(parseInt) {
		return errors.New("timestamp is out of allowed range: " + timestamp)
	}

	return nil
}

// isValidTimeStamp determine whether the gap between the timestamp and the current timestamp is within the allowable range.
func (s *SsoEnforcer) isValidTimeStamp(timestamp int64) bool {
	allowDisparity := s.signConfig.TimeStampDisparity
	nowDisparity := time.Now().UnixMilli() - timestamp

	return allowDisparity == 1 || nowDisparity <= allowDisparity
}

// checkNonce the same nonce can only be verified once, cannot be used again for a period of time after use
func (s *SsoEnforcer) checkNonce(nonce string) error {
	if nonce == "" {
		return errors.New("nonce is nil")
	}
	// if nonce exists in adapter
	if !s.isValidNonce(nonce) {
		return errors.New("the nonce has been used: " + nonce)
	}
	// set nonce after nonce
	err := s.enforcer.GetAdapter().SetStr(s.spliceNonceSaveKey(nonce), nonce, s.signConfig.GetSaveNonceExpire()*2+2)
	if err != nil {
		return err
	}
	return nil
}

// isValidNonce determine random string, if not exist in adapter, return true.
func (s *SsoEnforcer) isValidNonce(nonce string) bool {
	if nonce == "" {
		return false
	}
	key := s.spliceNonceSaveKey(nonce)
	// if not exist in adapter, return true.
	return s.enforcer.GetAdapter().GetStr(key) == ""
}

func (s *SsoEnforcer) checkSign(timestamp string, nonce string, sign string) error {
	valid, err := s.isValidSign(timestamp, nonce, sign)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("invalid sign: " + sign)
	}
	return nil
}

// isValidSign use timestamp,nonce,sign to createSign to compare.
func (s *SsoEnforcer) isValidSign(timestamp string, nonce string, sign string) (bool, error) {
	recreateSign, err := s.createSign(timestamp, nonce)
	if err != nil {
		return false, err
	}
	return recreateSign == sign, nil
}

// createSign use util.MD5() to generate str.
func (s *SsoEnforcer) createSign(timestamp string, nonce string) (string, error) {
	secretKey := s.signConfig.SecretKey
	if secretKey == "" {
		return "", errors.New("please check SignConfig.SecretKey, SecretKey can not be nil")
	}
	str := s.paramName.Nonce + "=" + nonce + "&" + s.paramName.TimeStamp + "=" + timestamp + "&" + s.paramName.SecretKet + "=" + secretKey

	return util.MD5(str), nil
}

// spliceNonceSaveKey splice nonce store key.
func (s *SsoEnforcer) spliceNonceSaveKey(nonce string) string {
	return s.enforcer.GetTokenConfig().TokenName + ":sign:nonce:" + nonce
}
