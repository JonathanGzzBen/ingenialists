package models

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open(sqlite.Open("test.db"))
	if err != nil {
		log.Panic("could not open database file")
	}
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Category{})
	db.AutoMigrate(&Article{})
	DB = db
}
