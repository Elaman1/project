package middlewares

import (
	"log"
	"myproject/config"
	"myproject/internal/functions"
	"net/http"
)

type Logging struct {
	Method  string
	Address string
}

func (logging *Logging) Handle(next functions.CustomHttpHandler) functions.CustomHttpHandler {
	return func(w http.ResponseWriter, r *http.Request, ctxApp config.CtxApp) {
		log.Printf("Method: %s, Address: %s \n", logging.Method, logging.Address)

		next(w, r, ctxApp)
	}
}

func (logging *Logging) Err() error {
	return nil
}
