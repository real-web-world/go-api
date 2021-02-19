package main

import (
	"log"

	"github.com/real-web-world/go-web-api/bootstrap"
	"github.com/real-web-world/go-web-api/global"
	"github.com/real-web-world/go-web-api/models"
)

func main() {
	bootstrap.InitMigration()
	db := global.DB
	err := db.Set("gorm:table_options",
		"ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci").
		AutoMigrate(&models.User{},
			&models.City{},
			&models.File{},
			&models.Category{},
			&models.Tag{},
			&models.LoginHistory{})
	log.Println(err)
}
