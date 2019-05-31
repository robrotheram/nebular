package main

import (
	"fmt"


	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error

//GalaxyMeta contains all the data from meta/main.yml of a galexy role info
type GalaxyMeta struct {
	Dependencies []string `yaml:"dependencies"`
	GalaxyInfo   struct {
		Author            string  `yaml:"author"`
		Description       string  `yaml:"description"`
		Company           string  `yaml:"company"`
		License           string  `yaml:"license"`
		MinAnsibleVersion float64 `yaml:"min_ansible_version"`
		Platforms         []struct {
			Name     string   `yaml:"name"`
			Versions []string `yaml:"versions"`
		} `yaml:"platforms"`
		GalaxyTags []string `yaml:"galaxy_tags"`
	} `yaml:"galaxy_info"`
}

type UserVotes struct {
	username string
	vote     string
}

type GalaxyRole struct {
	gorm.Model
	Server    string
	Namespace string
	Repo      string
	Meta      GalaxyMeta
	Readme    string
	Supported bool
	Rated     int
	userVotes []UserVotes
}

func InitMigrate() {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to db")
	}
	defer db.Close()

	db.AutoMigrate(&GalaxyRole{})
}

// func GetRoles(w http.ResponseWriter, r *http.Request) {
// 	db, err := gorm.Open("sqlite3", "test.db")
// 	if err != nil {
// 		panic("failed to connect database")
// 	}
// 	defer db.Close()

// 	var roles []GalaxyRole
// 	db.Find(&roles)

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)

// 	json.NewEncoder(w).Encode(roles)
// }

// func CreateRole(w http.ResponseWriter, req *http.Request) {
// 	var role GalaxyRole
// 	_ = json.NewDecoder(req.Body).Decode(&role)

// 	db, err := gorm.Open("sqlite3", "test.db")
// 	if err != nil {
// 		panic("failed to connect database")
// 	}
// 	defer db.Close()
// 	var doUpdate = false
// 	if err := db.Where("repo = ?", role.Repo).Find(&role).Error; err != nil {
// 		doUpdate = true
// 	}

// 	role.cloneRepo()
// 	role.getMeta()
// 	role.getReadme()
// 	if doUpdate {
// 		fmt.Println(role)
// 		db.Update(&role)
// 		fmt.Fprintf(w, "Role Successfully Updated")
// 	} else {
// 		db.Create(&role)
// 		fmt.Fprintf(w, "New Role Successfully Created")
// 	}

// }

// func GetRoleById(w http.ResponseWriter, req *http.Request) {
// 	db, err := gorm.Open("sqlite3", "test.db")
// 	if err != nil {
// 		panic("failed to connect database")
// 	}
// 	defer db.Close()

// 	vars := mux.Vars(req)
// 	name := vars["id"]

// 	var role GalaxyRole
// 	db.Where("repo = ?", name).Find(&role)

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(role)
// }
