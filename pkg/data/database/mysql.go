package database

import (
	"blog/config"
	"blog/pkg/data/models"
	"blog/pkg/logging"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=%t",
		config.Cfg.Mysql.Username,
		config.Cfg.Mysql.Password,
		config.Cfg.Mysql.Host,
		config.Cfg.Mysql.Port,
		config.Cfg.Mysql.DBname,
		config.Cfg.Mysql.ParseTime,
	)
	db, err := gorm.Open(mysql.Open(dns))
	if err != nil {
		logging.MyLogger.Fatal(logging.General, logging.Startup, err.Error(), nil)
	}
	DB = db
	Migration(models.User{})
	Migration(models.Article{})
}
