package migrations

import (
	"godb/packages/categories"
	"godb/packages/product"
	"godb/packages/sales"
	"godb/packages/users"

	"gorm.io/gorm"
)

func init() {

}

func Migrations(db *gorm.DB) {   
	db.AutoMigrate(
		users.User{},
		categories.Category{},
		product.Product{},
		sales.Sale{},
	)
}
