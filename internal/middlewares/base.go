package middlewares

import (
	"myproject/internal/functions"
)

type BaseMiddleware interface {
	Handle(handlerFunc functions.CustomHttpHandler) functions.CustomHttpHandler
	Err() error
}
