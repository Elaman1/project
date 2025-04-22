package auth

import (
	"context"
	"database/sql"
	"myproject/internal/models"
	"time"
)

type Repository interface {
	Save(username, password string) (int64, error)
	GetUserByName(username string) (models.User, error)
}

type DbRepository struct {
	Db     *sql.DB
	ReqCtx context.Context
}

func (repo *DbRepository) Save(username, password string) (int64, error) {
	var newId int64

	execStr := "insert into users (name, password) values ($1, $2) returning id"

	ctx, cancel := context.WithTimeout(repo.ReqCtx, time.Second*2)
	defer cancel()

	execErr := repo.Db.QueryRowContext(ctx, execStr, username, password).Scan(&newId)
	return newId, execErr
}

func (repo *DbRepository) GetUserByName(username string) (models.User, error) {
	ctx, cancel := context.WithTimeout(repo.ReqCtx, time.Second*2)
	defer cancel()

	var selectedUser models.User

	sqlStr := "select id, name, password, created_at from users where name = $1"

	err := repo.Db.QueryRowContext(ctx, sqlStr, username).Scan(&selectedUser.Id, &selectedUser.Name, &selectedUser.Password, &selectedUser.CreatedAt)
	if err != nil {
		return selectedUser, err
	}

	return selectedUser, nil
}
