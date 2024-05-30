package utils

import "golang.org/x/crypto/bcrypt"


func HashPassword(password string)(string,error){
	hashedPassword,err:=bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)	
	if err != nil{
		return "",err
	}
	return string(hashedPassword),nil
}
func CheckPassword(password,hashedpassword string)(error){
	err := bcrypt.CompareHashAndPassword([]byte(hashedpassword),[]byte(password))
	return err
}