package database

import (
	"blog/config"
	"blog/internal/model"
	"blog/pkg/logging"
	"fmt"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	postgresInstance *gorm.DB
	mu      = &sync.Mutex{}
)

func GetPostgresqlDB(cfg config.Postgres) *gorm.DB {
	mu.Lock()
	defer mu.Unlock()
	if postgresInstance == nil {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", 
		cfg.Host, 
		cfg.Username, 
		cfg.Password, 
		cfg.DBName, 
		cfg.Port)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			logging.MyLogger.Fatal(logging.General, logging.Startup, err.Error(), nil)
		}

		postgresInstance = db
	}
	Migration(model.User{})
	Migration(model.Article{})
	Migration(model.Auth{})
	return postgresInstance
}
