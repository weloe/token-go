package model

import (
	"container/list"
	"encoding/json"
	"fmt"
	"time"
)

type TokenSign struct {
	Value  string
	Device string
}

func (t *TokenSign) String() string {
	return fmt.Sprintf("Value: %s, Device: %s", t.Value, t.Device)
}

type Session struct {
	Id            string
	Type          string
	LoginType     string
	LoginId       string
	Token         string
	CreateTime    int64
	DataMap       map[string]interface{}
	TokenSignList []*TokenSign `json:"TokenSignList"`
}

func DefaultSession(id string) *Session {
	return &Session{
		Id:            id,
		DataMap:       make(map[string]interface{}),
		CreateTime:    time.Now().UnixMilli(),
		TokenSignList: make([]*TokenSign, 0),
	}
}

func NewSession(id string, sessionType string, loginId string) *Session {
	return &Session{
		Id:            id,
		Type:          sessionType,
		LoginId:       loginId,
		CreateTime:    time.Now().UnixMilli(),
		DataMap:       make(map[string]interface{}),
		TokenSignList: make([]*TokenSign, 0),
	}
}

// GetFilterTokenSign filter by TokenSign.Device from all TokenSign
func (s *Session) GetFilterTokenSign(device string) *list.List {
	copyList := list.New()
	for _, tokenSign := range s.TokenSignList {
		if tokenSign.Device == device {
			copyList.PushBack(tokenSign)
		}
	}
	return copyList
}

// GetTokenSignListCopy find all TokenSign
func (s *Session) GetTokenSignListCopy() *list.List {
	copyList := list.New()
	for _, tokenSign := range s.TokenSignList {
		copyList.PushBack(tokenSign)
	}
	return copyList
}

// GetTokenSign find TokenSign by TokenSign.Value
func (s *Session) GetTokenSign(tokenValue string) *TokenSign {
	if tokenValue == "" {
		return nil
	}
	for _, tokenSign := range s.TokenSignList {
		if tokenSign.Value == tokenValue {
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
	s.TokenSignList = append(s.TokenSignList, sign)
}

// RemoveTokenSign remove TokenSign by TokenSign.Value
func (s *Session) RemoveTokenSign(tokenValue string) bool {
	if tokenValue == "" {
		return false
	}
	for i, tokenSign := range s.TokenSignList {
		if tokenSign.Value == tokenValue {
			// delete
			s.RemoveTokenSignByIndex(i)
			return true
		}
	}
	return true
}

// RemoveTokenSignByIndex delete by index
func (s *Session) RemoveTokenSignByIndex(i int) {
	s.TokenSignList = append(s.TokenSignList[:i], s.TokenSignList[i+1:]...)
}

// GetLastTokenByDevice get TokenSign.Value by device
func (s *Session) GetLastTokenByDevice(device string) string {
	tokenSignList := s.GetFilterTokenSign(device)
	if tokenSign, ok := tokenSignList.Back().Value.(*TokenSign); ok {
		return tokenSign.Value
	}
	return ""
}

// TokenSignSize get tokenSign size
func (s *Session) TokenSignSize() int {
	return len(s.TokenSignList)
}

// Json return json string
func (s *Session) Json() string {
	b, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(b)
}

// Get returns data from DataMap
func (s *Session) Get(key string) interface{} {
	value, ok := s.DataMap[key]
	if !ok {
		return nil
	}
	return value
}

func (s *Session) Set(key string, obj interface{}) {
	s.DataMap[key] = obj
}

// GetOrSet returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (s *Session) GetOrSet(key string, obj interface{}) (interface{}, bool) {
	value := s.Get(key)
	if value == nil {
		s.Set(key, obj)
		return obj, false
	}
	return value, true
}

func (s *Session) String() string {
	tokenSigns := ""
	for _, ts := range s.TokenSignList {
		tokenSigns += ts.String() + "\n"
	}

	return fmt.Sprintf("Id: %s, Type: %s, LoginType: %s, LoginId: %s, Token: %s, CreateTime: %d, DataMap: %+v, \nTokenSignList:\n%s",
		s.Id, s.Type, s.LoginType, s.LoginId, s.Token, s.CreateTime, s.DataMap, tokenSigns)
}
