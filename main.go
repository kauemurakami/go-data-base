package main

import (
	"fmt"
	db_config "godb/db/config"
	"godb/db/migrations"
	"godb/packages/users"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	//start db
	db := db_config.StartConfig()
	//migrations
	migrations.Migrations(db)
	//router
	router := mux.NewRouter()
	router.HandleFunc("/users", users.CreateUser).Methods(http.MethodPost)

	fmt.Println("Listenning port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))

}
