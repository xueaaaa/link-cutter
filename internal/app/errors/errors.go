package errors

import (
	"net/http"
	"time"
)

type ErrorResponse struct {
	RequestId string    `json:"reqId"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type APIError struct {
	Code    int
	Message string
}

func (e APIError) Error() string {
	return e.Message
}

var ErrDuplicateShortId = APIError{
	Code:    http.StatusInternalServerError,
	Message: "link with this shortId already exists",
}
var ErrShortIdLimitExceeded = APIError{
	Code:    http.StatusInternalServerError,
	Message: "shortId generation limit exceeded",
}
var ErrLinkNotFound = APIError{
	Code:    http.StatusNotFound,
	Message: "link not found",
}
var ErrMissingAuthHeader = APIError{
	Code:    http.StatusUnauthorized,
	Message: "missing authorization header",
}
var ErrInvalidAuthHeader = APIError{
	Code:    http.StatusUnauthorized,
	Message: "invalid authorization header",
}
var ErrUnexpectedSigningMethod = APIError{
	Code:    http.StatusUnauthorized,
	Message: "unexpected jwt signing method",
}
var ErrUserNotFound = APIError{
	Code:    http.StatusNotFound,
	Message: "user not found",
}
var ErrNotEnoughRights = APIError{
	Code:    http.StatusForbidden,
	Message: "not enough rights",
}
