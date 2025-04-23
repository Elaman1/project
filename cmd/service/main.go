package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"myproject/config"
	internalconfig "myproject/internal/config"
	internaldb "myproject/internal/db"
	"myproject/internal/handlers"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	err := internalconfig.LoadEnv(".env")
	if err != nil {
		panic(err)
	}

	connectConfig := config.Load()

	db, err := internaldb.NewDB(
		connectConfig.DBHost,
		connectConfig.DBPort,
		connectConfig.DBUser,
		connectConfig.DBPass,
		connectConfig.DBName,
	)

	if err != nil {
		panic(err)
	}

	defer func(Conn *sql.DB) {
		err = Conn.Close()
		if err != nil {
			panic(err)
		}
	}(db.Conn)

	srv := handlers.InitRoutes(db.Conn)

	// Канал для завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Starting server on :8080")
		config.InitConfig()
		config.InitGlobalConfig()

		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Ошибка запуска сервера: %v", err)
			os.Exit(1)
		}

	}()

	// Ожидание завершения
	<-quit
	log.Println("Получен сигнал завершения, закрываем сервер...")

	// Отмена глобального контекст
	config.CancelGlobalCtx()

	// Контекст с тайм-аутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Корректное завершение сервера
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Ошибка завершения сервера: %v", err)
		os.Exit(1)
	}
	log.Println("Сервер завершил работу корректно")
}
