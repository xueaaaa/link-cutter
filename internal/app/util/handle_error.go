package util

import (
	"errors"
	errors2 "link-cutter/internal/app/errors"
	"net/http"

	"go.uber.org/zap"
)

func HandleError(w http.ResponseWriter, r *http.Request, logger *zap.Logger, err error, defaultCode ...int) {
	var apiErr errors2.APIError

	logger.Error(err.Error(),
		zap.String("req_id", GetRequestId(r)),
	)
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
