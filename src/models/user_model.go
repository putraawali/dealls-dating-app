package models

import (
	"dealls-dating-app/src/dtos"
	"dealls-dating-app/src/pkg/helpers"
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserID     int64     `json:"user_id" gorm:"primaryKey;column:user_id"`
	Email      string    `json:"email" gorm:"unique;not null;type:varchar(255)"`
	IsVerified bool      `json:"is_verified" gorm:"is_verified"`
	Sex        string    `json:"sex" gorm:"column:sex;type:varchar(6)"`
	FirstName  string    `json:"first_name" gorm:"not null;column:first_name;type:varchar(255)"`
	LastName   string    `json:"last_name" gorm:"type:varchar(255)"`
	IsPremium  bool      `json:"is_premium" gorm:"column:is_premium;default:false"`
	Password   string    `json:"-" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	sex := map[string]bool{
		"male":   true,
		"female": true,
	}

	if !sex[u.Sex] {
		return errors.New("possible value for sex is only male or female")
	}

	u.Password = helpers.HashPassword(u.Password)

	return
}

func (u *User) RegisterToModel(data dtos.RegisterParam) {
	u.Email = data.Email
	u.Sex = data.Sex
	u.Password = data.Password
	u.FirstName = data.FirstName
	u.LastName = data.LastName
}
