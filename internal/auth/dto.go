package auth

import "errors"

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
