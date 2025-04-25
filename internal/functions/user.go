package functions

import (
	"context"
	"myproject/config"
	"myproject/internal/models"
)

func GetUserFromContext(ctx context.Context) (*models.User, bool) {
	userVal := ctx.Value(config.CtxUserKey)
	user, ok := userVal.(*models.User)
	return user, ok
}
