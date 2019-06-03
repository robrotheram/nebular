package main

import (
	"encoding/json"
	"fmt"
)

type GalaxyRole struct {
	Namespace string `json:"Namespace"`
	Rated     int    `json:"Rated"`
	Readme    string `json:"Readme"`
	Repo      string `json:"Repo"`
	Server    string `json:"Server"`
}

func main() {
	s := string(`{"Namespace":"test","Rated":0,"Readme":"","Repo":"docker","Server":"https://github.com"}`)
	data := GalaxyRole{}
	json.Unmarshal([]byte(s), &data)
	fmt.Printf("Operation: %s", data.Server)
}
