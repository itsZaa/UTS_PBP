package controller

import (
	"database/sql"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func connect() *sql.DB {
	db, err := sql.Open(
		"mysql",
		"root:@tcp(localhost:3306)/uts_pbp")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func connectGorm() (*gorm.DB, error) {
	dsn := "root:@tcp(localhost:3306)/uts_pbp?parseTime=true&loc=Asia%2FJakarta"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
