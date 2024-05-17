package product

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupProductsRoutes(router *mux.Router) {

	router.HandleFunc("/products", CreateProduct).Methods(http.MethodPost)
	router.HandleFunc("/products", GetProducts).Methods(http.MethodGet)
	router.HandleFunc("/products/{id}", DeleteProductById).Methods(http.MethodDelete)

}
