package database

import (
	"blog/config"
	"blog/pkg/data/models"
	"fmt"
	"log"

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
		log.Fatalln(err)
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalln(err)
	}
	DB = db
}

