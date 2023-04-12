package constants

import (
	"errors"
)

var (
	ErrObjectNil     = errors.New("object is nil")
	ErrInvalidMethod = errors.New("method not allowed")
	ErrObjectExists  = errors.New("object already exists")
	ErrNoSuchObject  = errors.New("object does not exist")
	ErrInvalidFilter = errors.New("invalid filter key")
	ErrInvalidPath   = errors.New("invalid request path")
)
