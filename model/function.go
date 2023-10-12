package model

import (
	"errors"
	"github.com/weloe/token-go/constant"
	"github.com/weloe/token-go/util"
	"sync"
)

type GenerateFunc func() (string, error)

type GenerateTokenFunc struct {
	fns *sync.Map
}

func (g *GenerateTokenFunc) Exec(tokenForm string) (string, error) {
	if tokenForm == "" {
		return "", errors.New("GenerateToken failed: tokenStyle is nil")
	}
	handlerFunc, err := g.GetFunction(tokenForm)
	if err != nil {
		return "", err
	}
	s, err := handlerFunc()
	if err != nil {
		return "", err
	}
	return s, nil
}

func (g *GenerateTokenFunc) GetFunction(tokenForm string) (GenerateFunc, error) {
	value, ok := g.fns.Load(tokenForm)
	if !ok {
		return nil, errors.New("GetFunction() failed: load func error")
	}
	if value == nil {
		return nil, errors.New("GetFunction() failed: this tokenStyle generate func doesn't exist")
	}
	handlerFunc := value.(GenerateFunc)
	return handlerFunc, nil
}

func LoadFunctionMap() GenerateTokenFunc {
	fm := &GenerateTokenFunc{fns: &sync.Map{}}
	fm.AddFunc(constant.UUID, util.GenerateUUID)
	fm.AddFunc(constant.SimpleUUID, util.GenerateSimpleUUID)
	fm.AddFunc(constant.RandomString32, util.GenerateRandomString32)
	fm.AddFunc(constant.RandomString64, util.GenerateRandomString64)
	fm.AddFunc(constant.RandomString128, util.GenerateRandomString128)

	return *fm
}

func (g *GenerateTokenFunc) AddFunc(key string, f GenerateFunc) {
	g.fns.LoadOrStore(key, f)
}
