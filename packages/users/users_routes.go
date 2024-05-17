package users

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupUserRoutes(router *mux.Router) {

	router.HandleFunc("/users", CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/users", GetUsers).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", GetUserById).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", DeleteUserById).Methods(http.MethodDelete)
	router.HandleFunc("/users", UpdateUserById).Methods(http.MethodPut)

}
