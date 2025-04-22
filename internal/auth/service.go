package auth

import (
	"errors"
	"myproject/internal/lib"
	"myproject/pkg/passwordhasher"
)

type Service struct {
	Rep Repository
}

func (s *Service) IsRouteService() {
	// Затычка для определения структуры
}

func (s *Service) Registration(username, password string) error {
	hashedPassword, err := passwordhasher.HashPassword(password)
	if err != nil {
		return err
	}

	_, err = s.Rep.Save(username, hashedPassword)
	if lib.IsUniqueViolation(err) {
		return errors.New("пользователь уже существует")
	}

	return nil
}

func (s *Service) Login(username, password string) error {
	selectedUser, err := s.Rep.GetUserByName(username)

	checkPassword, err := passwordhasher.CheckPassword(password, selectedUser.Password)
	if err != nil {
		return err
	}

	if !checkPassword {
		return errors.New("пароль или логин непавильный")
	}

	return nil
}
