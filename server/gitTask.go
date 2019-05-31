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
	path := fmt.Sprintf("%s/%s", config.TmpDIR, role.Repo)
	url := fmt.Sprintf("%s/%s/%s", role.Server, role.Namespace, role.Repo)

	os.RemoveAll(path)

	Info("git clone %s %s", url, path)
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: config.GitUsername,
			Password: config.GitPassword,
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
	path := fmt.Sprintf("%s/%s/meta/main.yml", config.TmpDIR, role.Repo)
	fmt.Println("PATH: " + path)
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &role.Meta)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return role
}

func (role *GalaxyRole) getReadme() *GalaxyRole {
	path := fmt.Sprintf("%s/%s/README.md", config.TmpDIR, role.Repo)
	b, err := ioutil.ReadFile(path) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	role.Readme = string(b)
	return role
}
