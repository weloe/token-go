package model

import (
	"encoding/json"
	"sync"
)

func JsonByteToSyncSession(jsonByte []byte) (*SyncSession, error) {
	s := &Session{}
	err := json.Unmarshal(jsonByte, s)
	if err != nil {
		return nil, err
	}

	return ConvertSyncSession(s), nil
}

func JsonToSyncSession(jsonStr string) (*SyncSession, error) {
	return JsonByteToSyncSession([]byte(jsonStr))
}

// ConvertSyncSession convert to Session
func ConvertSyncSession(s *Session) *SyncSession {
	session := &SyncSession{
		Id:            s.Id,
		Type:          s.Type,
		LoginType:     s.LoginType,
		LoginId:       s.LoginId,
		Token:         s.Token,
		CreateTime:    s.CreateTime,
		DataMap:       &sync.Map{},
		TokenSignList: s.TokenSignList,
	}
	// copy DataMap
	for k, v := range s.DataMap {
		session.DataMap.Store(k, v)
	}
	return session
}
