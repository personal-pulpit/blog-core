package repository
import "blog/internal/model"
type AuthMysqlRepository interface {
	Create(authModel *model.Auth) (*model.Auth, error)
	GetUserAuth(ID model.ID) (*model.Auth, error)
	ChangePassword(ID model.ID,hashedPassword string)error
	DeleteByID(ID model.ID) error
}
type AuthRedisRepository interface {
}
