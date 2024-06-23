package model

type Auth struct {
	ID             ID `gorm:"type:string;NOT NULL"`
	HashedPassword string `gorm:"size:300;NOT NULL"`
}
func NewAuth(ID ID,hashedPassword string)*Auth{
	return &Auth{
		ID: ID,
		HashedPassword: hashedPassword,
	}
}