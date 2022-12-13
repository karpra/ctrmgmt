package controllers

import (
	"encoding/json"
	"net/http"

	"ctrmgmt/models"
)

func GetVersion(w http.ResponseWriter, r *http.Request) {
	version := "0.0.1"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(version)
}

func GetContainers(w http.ResponseWriter, r *http.Request) {
	var ctrs []models.CtrMgt
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ctrs)
}
