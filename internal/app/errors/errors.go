package errors

import (
	"errors"
	"time"
)

type ErrorResponse struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

var ErrDuplicateShortId = errors.New("link with this shortId already exists")
var ErrShortIdLimitExceeded = errors.New("shortId generation limit exceeded")
