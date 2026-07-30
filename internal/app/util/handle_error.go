package util

import (
	"errors"
	errors2 "link-cutter/internal/app/errors"
	"net/http"
)

func HandleError(w http.ResponseWriter, r *http.Request, err error, defaultCode ...int) {
	var apiErr errors2.APIError

	if errors.As(err, &apiErr) {
		WriteError(
			w,
			apiErr.Code,
			err.Error(),
			GetRequestId(r),
		)
	} else {
		code := http.StatusInternalServerError
		if len(defaultCode) > 0 {
			code = defaultCode[0]
		}

		WriteError(
			w,
			code,
			err.Error(),
			GetRequestId(r),
		)
	}
}
