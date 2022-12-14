package main

import (
	"ctrmgmt/config"
	"ctrmgmt/controllers"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterProductRoutes(router *mux.Router) {
	router.HandleFunc("/version", controllers.GetVersion).Methods("GET")
	router.HandleFunc("/api/containers", controllers.GetContainers).Methods("GET")
	router.HandleFunc("/api/containers", controllers.CreateContainers).Methods("POST")
	router.HandleFunc("/api/containers/{name}", controllers.StopContainers).Methods("DELETE")
	router.HandleFunc("/api/containers/{name}/logs", controllers.GetContainerLogs).Methods("GET")
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
	log.Printf("Starting Server on port %s", config.AppConfig.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%+v", config.AppConfig.Port), router))
}
