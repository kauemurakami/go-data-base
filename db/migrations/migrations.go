package migrations

import (
	"godb/packages/users"

	"gorm.io/gorm"
)

func init() {

}

func Migrations(db *gorm.DB) {
	db.AutoMigrate(
		users.User{},
	)
}
