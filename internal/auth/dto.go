package auth

import (
	"errors"
	"myproject/internal/models"
	"sync"
	"time"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (req *CreateUserRequest) Validate() error {
	if req.Username == "" {
		return errors.New("имя пользователя не заполнено")
	}

	if req.Password == "" {
		return errors.New("пароль не заполнено")
	}

	if len(req.Password) < 8 {
		return errors.New("пароль слишком короткий")
	}

	if len(req.Username) > 100 {
		return errors.New("имя пользователя слишком длинное")
	}

	if len(req.Password) > 500 {
		return errors.New("пароль слишком длинный")
	}

	return nil
}

type RegisterUserResponse struct {
	Message string `json:"message"`
}

type LoginUserResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type CachedUser struct {
	User      *models.User
	ExpiresAt time.Time
}

func (user *CachedUser) Expired() bool {
	return time.Now().After(user.ExpiresAt)
}

var UserCaches = make(map[int64]CachedUser)

var UserCachesMu sync.RWMutex
