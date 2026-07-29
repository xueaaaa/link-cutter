package errors

import (
	"errors"
	"time"
)

type ErrorResponse struct {
	RequestId string    `json:"reqId"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

var ErrDuplicateShortId = errors.New("link with this shortId already exists")
var ErrShortIdLimitExceeded = errors.New("shortId generation limit exceeded")
var ErrLinkNotFound = errors.New("link not found")
var ErrMissingAuthHeader = errors.New("missing authorization header")
var ErrInvalidAuthHeader = errors.New("invalid authorization header")
var ErrUnexpectedSigningMethod = errors.New("unexpected jwt signing method")
var ErrUserNotFound = errors.New("user not found")
var ErrNotEnoughRights = errors.New("not enough rights")
