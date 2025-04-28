package auth

import (
	"context"
	"database/sql"
	"errors"
	"myproject/internal/models"
	"time"
)

type Repository interface {
	Save(reqCtx context.Context, username, password, roleCode string) (int64, error)
	GetUserByName(reqCtx context.Context, username string) (models.User, error)
	GetUserById(reqCtx context.Context, userId int64) (models.User, error)
	GetAllUsers(reqCtx context.Context) ([]models.User, error)
	DeleteUser(reqCtx context.Context, userId int64) error
	ChangeBlockUser(reqCtx context.Context, userId int64) error
}

type DbRepository struct {
	Db *sql.DB
}

func (repo *DbRepository) Save(reqCtx context.Context, username, password, roleCode string) (int64, error) {
	var newId int64

	// Получаем все роли чтобы определить ИД роли
	roles, err := repo.getAllRoles(reqCtx)
	if err != nil {
		return newId, err
	}

	role, ok := roles[roleCode]
	if !ok {
		return newId, errors.New("role not found")
	}

	execStr := "insert into users (name, password, role_id) values ($1, $2, $3) returning id"

	ctx, cancel := context.WithTimeout(reqCtx, time.Second*2)
	defer cancel()

	execErr := repo.Db.QueryRowContext(ctx, execStr, username, password, role.Id).Scan(&newId)
	return newId, execErr
}

func (repo *DbRepository) GetUserByName(reqCtx context.Context, username string) (models.User, error) {
	ctx, cancel := context.WithTimeout(reqCtx, time.Second*2)
	defer cancel()

	var selectedUser models.User

	sqlStr := `select u.id, u.name, u.password, u.created_at, u.role_id,
       			r.id, r.code, r.name, r.created_at
			from users u
			left join roles r on u.role_id = r.id
			where u.name = $1`

	err := repo.Db.QueryRowContext(ctx, sqlStr, username).Scan(
		&selectedUser.Id,
		&selectedUser.Name,
		&selectedUser.Password,
		&selectedUser.CreatedAt,
		&selectedUser.RoleID,
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

	sqlStr := `select u.id, u.name, u.password, u.created_at, u.role_id, u.blocked,
       			r.id, r.code, r.name, r.created_at
			from users u
			left join roles r on u.role_id = r.id
			where u.id = $1`

	err := repo.Db.QueryRowContext(ctx, sqlStr, userId).Scan(
		&selectedUser.Id,
		&selectedUser.Name,
		&selectedUser.Password,
		&selectedUser.CreatedAt,
		&selectedUser.RoleID,
		&selectedUser.Blocked,
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

func (repo *DbRepository) getAllRoles(reqCtx context.Context) (map[string]models.Role, error) {
	ctx, cancel := context.WithTimeout(reqCtx, time.Second*2)
	defer cancel()
	var roles = make(map[string]models.Role)

	RolesCachesMu.Lock()
	defer RolesCachesMu.Unlock()

	if len(RolesCached.Roles) != 0 && RolesCached.ExpiresAt.After(time.Now()) {
		for _, role := range RolesCached.Roles {
			roles[role.Code] = role
		}

		return roles, nil
	}

	sqlStr := `select id, code, name from roles`
	rows, err := repo.Db.QueryContext(ctx, sqlStr)
	if err != nil {
		return roles, err
	}

	defer rows.Close()
	for rows.Next() {
		var role models.Role
		err = rows.Scan(&role.Id, &role.Code, &role.Name)
		if err != nil {
			return roles, err
		}

		roles[role.Code] = role
	}

	RolesCached.Roles = make([]models.Role, 0, len(roles))
	for _, role := range roles {
		RolesCached.Roles = append(RolesCached.Roles, role)
	}
	RolesCached.ExpiresAt = time.Now().Add(10 * time.Minute)

	return roles, nil
}

func (repo *DbRepository) GetAllUsers(reqCtx context.Context) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(reqCtx, time.Second*2)
	defer cancel()
	var users []models.User

	sqlStr := `select u.id, u.name, u.password, u.created_at, u.role_id, u.blocked,
       			r.id, r.code, r.name, r.created_at
			from users u
			left join roles r on u.role_id = r.id`

	rows, err := repo.Db.QueryContext(ctx, sqlStr)
	defer rows.Close()
	if err != nil {
		return users, err
	}

	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.Id, &user.Name, &user.Password, &user.CreatedAt, &user.RoleID, &user.Blocked,
			&user.Role.Id, &user.Role.Code, &user.Role.Name, &user.CreatedAt)

		if err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (repo *DbRepository) DeleteUser(reqCtx context.Context, userId int64) error {
	ctx, cancel := context.WithTimeout(reqCtx, time.Second*2)
	defer cancel()

	_, err := repo.Db.ExecContext(ctx, `delete from users where id = $1`, userId)
	return err
}

func (repo *DbRepository) ChangeBlockUser(reqCtx context.Context, userId int64) error {
	ctx, cancel := context.WithTimeout(reqCtx, time.Second*2)
	defer cancel()

	_, err := repo.Db.ExecContext(ctx, `update users set blocked = !blocked where id = $1`, userId)
	return err
}
