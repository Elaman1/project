package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"myproject/config"
	"myproject/internal/dto"
	customerrors "myproject/internal/errors"
	"myproject/internal/lib"
	"myproject/internal/models"
	"myproject/pkg/passwordhasher"
	"net/http"
	"time"
)

func validateUserRequest(w http.ResponseWriter, r *http.Request, op string) (dto.CreateUserRequest, bool) {
	var user dto.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusBadRequest, op)
		return user, false
	}

	err = user.Validate()
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusBadRequest, op)
		return user, false
	}

	return user, true
}

func generateSessionId(login string) string {
	return fmt.Sprintf("%s-%d", login, time.Now().Unix())
}

func Register(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	const op = "register function"

	user, ok := validateUserRequest(w, r, op)
	if !ok {
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
		if lib.IsUniqueViolation(execErr) {
			execErr = errors.New("пользователь уже существует")
		}

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

func Login(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	const op = "login function"

	user, ok := validateUserRequest(w, r, op)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	var selectedUser models.User
	sqlStr := "select id, name, password, created_at from users where name = $1"

	err := db.QueryRowContext(ctx, sqlStr, user.Username).Scan(&selectedUser.Id, &selectedUser.Name, &selectedUser.Password, &selectedUser.CreatedAt)
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusBadRequest, op)
		return
	}

	checkPassword, err := passwordhasher.CheckPassword(user.Password, selectedUser.Password)
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusBadRequest, op)
		return
	}

	if !checkPassword {
		err = errors.New("пароль или логин непавильный")
		customerrors.HandleJsonErrors(w, err, http.StatusBadRequest, op)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	sessionId := generateSessionId(user.Username)
	config.Sessions[sessionId] = config.Session{
		Name:      user.Username,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	fmt.Println("Устанавливается сессия: ", sessionId)
	http.SetCookie(w, &http.Cookie{
		Name:     "sessionId",
		Value:    sessionId,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
	})

	response := map[string]string{
		"message": "Вы авторизовались",
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusBadRequest, op)
		return
	}

	log.Println("Пользователь Авторизован: ", user.Username)
}

func Protected(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	const op = "protected function"

	fmt.Println(r.Context().Value("user"))
}
