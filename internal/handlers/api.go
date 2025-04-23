package handlers

import (
	"encoding/json"
	"myproject/config"
	customerrors "myproject/internal/errors"
	"net/http"
)

func PingFunc(w http.ResponseWriter, r *http.Request, app config.CtxApp) {
	const op = "ping func"

	w.Header().Set("Content-Type", "application/json")
	message := map[string]string{
		"message": "pong",
	}

	response, _ := json.Marshal(message)
	_, err := w.Write(response)
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusInternalServerError, op)
	}
}
