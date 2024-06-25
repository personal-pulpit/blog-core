package database

import (
	"blog/config"
	"blog/internal/model"
	"blog/pkg/logging"
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mysqlInstance *gorm.DB
	mysqlMutex    = &sync.Mutex{}
)

func GetMysqlDB(cfg config.MysqlConfig) *gorm.DB {
	mysqlMutex.Lock()
	defer mysqlMutex.Unlock()
	if mysqlInstance == nil {
		var dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=%t",
			cfg.Username,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.DBname,
			cfg.ParseTime)
		db, err := gorm.Open(mysql.Open(dns))
		if err != nil {
			logging.MyLogger.Fatal(logging.General, logging.Startup, err.Error(), nil)
		}

		mysqlInstance = db
	}
	Migration(model.User{})
	Migration(model.Article{})
	Migration(model.Auth{})
	return mysqlInstance
}
