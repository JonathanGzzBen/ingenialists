package models

import (
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func DB() (*gorm.DB, error) {
	if db != nil {
		return db, nil
	}
	db, err := gorm.Open(sqlite.Open("test.db"))
	if err != nil {
		return nil, errors.New("could not connect to database")
	}
	return db, nil
}
