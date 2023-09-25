package models

import (
	"time"
)

type Photo struct {
	ID      uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title     string     `gorm:"column:title" json:"title"`
	Caption  string     `gorm:"column:caption" json:"caption"`
	PhotoUrl string     `gorm:"column:photo_url" json:"photo_url"`
	UserID    uint       `gorm:"column:user_id" json:"user_id"`
	User      User       `gorm:"foreignKey:UserID"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	
}