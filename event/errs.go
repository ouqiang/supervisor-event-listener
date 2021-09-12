package event

import (
	"errors"
)

var (
	ErrParseHeader   = errors.New("parse header failed")
	ErrParsePayload  = errors.New("parse payload failed")
	ErrPayloadLength = errors.New("invalid payload length")
)
