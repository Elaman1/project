package config

import (
	"context"
	"log"
	"os"
	"time"
)

var StartTime time.Time

var GlobalCtx context.Context
var CancelGlobalCtx context.CancelFunc

func InitConfig() {
	StartTime = time.Now()
	log.Printf("Конфигурация инициализирована. Время запуска: %v", StartTime)
}

func InitGlobalConfig() {
	GlobalCtx, CancelGlobalCtx = context.WithCancel(context.Background())
	log.Println("Глобальный контекст инициализирован")
}

type Config struct {
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
}

func Load() *Config {
	return &Config{
		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASS"),
		DBName: os.Getenv("DB_NAME"),
	}
}
