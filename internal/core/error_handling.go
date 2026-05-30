package core

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Details string `json:"details,omitempty"`
}

func NewErrorResponse(writer http.ResponseWriter, error string, code int, details string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)

	res := &ErrorResponse{
		Error:   error,
		Code:    code,
		Status:  http.StatusText(code),
		Details: details,
	}
	_ = json.NewEncoder(writer).Encode(res)
}
