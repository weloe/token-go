package errors

import (
	"errors"
)

var (
	BeReplaced = errors.New("this account is replaced")
	BeKicked   = errors.New("this account is kicked out")
	BeBanned   = errors.New("this account is banned")
)
