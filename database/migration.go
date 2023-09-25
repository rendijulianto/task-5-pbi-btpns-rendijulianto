package database

import (
	"final-project-rakamin/models"
	"log"
)

func Migrate() {
	Instance.AutoMigrate(&models.User{}, &models.Photo{})
	log.Println("Database Migration Completed!")
}