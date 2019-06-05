package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type InstallRole struct {
	DownloadUrl string `json:"download_url"`
	Version     string `json:"version"`
}

type SearchResp struct {
	Count   int          `json:"count"`
	Results []GalaxyRole `json:"results"`
}

type InstallResp struct {
	Count   int           `json:"count"`
	Results []InstallRole `json:"results"`
}
type ApiVersion struct {
	Description    string `json:"description"`
	CurrentVersion string `json:"current_version"`
}

func GetApiVersion(w http.ResponseWriter, req *http.Request) {
	apiVersion := ApiVersion{
		Description:    "GALAXY REST API",
		CurrentVersion: "v1",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(apiVersion)
	return
}

func searchRoles(w http.ResponseWriter, req *http.Request) {
	queryValues := req.URL.Query()
	term := queryValues.Get("autocomplete")
	fmt.Println()
	roles := kv.SearchTerm(term)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := SearchResp{
		Count:   len(roles),
		Results: roles,
	}
	json.NewEncoder(w).Encode(response)
}

func createApi(router *mux.Router) {
	router.HandleFunc("/api/", GetApiVersion).Methods("GET")
	router.HandleFunc("/api/v1/search/roles/", searchRoles).Methods("GET")
}
