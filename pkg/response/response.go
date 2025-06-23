package response

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

type ErrorResponse struct {
	Success   bool   `json:"success"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
	Details   any    `json:"details,omitempty"`
}

func RespondSuccess(w http.ResponseWriter, code int, msg string, data any) {
	response := SuccessResponse{Success: true, Data: data, Message: msg}
	respond(w, code, response)
}

func RespondError(w http.ResponseWriter, status int, code, msg string, details any) {
	response := ErrorResponse{Success: false, Code: status, Message: msg, ErrorCode: code, Details: details}
	respond(w, status, response)
}

func respond[T any](w http.ResponseWriter, statusCode int, data T) {
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
