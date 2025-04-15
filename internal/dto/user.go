package dto

import "errors"

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (req *CreateUserRequest) Validate() error {
	if req.Username == "" {
		return errors.New("username is empty")
	}

	if req.Password == "" {
		return errors.New("password is empty")
	}

	if len(req.Password) < 8 {
		return errors.New("password is too short")
	}

	if len(req.Username) > 100 {
		return errors.New("username is too long")
	}

	if len(req.Password) < 2 {
		return errors.New("password is too short")
	}

	return nil
}
