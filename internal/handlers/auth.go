package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"myproject/internal/dto"
	customerrors "myproject/internal/errors"
	"myproject/internal/lib"
	"myproject/pkg/passwordhasher"
	"net/http"
	"time"
)

func Register(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	const op = "register function"

	var user dto.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusBadRequest, op)
		return
	}

	err = user.Validate()
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusBadRequest, op)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	var newId int64

	execStr := "insert into users (name, password) values ($1, $2) returning id"

	hashedPassword, err := passwordhasher.HashPassword(user.Password)
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusBadRequest, op)
		return
	}

	execErr := db.QueryRowContext(ctx, execStr, user.Username, hashedPassword).Scan(&newId)
	if execErr != nil {
		fmt.Println("before error:", execErr)
		if lib.IsUniqueViolation(execErr) {
			execErr = errors.New("пользователь уже существует")
		}
		fmt.Println("after error:", execErr.Error())

		customerrors.HandleJsonErrors(w, execErr, http.StatusBadRequest, op)
		return
	}

	response := map[string]any{
		"message": "Успешно добавлено",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)

	log.Printf("Добавлен новый пользователь: %d", newId)
}
