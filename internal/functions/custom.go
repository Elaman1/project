package functions

import (
	"database/sql"
	"net/http"
)

type CustomHttpHandler func(http.ResponseWriter, *http.Request, *sql.DB)
