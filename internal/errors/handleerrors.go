package errors

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

const NotAllowed = "Not Allowed"

type MiddlewareNotAllowedError struct {
	Code    int
	Message string
}

func (e MiddlewareNotAllowedError) Error() string {
	return fmt.Sprintf("Код: %d Ошибка: %s", e.Code, e.Message)
}

func HandleJsonErrors(w http.ResponseWriter, err error, code int) {
	if err == nil {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	resp := ErrorResponse{
		Message: err.Error(),
		Error:   http.StatusText(code),
	}

	log.Printf("Ошибка [%d - %s] Текст ошибки: %s", code, http.StatusText(code), err.Error())

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		return
	}
}
