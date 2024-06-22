package mysql_repository

import "gorm.io/gorm"

func NewTx(DB *gorm.DB) *gorm.DB {
	return DB.Begin()
}
