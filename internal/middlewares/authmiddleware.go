package middlewares

import (
	"database/sql"
	"fmt"
	"myproject/internal/functions"
	"net/http"
)

type Auth struct {
}

func (a *Auth) Handle(next functions.CustomHttpHandler) functions.CustomHttpHandler {
	return func(w http.ResponseWriter, r *http.Request, db *sql.DB) {
		fmt.Println("Авторизация")
		next(w, r, db)
	}
}

func (a *Auth) Err() error {
	return nil
}
