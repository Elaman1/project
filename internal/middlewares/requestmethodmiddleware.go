package middlewares

import (
	"database/sql"
	"myproject/internal/errors"
	"myproject/internal/functions"
	"net/http"
)

type RequestMethod struct {
	Method    string
	ErrorType error
}

func (a *RequestMethod) Handle(next functions.CustomHttpHandler) functions.CustomHttpHandler {
	return func(w http.ResponseWriter, r *http.Request, db *sql.DB) {
		if r.Method != a.Method {
			a.ErrorType = errors.MiddlewareNotAllowedError{Code: http.StatusMethodNotAllowed, Message: "Not allowed method"}
			return
		}

		next(w, r, db)
	}
}

func (a *RequestMethod) Err() error {
	return a.ErrorType
}
