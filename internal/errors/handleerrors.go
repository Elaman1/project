package errors

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func HandleJsonErrors(w http.ResponseWriter, err error, code int, op string) {
	if err == nil {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	resp := ErrorResponse{
		Message: err.Error(),
		Error:   http.StatusText(code),
	}

	log.Printf("Ошибка [%d - %s] Текст ошибки: %s в %s", code, http.StatusText(code), err.Error(), op)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		return
	}
}
