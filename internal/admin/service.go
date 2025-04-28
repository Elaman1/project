package admin

import (
	"context"
	"myproject/internal/models"
)

type Service struct {
	Repo Repository
}

func (s *Service) LoggingAdminAction(ctx context.Context, adminId int64, action string) error {
	err := s.Repo.LogAction(ctx, adminId, action)
	return err
}

func (s *Service) GetLogs(ctx context.Context) ([]models.Log, error) {
	return s.Repo.GetLogs(ctx)
}
