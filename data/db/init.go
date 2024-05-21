package db

import (
	"blog/data/models"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var dns = "root:password@tcp(db:3306)/blog"
	fmt.Println(dns)
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

func GetUser(id string)(models.User,error){
	var u models.User
	err:=DB.Find(&u,id).Error
	return u,err
}
func DeleteUser(id string)(error){
	var u models.User
	err:=DB.Delete(&u,id).Error
	return err
}
func PostUser(username,password string)(models.User,error){
	var u models.User
	u.Username = username
	u.Password = password
	err:=DB.Create(&u).Error
	return u,err
}
func PutUser(id,username,password string)(error){
	u,err := GetUser(id)
	if err != nil{
		return err
	}
	u.Username = username
	u.Password = password
	err=DB.Save(&u).Error
	return err
}