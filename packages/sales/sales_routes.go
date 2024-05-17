package sales

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupSalesRoutes(router *mux.Router) {

	router.HandleFunc("/sales", CreateSale).Methods(http.MethodPost)
	router.HandleFunc("/sales", GetSales).Methods(http.MethodGet)
	router.HandleFunc("/sales/{id}", DeleteSaleById).Methods(http.MethodDelete)

}
