package main

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var kv *KeyValueDB
var handler http.Handler

func GetRoles(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	roles := kv.SearchAll()
	json.NewEncoder(w).Encode(roles)
}

func SearchRoles(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	params := mux.Vars(req)
	roles := kv.SearchTerm(params["term"])
	json.NewEncoder(w).Encode(roles)
}

func GetRoleById(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	role := kv.Get(params["id"])
	json.NewEncoder(w).Encode(role)
}

func UpdateRole(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	role := kv.Get(params["id"])
	role.cloneRepo()
	role.getMeta()
	role.getReadme()
	kv.Save(role)
	json.NewEncoder(w).Encode(role)
}

func CreateRole(w http.ResponseWriter, req *http.Request) {
	var role GalaxyRole
	_ = json.NewDecoder(req.Body).Decode(&role)

	role.cloneRepo()
	role.getMeta()
	role.getReadme()

	kv.Save(role)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(role)

}

func DeleteRole(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]
	role := kv.Get(id)
	kv.Delete(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(role)
}

func GetUser(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Usr)
}

func GetSettings(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(settings{
		GitServer:    configuration.DefaultGitServer,
		GitNamespace: configuration.DefaultGitNamespace,
	})
}

func main() {
	configFromFile()
	configFromEnv()

	url := configuration.Hostname + ":" + configuration.Port
	fmt.Printf("magic is happening on %s \n", url)

	kv = NewKVDB(configuration.DbPath)

	router := mux.NewRouter()
	router.HandleFunc("/roles/update/{id}", UpdateRole).Methods("GET")
	router.HandleFunc("/roles/{id}", GetRoleById).Methods("GET")
	router.HandleFunc("/roles/{id}", DeleteRole).Methods("DELETE")
	router.HandleFunc("/roles", GetRoles).Methods("GET")
	router.HandleFunc("/roles", CreateRole).Methods("POST")
	router.HandleFunc("/search/{term}", SearchRoles).Methods("GET")
	router.HandleFunc("/search", GetRoles).Methods("GET")
	router.HandleFunc("/settings", GetSettings).Methods("GET")

	router.HandleFunc("/user", GetUser).Methods("GET")

	createApi(router)

	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./public/"))))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT"},
	})
	handler = c.Handler(router)

	log.Fatal(http.ListenAndServe(url, handler))
}
