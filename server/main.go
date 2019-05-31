package main

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"
	"github.com/gorilla/mux"
)

// Config stores all setting about the application
type Config struct {
	GitUsername string
	GitPassword string
	TmpDIR      string
}

var config = Config{GitUsername: "", GitPassword: "", TmpDIR: "/tmp/repoCache"}

var roles = []GalaxyRole{}

func findRoleByRepo(repo string) GalaxyRole {
	var role GalaxyRole
	for _, erole := range roles {
		if erole.Repo == repo {
			role = erole
		}
	}
	return role
}

func GetRoles(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(roles)
}

func GetRoleById(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(findRoleByRepo(params["id"]))
}

func CreateRole(w http.ResponseWriter, req *http.Request) {
	var role GalaxyRole
	_ = json.NewDecoder(req.Body).Decode(&role)

	if findRoleByRepo(role.Repo).Repo == role.Repo {
		fmt.Fprintf(w, "Repo %s already exists", role.Repo)
	} else {
		role.cloneRepo()
		role.getMeta()
		role.getReadme()
		roles = append(roles, role)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(roles)
	}

}

func DeleteRole(w http.ResponseWriter, req *http.Request) {
	// params := mux.Vars(req)
	// var user User
	// _ = json.NewDecoder(req.Body).Decode(&user)
	// user.ID = params["id"]
	// users = append(users, user)

	// json.NewEncoder(w).Encode(users)
}

func main() {
	fmt.Println("magic is happening on port 8081")
	InitMigrate()

	router := mux.NewRouter()
	router.HandleFunc("/roles", GetRoles).Methods("GET")
	router.HandleFunc("/roles", CreateRole).Methods("POST")
	router.HandleFunc("/roles/{id}", GetRoleById).Methods("GET")
	router.HandleFunc("/roles/{id}", DeleteRole).Methods("DELETE")

	// router.HandleFunc("/roles/{id}/vote/up", VoteRoleUp).Methods("GET")
	// router.HandleFunc("/roles/{id}/vote/down", VoteRoleDown).Methods("GET")

	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./public/"))))

	log.Fatal(http.ListenAndServe(":8081", router))
}
