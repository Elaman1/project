package auth

import (
	"context"
	"errors"
	"myproject/internal/lib"
	"myproject/internal/models"
	"myproject/pkg/passwordhasher"
)

type Service struct {
	Repo Repository
}

func (s *Service) Registration(ctx context.Context, username, password string) error {
	hashedPassword, err := passwordhasher.HashPassword(password)
	if err != nil {
		return err
	}

	_, err = s.Repo.Save(ctx, username, hashedPassword)
	if lib.IsUniqueViolation(err) {
		return errors.New("пользователь уже существует")
	}

	return nil
}

func (s *Service) Login(ctx context.Context, username, password string) (models.User, error) {
	selectedUser, err := s.Repo.GetUserByName(ctx, username)
	if err != nil {
		return selectedUser, err
	}

	checkPassword, err := passwordhasher.CheckPassword(password, selectedUser.Password)
	if err != nil {
		return selectedUser, err
	}

	if !checkPassword {
		return selectedUser, errors.New("неверный логин или пароль")
	}

	return selectedUser, nil
}
