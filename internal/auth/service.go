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

func (s *Service) Registration(ctx context.Context, username, password, role string) error {
	hashedPassword, err := passwordhasher.HashPassword(password)
	if err != nil {
		return err
	}

	_, err = s.Repo.Save(ctx, username, hashedPassword, role)
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

func (s *Service) GetAllUsers(ctx context.Context) ([]models.User, error) {
	selectedUsers, err := s.Repo.GetAllUsers(ctx)
	if err != nil {
		return selectedUsers, err
	}

	return selectedUsers, nil
}

func (s *Service) DeleteUser(ctx context.Context, userId int64) error {
	return s.Repo.DeleteUser(ctx, userId)
}

func (s *Service) ChangeBlockUser(ctx context.Context, userId int64) error {
	return s.Repo.ChangeBlockUser(ctx, userId)
}
