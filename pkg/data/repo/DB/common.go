package db

import "gorm.io/gorm"

func NewTx(db *gorm.DB)*gorm.DB{
	return db.Begin() 
}