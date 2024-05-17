package routes

import (
	"godb/packages/categories"
	"godb/packages/product"
	"godb/packages/sales"
	"godb/packages/users"

	"github.com/gorilla/mux"
)

func SetupAppRoutes(router *mux.Router) {
	users.SetupUserRoutes(router)
	product.SetupProductsRoutes(router)
	sales.SetupSalesRoutes(router)
	categories.SetupCategoriesRoutes(router)
}
