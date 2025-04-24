package handlers

import (
	"context"
	"database/sql"
	"myproject/config"
	"myproject/internal/auth"
	"myproject/internal/functions"
	"myproject/internal/middlewares"
	"net/http"
	"strings"
	"time"
)

type Route struct {
	Address     string
	Method      string
	Handler     functions.CustomHttpHandler
	Middlewares []middlewares.BaseMiddleware
}

func (r *Route) GetMiddlewares() []middlewares.BaseMiddleware {
	r.Middlewares = append(
		r.Middlewares,
		&middlewares.Logging{
			Method:  r.Method,
			Address: r.Address,
		},
	)
	return r.Middlewares
}

type Router struct {
	routes []Route
	db     *sql.DB
}

func NewRoutes(db *sql.DB) *Router {
	return &Router{
		db: db,
	}
}

func (r *Router) Handle(route Route) {
	r.routes = append(r.routes, route)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		// Установка ключей
		params, ok := matchPattern(route.Address, req.URL.Path)
		if !ok {
			continue
		}

		if req.Method != route.Method {
			continue
		}

		middlewaresList := route.GetMiddlewares()
		finalHandler := route.Handler
		// Проверка мидлваров
		for i := len(middlewaresList) - 1; i >= 0; i-- {
			mw := middlewaresList[i]
			finalHandler = mw.Handle(finalHandler)

			if err := mw.Err(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				http.Error(w, "internal middleware error", http.StatusInternalServerError)
				return
			}
		}

		// Передаём параметры через context
		ctx := context.WithValue(req.Context(), "params", params)
		finalHandler(w, req.WithContext(ctx), config.CtxApp{
			Db: r.db,
		})
		return
	}
	http.NotFound(w, req)
}

func InitRoutes(db *sql.DB) *http.Server {
	newRoutes := NewRoutes(db)

	newRoutes.Handle(Route{
		Address: "/ping",
		Method:  http.MethodGet,
		Handler: PingFunc,
	})

	newRoutes.Handle(Route{
		Address: "/auth/register",
		Method:  http.MethodPost,
		Handler: auth.RegisterHandler,
	})

	newRoutes.Handle(Route{
		Address: "/auth/login",
		Method:  http.MethodPost,
		Handler: auth.LoginHandler,
	})

	newRoutes.Handle(Route{
		Address: "/auth/me",
		Method:  http.MethodGet,
		Handler: auth.MeHandler,
		Middlewares: []middlewares.BaseMiddleware{
			&middlewares.Auth{},
		},
	})

	return &http.Server{
		Addr:         ":8080",
		Handler:      newRoutes,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func matchPattern(pattern, path string) (map[string]string, bool) {
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(patternParts) != len(pathParts) {
		return nil, false
	}

	params := make(map[string]string)

	for i := range patternParts {
		pp := patternParts[i]
		pv := pathParts[i]

		if strings.HasPrefix(pp, "{") && strings.HasSuffix(pp, "}") {
			key := pp[1 : len(pp)-1]
			params[key] = pv
		} else if pp != pv {
			return nil, false
		}
	}

	return params, true
}

func Param(r *http.Request, key string) string {
	if params, ok := r.Context().Value("params").(map[string]string); ok {
		return params[key]
	}
	return ""
}
