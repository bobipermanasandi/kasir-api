package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)


type APIResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}


func writeJSON(w http.ResponseWriter, httpStatus int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	resp := APIResponse{
		Code:    httpStatus,
		Status:  http.StatusText(httpStatus),
		Message: message,
		Data:    data,
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, msg, nil)
}

func decodeJSON(r *http.Request, v any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

func getIDFromURL(r *http.Request, basePath string) (int, error) {
	idStr := strings.TrimPrefix(r.URL.Path, basePath)
	return strconv.Atoi(idStr)
}

func mustGetID(w http.ResponseWriter, r *http.Request, basePath string) (int, bool) {
	id, err := getIDFromURL(r, basePath)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid ID")
		return 0, false
	}
	return id, true
}
