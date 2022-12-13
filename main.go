package main

import (
	"ctrmgmt/controllers"
	"fmt"
	"log"
	"net/http"

	"ctrmgmt/config"

	"github.com/gorilla/mux"
)

func RegisterProductRoutes(router *mux.Router) {
	router.HandleFunc("/version", controllers.GetVersion).Methods("GET")
	router.HandleFunc("/api/containers", controllers.GetContainers).Methods("GET")
	//router.HandleFunc("/api/products/{id}", controllers.GetProductById).Methods("GET")
	//router.HandleFunc("/api/products", controllers.CreateProduct).Methods("POST")
	//router.HandleFunc("/api/products/{id}", controllers.UpdateProduct).Methods("PUT")
	//router.HandleFunc("/api/products/{id}", controllers.DeleteProduct).Methods("DELETE")
}

func main() {
	// Load Configurations from config.json using Viper
	config.LoadAppConfig()
	// Initialize the router
	router := mux.NewRouter().StrictSlash(true)

	// Register Routes
	RegisterProductRoutes(router)

	// Start the server
	log.Println(fmt.Sprintf("Starting Server on port %s", config.AppConfig.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.AppConfig.Port), router))
}
