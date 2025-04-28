package admin

import (
	"context"
	"encoding/json"
	"errors"
	"myproject/config"
	"myproject/internal/auth"
	customerrors "myproject/internal/errors"
	"myproject/internal/functions"
	"myproject/internal/models"
	"net/http"
	"strconv"
)

func PanelHandler(w http.ResponseWriter, r *http.Request, ctxApp config.CtxApp) {
	const op = "panelHandler"

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{"message": "admin panel"})
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusInternalServerError, op)
		return
	}
}

func UsersHandler(w http.ResponseWriter, r *http.Request, ctxApp config.CtxApp) {
	const op = "usersHandler"

	w.Header().Set("Content-Type", "application/json")
	authService := auth.Service{
		Repo: &auth.DbRepository{
			Db: ctxApp.Db,
		},
	}

	users, err := authService.GetAllUsers(r.Context())
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusInternalServerError, op)
		return
	}

	// Логирование действие админа
	err = logAction(r.Context(), ctxApp, models.GetUsersAction)
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusInternalServerError, op)
		return
	}

	err = json.NewEncoder(w).Encode(users)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request, ctxApp config.CtxApp) {
	const op = "deleteUserHandler"
	authService := auth.Service{
		Repo: &auth.DbRepository{
			Db: ctxApp.Db,
		},
	}

	userId, err := getId(r)
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusInternalServerError, op)
		return
	}

	err = authService.DeleteUser(r.Context(), userId)
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusInternalServerError, op)
		return
	}

	err = logAction(r.Context(), ctxApp, models.DeleteAction)
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusInternalServerError, op)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(Response{
		Message: "Успешно удалено",
	})
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusInternalServerError, op)
		return
	}
}

func BlockUserHandler(w http.ResponseWriter, r *http.Request, ctxApp config.CtxApp) {
	const op = "blockUserHandler"
	w.Header().Set("Content-Type", "application/json")

	authService := auth.Service{
		Repo: &auth.DbRepository{
			Db: ctxApp.Db,
		},
	}

	userId, err := getId(r)
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusInternalServerError, op)
		return
	}

	err = authService.ChangeBlockUser(r.Context(), userId)
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusInternalServerError, op)
		return
	}

	err = logAction(r.Context(), ctxApp, models.BlockAction)
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusInternalServerError, op)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(Response{
		Message: "Успешно изменено",
	})

	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusInternalServerError, op)
		return
	}
}

func LogsHandler(w http.ResponseWriter, r *http.Request, ctxApp config.CtxApp) {
	const op = "logsHandler"
	w.Header().Set("Content-Type", "application/json")
	adminService := Service{
		Repo: &DbRepository{
			Db: ctxApp.Db,
		},
	}

	res, err := adminService.GetLogs(r.Context())
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusInternalServerError, op)
		return
	}

	err = logAction(r.Context(), ctxApp, models.ShowLogAction)
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusInternalServerError, op)
		return
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		customerrors.HandleJsonErrors(w, err, http.StatusInternalServerError, op)
		return
	}
}

func logAction(reqCtx context.Context, ctxApp config.CtxApp, action string) error {
	adminService := Service{
		Repo: &DbRepository{
			Db: ctxApp.Db,
		},
	}

	admin, ok := functions.GetUserFromContext(reqCtx)
	if !ok {
		return errors.New("no admin user found")
	}

	err := adminService.LoggingAdminAction(reqCtx, admin.Id, action)
	if err != nil {
		return err
	}

	return nil
}

func getId(r *http.Request) (int64, error) {
	userIdVal := functions.Param(r, "id")
	if userIdVal == "" {
		return 0, errors.New("invalid user id")
	}

	userId, err := strconv.Atoi(userIdVal)
	if err != nil {
		return 0, err
	}

	return int64(userId), nil
}
