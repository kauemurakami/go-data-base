package categories

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupCategoriesRoutes(router *mux.Router) {

	router.HandleFunc("/categories", CreateCategory).Methods(http.MethodPost)
	router.HandleFunc("/categories", GetCategories).Methods(http.MethodGet)

}
