package auth

import (
	"context"
	"database/sql"
	"myproject/internal/models"
	"time"
)

type Repository interface {
	Save(reqCtx context.Context, username, password string) (int64, error)
	GetUserByName(reqCtx context.Context, username string) (models.User, error)
	GetUserById(reqCtx context.Context, userId int64) (models.User, error)
}

type DbRepository struct {
	Db *sql.DB
}

func (repo *DbRepository) Save(reqCtx context.Context, username, password string) (int64, error) {
	var newId int64

	execStr := "insert into users (name, password) values ($1, $2) returning id"

	ctx, cancel := context.WithTimeout(reqCtx, time.Second*2)
	defer cancel()

	execErr := repo.Db.QueryRowContext(ctx, execStr, username, password).Scan(&newId)
	return newId, execErr
}

func (repo *DbRepository) GetUserByName(reqCtx context.Context, username string) (models.User, error) {
	ctx, cancel := context.WithTimeout(reqCtx, time.Second*2)
	defer cancel()

	var selectedUser models.User

	sqlStr := `select u.id, u.name, u.password, u.created_at, 
       			r.id, r.code, r.name, r.created_at
			from users u
			left join roles r on u.role_id = r.id
			where u.name = $1`

	err := repo.Db.QueryRowContext(ctx, sqlStr, username).Scan(
		&selectedUser.Id,
		&selectedUser.Name,
		&selectedUser.Password,
		&selectedUser.CreatedAt,
		&selectedUser.Role.Id,
		&selectedUser.Role.Code,
		&selectedUser.Role.Name,
		&selectedUser.Role.CreatedAt,
	)

	if err != nil {
		return selectedUser, err
	}

	return selectedUser, nil
}

func (repo *DbRepository) GetUserById(reqCtx context.Context, userId int64) (models.User, error) {
	ctx, cancel := context.WithTimeout(reqCtx, time.Second*2)
	defer cancel()

	var selectedUser models.User

	sqlStr := `select u.id, u.name, u.password, u.created_at, 
       			r.id, r.code, r.name, r.created_at
			from users u
			left join roles r on u.role_id = r.id
			where u.id = $1`

	err := repo.Db.QueryRowContext(ctx, sqlStr, userId).Scan(
		&selectedUser.Id,
		&selectedUser.Name,
		&selectedUser.Password,
		&selectedUser.CreatedAt,
		&selectedUser.Role.Id,
		&selectedUser.Role.Code,
		&selectedUser.Role.Name,
		&selectedUser.Role.CreatedAt,
	)

	if err != nil {
		return selectedUser, err
	}

	return selectedUser, nil
}
