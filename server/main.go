package main

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Config stores all setting about the application
type Config struct {
	GitUsername string
	GitPassword string
	TmpDIR      string
}

var config = Config{GitUsername: "", GitPassword: "", TmpDIR: "/tmp/repoCache"}
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

func main() {
	fmt.Println("magic is happening on port 8081")
	kv = NewKVDB("/tmp/nebular")

	// kv.Save(createData("test", "docker"))
	// kv.Save(createData("test", "jekins"))
	// kv.Save(createData("doom", "jekins"))

	// fmt.Println(kv.SearchTerm("docker"))
	// fmt.Println(kv.SearchTerm("jekins"))

	//kv.SearchAll()

	router := mux.NewRouter()
	router.HandleFunc("/roles", GetRoles).Methods("GET")
	router.HandleFunc("/roles", CreateRole).Methods("POST")
	router.HandleFunc("/roles/{id}", GetRoleById).Methods("GET")
	router.HandleFunc("/search/{term}", SearchRoles).Methods("GET")
	router.HandleFunc("/search", GetRoles).Methods("GET")
	router.HandleFunc("/roles/{id}", DeleteRole).Methods("DELETE")
	router.HandleFunc("/user", GetUser).Methods("GET")

	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./public/"))))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT"},
	})
	handler = c.Handler(router)

	log.Fatal(http.ListenAndServe(":8081", handler))
}
