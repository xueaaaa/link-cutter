package util

import (
	"encoding/json"
	errors2 "link-cutter/internal/app/errors"
	"net/http"
	"time"
)

func WriteError(w http.ResponseWriter, statusCode int, msg string, reqId string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonErr := errors2.ErrorResponse{
		RequestId: reqId,
		Message:   msg,
		Timestamp: time.Now().UTC(),
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	encoder.Encode(jsonErr)
}
