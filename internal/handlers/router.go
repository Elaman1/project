package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	custom_errors "myproject/internal/errors"
	"myproject/internal/functions"
	"myproject/internal/middlewares"
	"net/http"
	"strings"
	"time"
)

type Route struct {
	Address     string                       // Алрес роута
	Method      string                       // Метод роута
	Handler     functions.CustomHttpHandler  // Это чтобы сохранял метод функции
	Middlewares []middlewares.BaseMiddleware // Список мидлваров для роута
}

func (r *Route) GetMiddlewares() []middlewares.BaseMiddleware {
	r.Middlewares = append(
		r.Middlewares,
		&middlewares.RequestMethod{Method: r.Method},
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
	fmt.Println(req.Method, req.URL.Path)
outer:
	for _, route := range r.routes {
		fmt.Printf("%s %s\n", route.Method, route.Address)

		middlewaresList := route.GetMiddlewares()
		finalHandler := route.Handler
		// Проверка мидлваров
		for i := len(middlewaresList) - 1; i >= 0; i-- {
			mw := middlewaresList[i]
			finalHandler = mw.Handle(finalHandler)

			fmt.Println("middleware: ", mw)
			if err := mw.Err(); err != nil {
				var mwErr *custom_errors.MiddlewareNotAllowedError
				if errors.As(err, &mwErr) {
					continue outer
				}

				// Любая другая ошибка
				http.Error(w, "internal middleware error", http.StatusInternalServerError)
				return
			}
		}

		fmt.Println("route: ", route)
		// Установка ключей
		params, ok := matchPattern(route.Address, req.URL.Path)
		if ok {
			// Передаём параметры через context
			ctx := context.WithValue(req.Context(), "params", params)

			fmt.Println("params: ", params)
			fmt.Println("before end: ")
			route.Handler(w, req.WithContext(ctx), r.db)
			return
		}
	}
	http.NotFound(w, req)
}

func InitRoutes(db *sql.DB) *http.Server {
	newRoutes := NewRoutes(db)

	newRoutes.Handle(Route{
		Address: "/ping",
		Method:  http.MethodGet,
		Handler: PingFunc,
		Middlewares: []middlewares.BaseMiddleware{
			&middlewares.Auth{},
		},
	})

	newRoutes.Handle(Route{
		Address: "/auth/register",
		Method:  http.MethodPost,
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
