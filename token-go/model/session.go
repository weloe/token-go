package model

import (
	"container/list"
	"sync"
	"time"
)

type TokenSign struct {
	Value  string
	Device string
}

type Session struct {
	Id            string
	Type          string
	LoginType     string
	LoginId       string
	Token         string
	CreateTime    int64
	DataMap       *sync.Map
	TokenSignList *list.List
}

func DefaultSession(id string) *Session {
	return &Session{
		Id:         id,
		CreateTime: time.Now().UnixMilli(),
	}
}

func NewSession(id string, sessionType string, loginId string) *Session {
	return &Session{
		Id:            id,
		Type:          sessionType,
		LoginId:       loginId,
		CreateTime:    time.Now().UnixMilli(),
		TokenSignList: list.New(),
	}
}

// GetFilterTokenSign filter by TokenSign.Device from all TokenSign
func (s *Session) GetFilterTokenSign(device string) *list.List {
	if device == "" {
		return s.GetTokenSignListCopy()
	}
	copyList := list.New()
	for e := s.TokenSignList.Front(); e != nil; e = e.Next() {
		if tokenSign, ok := e.Value.(*TokenSign); ok && tokenSign.Device == device {
			copyList.PushBack(tokenSign)
		}
	}
	return copyList
}

// GetTokenSignListCopy find all TokenSign
func (s *Session) GetTokenSignListCopy() *list.List {
	copyList := list.New()
	for e := s.TokenSignList.Front(); e != nil; e = e.Next() {
		copyList.PushBack(e.Value)
	}
	return copyList
}

// GetTokenSign find TokenSign by TokenSign.Value
func (s *Session) GetTokenSign(tokenValue string) *TokenSign {
	if tokenValue == "" {
		return nil
	}
	for e := s.TokenSignList.Front(); e != nil; e = e.Next() {
		if tokenSign, ok := e.Value.(*TokenSign); ok && tokenSign.Value == tokenValue {
			return tokenSign
		}
	}
	return nil
}

// AddTokenSign add TokenSign
func (s *Session) AddTokenSign(sign *TokenSign) {
	if s.GetTokenSign(sign.Value) != nil {
		return
	}
	s.TokenSignList.PushBack(sign)
}

// RemoveTokenSign remove TokenSign by TokenSign.Value
func (s *Session) RemoveTokenSign(tokenValue string) bool {
	if tokenValue == "" {
		return false
	}
	for e := s.TokenSignList.Front(); e != nil; e = e.Next() {
		if tokenSign, ok := e.Value.(*TokenSign); ok && tokenSign.Value == tokenValue {
			s.TokenSignList.Remove(e)
		}
	}
	return true
}

// GetLastTokenByDevice get TokenSign.Value by device
func (s *Session) GetLastTokenByDevice(device string) string {
	if device == "" {
		return ""
	}
	tokenSignList := s.GetFilterTokenSign(device)
	if tokenSign, ok := tokenSignList.Back().Value.(*TokenSign); ok && tokenSign.Device == device {
		return tokenSign.Value
	}
	return ""
}
