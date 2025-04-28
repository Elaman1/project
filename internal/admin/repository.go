package admin

import (
	"context"
	"database/sql"
	"myproject/internal/models"
	"time"
)

type Repository interface {
	LogAction(reqCtx context.Context, adminId int64, action string) error
	GetLogs(reqCtx context.Context) ([]models.Log, error)
}

type DbRepository struct {
	Db *sql.DB
}

func (repo *DbRepository) LogAction(reqCtx context.Context, adminId int64, action string) error {
	ctx, cancel := context.WithTimeout(reqCtx, time.Second*2)
	defer cancel()

	sqlStr := "insert into logs (user_id, action) values ($1, $2)"
	execErr := repo.Db.QueryRowContext(ctx, sqlStr, adminId, action)
	if execErr != nil {
		return execErr.Err()
	}

	return nil
}

func (repo *DbRepository) GetLogs(reqCtx context.Context) ([]models.Log, error) {
	ctx, cancel := context.WithTimeout(reqCtx, time.Second*2)
	defer cancel()
	rows, err := repo.Db.QueryContext(ctx, `select * from logs`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var logs []models.Log
	for rows.Next() {
		var log models.Log
		err = rows.Scan(&log.Id, &log.UserId, &log.Action, &log.CreatedAt)
		if err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}

	return logs, nil
}
