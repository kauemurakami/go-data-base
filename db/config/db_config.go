package db_config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func StartConfig() *gorm.DB {
	// db, err = gorm.Open( "postgres", "host=127.0.0.1 port=5432 user=postgres dbname=go-db sslmode=disable password=root")
	dsn := "host=localhost user=postgres password=root dbname=go-db port=5432 sslmode=disable TimeZone=America/Sao_Paulo"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(err.Error())
	}
	DB = db
	return db
}
