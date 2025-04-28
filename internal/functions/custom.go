package functions

import (
	"myproject/config"
	"net/http"
)

type CustomHttpHandler func(http.ResponseWriter, *http.Request, config.CtxApp)

func Param(r *http.Request, key string) string {
	if params, ok := r.Context().Value("params").(map[string]string); ok {
		return params[key]
	}
	return ""
}
