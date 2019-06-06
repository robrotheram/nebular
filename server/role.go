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
//GalaxyInfo Internal Info Struct
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
//GalaxyMetaComplex Contains info with a complex dependacy format
type GalaxyMetaComplex struct {
	GalaxyInfo   GalaxyInfo `yaml:"galaxy_info"`
	Dependencies []struct {
		Name string `yaml:"name"`
		Src  string `yaml:"src"`
		Scm  string `yaml:"scm"`
	} `yaml:"dependencies"`
}
//GalaxyMetaSimple Contains info with a simple dependacy format
type GalaxyMetaSimple struct {
	GalaxyInfo   GalaxyInfo `yaml:"galaxy_info"`
	Dependencies []string   `yaml:"dependencies"`
}

//GalaxyRole contains the info about a role got from meta/main.yml
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

func (role *GalaxyRole) cloneRepo() error {
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
		return err
	}
	return nil
}

func (role *GalaxyRole) getMeta() error {
	path := fmt.Sprintf("%s/%s/meta/main.yml", configuration.GitTmpDir, role.Repo)
	fmt.Println("PATH: " + path)
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	meta := GalaxyMetaComplex{}
	err = yaml.Unmarshal(yamlFile, &meta)
	if err != nil {
		log.Println("Unmarshal to Complex failed, trying simple")
		meta := GalaxyMetaSimple{}
		err = yaml.Unmarshal(yamlFile, &meta)
		if err != nil {
			return err
		}
		log.Println("Unmarshal to Simple completed")
		role.MetaType = "SIMPLE"
		role.Meta = meta
	} else {
		role.MetaType = "COMPLEX"
		role.Meta = meta
	}
	return nil
}

func (role *GalaxyRole) getReadme() error {
	path := fmt.Sprintf("%s/%s/README.md", configuration.GitTmpDir, role.Repo)
	b, err := ioutil.ReadFile(path) // just pass the file name
	if err != nil {
		return err
	}
	role.Readme = string(b)
	return nil
}
