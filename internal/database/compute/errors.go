package compute

import "errors"

var (
	ErrInvalidSymbol    = errors.New("invalid symbol")
	ErrInvalidCommand   = errors.New("invalid command")
	ErrInvalidArguments = errors.New("invalid arguments")
)
