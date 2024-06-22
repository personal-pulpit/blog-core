package model

type Auth struct {
	ID string `gorm:"type:string;NOT NULL"`
	HashedPassword string `gorm:"size:300;NOT NULL"`
}