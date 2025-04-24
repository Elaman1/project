package middlewares

import (
	"context"
	"errors"
	"fmt"
	"myproject/config"
	"myproject/internal/auth"
	customerrors "myproject/internal/errors"
	"myproject/internal/functions"
	"myproject/internal/models"
	"net/http"
	"time"
)

type Auth struct {
}

func (a *Auth) Handle(next functions.CustomHttpHandler) functions.CustomHttpHandler {
	const op = "middleware auth"
	return func(w http.ResponseWriter, r *http.Request, ctxApp config.CtxApp) {
		sessionId, err := r.Cookie("sessionId")
		if err != nil {
			customerrors.HandleJsonErrors(w, errors.New("unauthorized"), http.StatusUnauthorized, op)
			return
		}

		config.SessionsMu.RLock()
		session, exists := config.Sessions[sessionId.Value]
		config.SessionsMu.RUnlock()
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

		user, err := a.GetUser(r.Context(), session.UserId, ctxApp)
		if err != nil {
			customerrors.HandleJsonErrors(w, err, http.StatusUnauthorized, op)
			return
		}

		ctx := context.WithValue(r.Context(), config.CtxUserKey, &user)
		fmt.Println("Найден пользователь: ", sessionId.Value)
		next(w, r.WithContext(ctx), ctxApp)
	}
}

func (a *Auth) Err() error {
	return nil
}

func (a *Auth) GetUser(ctx context.Context, userId int64, ctxApp config.CtxApp) (models.User, error) {
	auth.UserCachesMu.RLock()
	cachedUser, ok := auth.UserCaches[userId]
	auth.UserCachesMu.RUnlock()
	// Если нашлось и не истек возвращаем по кэшу
	if ok && !cachedUser.Expired() {
		return *cachedUser.User, nil
	}

	// Если не нашлось или истек пробуем обращаться к базе
	service := auth.Service{
		Repo: &auth.DbRepository{
			Db: ctxApp.Db,
		},
	}

	resUser, err := service.Repo.GetUserById(ctx, userId)
	if err != nil {
		return resUser, err
	}

	auth.UserCachesMu.Lock()
	defer auth.UserCachesMu.Unlock()
	auth.UserCaches[userId] = auth.CachedUser{
		User:      &resUser,
		ExpiresAt: time.Now().Add(time.Minute * 5),
	}

	return resUser, nil

}
