package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Login    string `json:"login" gorm:"unique;not null;varchar(50)"`
	Password []byte `json:"-" gorm:"not null;varchar(100)"`
	Role     string `json:"role" gorm:"varchar(50)"`
}
