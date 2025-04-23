package functions

import (
	"myproject/config"
	"net/http"
)

type CustomHttpHandler func(http.ResponseWriter, *http.Request, config.CtxApp)
