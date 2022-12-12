package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ctrmgmt/models"
)

func GetContainers(w http.ResponseWriter, r *http.Request) {
	var ctrs []models.CtrMgt
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ctrs)
}
