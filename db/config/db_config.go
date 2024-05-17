package db_config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func StartConfig() *gorm.DB {
	// db, err = gorm.Open( "postgres", "host=127.0.0.1 port=5432 user=postgres dbname=go-db sslmode=disable password=root")
	// dsn := "host=DB_HOST user=postgres password=root dbname=go-db port=5432 sslmode=disable TimeZone=America/Sao_Paulo"
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSL"),
		os.Getenv("DB_TZ"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(err.Error())
	}
	DB = db
	return db
}
