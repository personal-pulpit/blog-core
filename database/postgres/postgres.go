package database

import (
	"blog/config"
	"blog/internal/model"
	"fmt"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	postgresInstance *gorm.DB
	mu      = &sync.Mutex{}
)

func GetPostgresqlDB(cfg *config.Postgres) (*gorm.DB,error) {
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
			return nil,err
		}

		postgresInstance = db
	}
	
	Migration(model.User{})
	Migration(model.Article{})
	Migration(model.Auth{})

	return postgresInstance,nil
}
