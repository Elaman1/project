package middlewares

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"myproject/config"
	customerrors "myproject/internal/errors"
	"myproject/internal/functions"
	"net/http"
	"time"
)

type Auth struct {
}

func (a *Auth) Handle(next functions.CustomHttpHandler) functions.CustomHttpHandler {
	const op = "middleware auth"
	return func(w http.ResponseWriter, r *http.Request, db *sql.DB) {
		sessionId, err := r.Cookie("sessionId")
		if err != nil {
			customerrors.HandleJsonErrors(w, errors.New("unauthorized"), http.StatusUnauthorized, op)
			return
		}

		session, exists := config.Sessions[sessionId.Value]
		fmt.Println("sessions:", config.Sessions)
		if !exists {
			customerrors.HandleJsonErrors(w, errors.New("unauthorized"), http.StatusUnauthorized, op)
			return
		}

		if session.ExpiresAt.Before(time.Now()) {
			err = errors.New("session expired")
			customerrors.HandleJsonErrors(w, err, http.StatusUnauthorized, op)
			return
		}

		ctx := context.WithValue(r.Context(), config.CtxUserKey, session.Name)
		fmt.Println("Найден пользователь: ", sessionId.Value)
		next(w, r.WithContext(ctx), db)
	}
}

func (a *Auth) Err() error {
	return nil
}
