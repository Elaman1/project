package middlewares

import (
	"database/sql"
	"log"
	"myproject/internal/functions"
	"net/http"
)

type Logging struct {
	Method  string
	Address string
}

func (logging *Logging) Handle(next functions.CustomHttpHandler) functions.CustomHttpHandler {
	return func(w http.ResponseWriter, r *http.Request, db *sql.DB) {
		log.Printf("Method: %s, Address: %s \n", logging.Method, logging.Address)

		next(w, r, db)
	}
}

func (logging *Logging) Err() error {
	return nil
}
