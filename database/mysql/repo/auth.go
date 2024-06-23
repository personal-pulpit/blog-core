package mysql_repository

import (
	"blog/internal/model"
	"blog/internal/repository"

	"gorm.io/gorm"
)

type authMysqlRepository struct {
	mysqlClient *gorm.DB
}
func NewAuthMysqlRepository(mysqlCLI *gorm.DB)repository.AuthMysqlRepository{
	return &authMysqlRepository{
		mysqlClient:mysqlCLI,
	}
}
func (a *authMysqlRepository) Create(authModel *model.Auth)(*model.Auth,error){
	tx := a.mysqlClient.Create(authModel)
	if tx.Error != nil{
		return nil,tx.Error
	}
	return authModel,nil
}
func (a *authMysqlRepository)GetUserAuth(ID model.ID) (*model.Auth, error){
	auth := &model.Auth{}
	tx := a.mysqlClient.Where("id= ?", ID).First(auth)
	if tx.Error != nil{
		return nil,tx.Error
	}
	return auth,nil
}
func (a *authMysqlRepository)ChangePassword(ID model.ID,hashedPassword string)error{
	authModel,err := a.GetUserAuth(ID)
	if err != nil{
		return err
	}
	authModel.HashedPassword =hashedPassword
	tx := a.mysqlClient.Save(authModel)
	if tx.Error != nil{
		return tx.Error
	}
	return nil
}
func (a *authMysqlRepository) DeleteByID(ID model.ID) error {
	authModel,err := a.GetUserAuth(ID)
	if err != nil{
		return err
	}
	tx := a.mysqlClient.Delete(&authModel)
	if tx.Error != nil{
		return tx.Error
	}
	return nil
}