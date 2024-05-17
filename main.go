package main

import (
	"fmt"
	"godb/core/middlewares"
	"godb/core/routes"
	db_config "godb/db/config"
	"godb/db/migrations"
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
	//middleware set content type
	router.Use(middlewares.JSONContentTypeMiddleware)
	//Setup app all routes
	routes.SetupAppRoutes(router)

	fmt.Println("Listenning port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))

}
