package model

import (
	"container/list"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// SyncSession sync SyncSession
type SyncSession struct {
	Id            string
	Type          string
	LoginType     string
	LoginId       string
	Token         string
	CreateTime    int64
	DataMap       *sync.Map
	TokenSignList []*TokenSign `json:"TokenSignList"`
}

func DefaultSyncSession(id string) *SyncSession {
	return &SyncSession{
		Id:            id,
		DataMap:       &sync.Map{},
		CreateTime:    time.Now().UnixMilli(),
		TokenSignList: make([]*TokenSign, 0),
	}
}

func NewSyncSession(id string, sessionType string, loginId string) *SyncSession {
	return &SyncSession{
		Id:            id,
		Type:          sessionType,
		LoginId:       loginId,
		CreateTime:    time.Now().UnixMilli(),
		DataMap:       &sync.Map{},
		TokenSignList: make([]*TokenSign, 0),
	}
}

// GetFilterTokenSign filter by TokenSign.Device from all TokenSign
func (s *SyncSession) GetFilterTokenSign(device string) *list.List {
	if device == "" {
		return s.GetTokenSignListCopy()
	}
	copyList := list.New()
	for _, tokenSign := range s.TokenSignList {
		if tokenSign.Device == device {
			copyList.PushBack(tokenSign)
		}
	}
	return copyList
}

// GetTokenSignListCopy find all TokenSign
func (s *SyncSession) GetTokenSignListCopy() *list.List {
	copyList := list.New()
	for _, tokenSign := range s.TokenSignList {
		copyList.PushBack(tokenSign)
	}
	return copyList
}

// GetTokenSign find TokenSign by TokenSign.Value
func (s *SyncSession) GetTokenSign(tokenValue string) *TokenSign {
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
func (s *SyncSession) AddTokenSign(sign *TokenSign) {
	if s.GetTokenSign(sign.Value) != nil {
		return
	}
	s.TokenSignList = append(s.TokenSignList, sign)
}

// RemoveTokenSign remove TokenSign by TokenSign.Value
func (s *SyncSession) RemoveTokenSign(tokenValue string) bool {
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
func (s *SyncSession) RemoveTokenSignByIndex(i int) {
	s.TokenSignList = append(s.TokenSignList[:i], s.TokenSignList[i+1:]...)
}

// GetLastTokenByDevice get TokenSign.Value by device
func (s *SyncSession) GetLastTokenByDevice(device string) string {
	if device == "" {
		return ""
	}
	tokenSignList := s.GetFilterTokenSign(device)
	if tokenSign, ok := tokenSignList.Back().Value.(*TokenSign); ok && tokenSign.Device == device {
		return tokenSign.Value
	}
	return ""
}

// TokenSignSize get tokenSign size
func (s *SyncSession) TokenSignSize() int {
	return len(s.TokenSignList)
}

// Json return json string
func (s *SyncSession) Json() string {
	myStruct := &Session{
		Id:            s.Id,
		Type:          s.Type,
		LoginType:     s.LoginType,
		LoginId:       s.LoginId,
		Token:         s.Token,
		CreateTime:    s.CreateTime,
		TokenSignList: s.TokenSignList,
		DataMap:       make(map[string]interface{}),
	}

	s.DataMap.Range(func(key, value any) bool {
		myStruct.DataMap[fmt.Sprintf("%v", key)] = value
		return true
	})

	b, err := json.Marshal(myStruct)
	if err != nil {
		return ""
	}
	return string(b)
}

// UnmarshalBytes convert bytes to SyncSession
func (s *SyncSession) UnmarshalBytes(jsonByte []byte) (*SyncSession, error) {
	return JsonByteToSyncSession(jsonByte)
}

// UnmarshalStr convert string to SyncSession
func (s *SyncSession) UnmarshalStr(jsonStr string) (*SyncSession, error) {
	return JsonToSyncSession(jsonStr)
}

// Get returns data from DataMap
func (s *SyncSession) Get(key string) interface{} {
	value, ok := s.DataMap.Load(key)
	if !ok {
		return nil
	}
	return value
}

func (s *SyncSession) Set(key string, obj interface{}) {
	s.DataMap.Store(key, obj)
}

// GetOrSet returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (s *SyncSession) GetOrSet(key string, obj interface{}) (interface{}, bool) {
	return s.DataMap.LoadOrStore(key, obj)
}
