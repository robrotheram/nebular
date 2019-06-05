package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/yaml.v2"
)

type GalaxyInfo struct {
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
}

type GalaxyMetaComplex struct {
	GalaxyInfo   GalaxyInfo `yaml:"galaxy_info"`
	Dependencies []struct {
		Name string `yaml:"name"`
		Src  string `yaml:"src"`
		Scm  string `yaml:"scm"`
	} `yaml:"dependencies"`
}

type GalaxyMetaSimple struct {
	GalaxyInfo   GalaxyInfo `yaml:"galaxy_info"`
	Dependencies []string   `yaml:"dependencies"`
}

type GalaxyRole struct {
	ID        string      `json:"ID"`
	Namespace string      `json:"Namespace"`
	MetaType  string      `json:"MetaType"`
	Meta      interface{} `json:"Meta"`
	Rated     int         `json:"Rated"`
	Readme    string      `json:"Readme"`
	Repo      string      `json:"Repo"`
	Server    string      `json:"Server"`

	Username    string `json:"username"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func createData(namespace string, repo string) GalaxyRole {
	return GalaxyRole{
		Server:    "https://github.com",
		Namespace: namespace,
		Repo:      repo,
	}
}

//CreateRole will clone a repo and parse it for meta and readme info
// func CreateRoleStruct(server string, namespace string, repo string) *GalaxyRole {
// 	role := GalaxyRole{server, namespace, repo, GalaxyMeta{}, "", false, 0, []UserVotes{}}
// 	role.cloneRepo()
// 	role.getMeta()
// 	role.getReadme()
// 	return &role
// }

func (role *GalaxyRole) cloneRepo() {
	// Clone the given repository to the given directory
	path := fmt.Sprintf("%s/%s", configuration.GitTmpDir, role.Repo)
	url := fmt.Sprintf("%s/%s/%s", role.Server, role.Namespace, role.Repo)

	os.RemoveAll(path)

	Info("git clone %s %s", url, path)
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: configuration.GitUser,
			Password: configuration.GitPassword,
		},
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
		return
	}
	// ... retrieving the branch being pointed by HEAD
	//ref, err := r.Head()
	if err != nil {
		fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
		return
	}
	// ... retrieving the commit object
	//commit, err := r.CommitObject(ref.Hash())

	if err != nil {
		fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
		return
	}

}

func (role *GalaxyRole) getMeta() *GalaxyRole {
	path := fmt.Sprintf("%s/%s/meta/main.yml", configuration.GitTmpDir, role.Repo)
	fmt.Println("PATH: " + path)
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	meta := GalaxyMetaComplex{}
	err = yaml.Unmarshal(yamlFile, &meta)
	if err != nil {
		log.Println("Unmarshal to Complex failed, trying simple")
		meta := GalaxyMetaSimple{}
		err = yaml.Unmarshal(yamlFile, &meta)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
		log.Println("Unmarshal to Simple completed")
		role.MetaType = "SIMPLE"
		role.Meta = meta
	} else {
		role.MetaType = "COMPLEX"
		role.Meta = meta
	}
	fmt.Println(role)
	return role
}

func (role *GalaxyRole) getReadme() *GalaxyRole {
	path := fmt.Sprintf("%s/%s/README.md", configuration.GitTmpDir, role.Repo)
	b, err := ioutil.ReadFile(path) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	role.Readme = string(b)
	return role
}
