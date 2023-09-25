package models

import (
	"final-project-rakamin/helpers"
	"time"
)

type User struct {
	ID       uint     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username string `json:"username" gorm:"unique"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Photos   []Photo `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}



func (user *User) HashPassword(password string) error {
	has, err := helpers.HashPassword(password)
	if err != nil {
		return err

	}
	user.Password = has
	return nil
}
func (user *User) CheckPassword(providedPassword string) error {
	result, err:= helpers.CheckPasswordHash(providedPassword, user.Password)
	if !result {
		return err
	}
	return nil
}
