package middlewares

import (
	"errors"
	"myproject/config"
	"myproject/internal/auth"
	customerrors "myproject/internal/errors"
	"myproject/internal/functions"
	"net/http"
)

type RoleMiddleware struct {
	RoleName string
}

func (role *RoleMiddleware) Handle(next functions.CustomHttpHandler) functions.CustomHttpHandler {
	const op = "roleMiddleware.Handle"

	return func(w http.ResponseWriter, r *http.Request, ctxApp config.CtxApp) {
		user, ok := auth.GetUserFromContext(r.Context())
		if !ok {
			err := errors.New("пользователь не найден")
			customerrors.HandleJsonErrors(w, err, http.StatusUnauthorized, op)
			return
		}

		if user.Role.Code != role.RoleName {
			err := errors.New("нет доступа")
			customerrors.HandleJsonErrors(w, err, http.StatusUnauthorized, op)
			return
		}

		next(w, r, ctxApp)
	}
}

func (role *RoleMiddleware) Err() error {
	return nil
}
