package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"myproject/config"
	customerrors "myproject/internal/errors"
	"net/http"
	"time"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request, ctxApp config.CtxApp) {
	const op = "register function"

	user, err := validateUserRequest(w, r, op)
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusBadRequest, op)
		return
	}

	service := Service{
		Repo: &DbRepository{
			Db: ctxApp.Db,
		},
	}

	execErr := service.Registration(r.Context(), user.Username, user.Password)
	if execErr != nil {
		customerrors.HandleJsonErrors(w, execErr, http.StatusBadRequest, op)
		return
	}

	response := RegisterUserResponse{
		Message: "Успешно добавлено",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusBadRequest, op)
		return
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request, ctxApp config.CtxApp) {
	const op = "login function"

	user, err := validateUserRequest(w, r, op)
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusBadRequest, op)
		return
	}

	service := Service{
		Repo: &DbRepository{
			Db: ctxApp.Db,
		},
	}

	selectedUser, err := service.Login(r.Context(), user.Username, user.Password)
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusBadRequest, op)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	sessionId := generateSessionId(user.Username)
	config.SessionsMu.Lock()
	config.Sessions[sessionId] = config.Session{
		UserId:    selectedUser.Id,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	config.SessionsMu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionId",
		Value:    sessionId,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
	})

	response := LoginUserResponse{
		Message: "Вы авторизовались",
		Token:   sessionId,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusBadRequest, op)
		return
	}

	log.Println("Пользователь Авторизован: ", user.Username)
}

func MeHandler(w http.ResponseWriter, r *http.Request, ctxApp config.CtxApp) {
	const op = "me function"

	user := r.Context().Value(config.CtxUserKey)

	fmt.Println("user: ", user)
}

func validateUserRequest(w http.ResponseWriter, r *http.Request, op string) (CreateUserRequest, error) {
	var user CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return user, err
	}

	err = user.Validate()
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusBadRequest, op)
		return user, err
	}

	return user, nil
}

func generateSessionId(login string) string {
	return fmt.Sprintf("%s-%d", login, time.Now().Unix())
}
