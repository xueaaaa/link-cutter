package errors

import "errors"

var ErrDuplicateShortId = errors.New("link with this shortId already exists")
var ErrShortIdLimitExceeded = errors.New("shortId generation limit exceeded")
