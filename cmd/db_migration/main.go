package main

import (
	"log"

	"github.com/real-web-world/go-api/bootstrap"
	"github.com/real-web-world/go-api/global"
	"github.com/real-web-world/go-api/models"
)

func main() {
	bootstrap.InitMigration()
	db := global.DB
	if err := db.Set("gorm:table_options",
		"ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci").
		AutoMigrate(&models.User{},
			&models.City{},
			&models.File{},
			&models.Category{},
			&models.Tag{},
			&models.LoginHistory{},
			&models.Article{},
			&models.ArticleProfilePicture{},
			&models.ArticleTag{},
		); err != nil {
		log.Fatalln(err)
	} else {
		log.Println("migration success")
	}
}
