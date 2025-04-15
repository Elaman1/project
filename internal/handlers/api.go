package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	customerrors "myproject/internal/errors"
	"net/http"
)

func PingFunc(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	const op = "ping func"

	w.Header().Set("Content-Type", "application/json")
	message := map[string]string{
		"message": "pong",
	}

	response, _ := json.Marshal(message)
	_, err := w.Write(response)
	if err != nil {
		returnErr := fmt.Errorf("вышла ошибка %w, Место %s", err, op)
		customerrors.HandleJsonErrors(w, returnErr, http.StatusInternalServerError)
	}
}
