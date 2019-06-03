package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/document"
	"github.com/blevesearch/bleve/index/scorch"
	"github.com/mitchellh/mapstructure"
)

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

z
type GalaxyRole struct {
	ID		  string	 `json:"ID"`
	Namespace string     `json:"Namespace"`
	Meta      GalaxyMeta `json:"Meta"`
	Rated     int        `json:"Rated"`
	Readme    string     `json:"Readme"`
	Repo      string     `json:"Repo"`
	Server    string     `json:"Server"`
}









func createData(namespace string, repo string) GalaxyRole {
	return GalaxyRole{
		Server:    "https://github.com",
		Namespace: namespace,
		Repo:      repo,
	}
}

func getDocsFromSearchResults(results *bleve.SearchResult, index bleve.Index) []GalaxyRole {
	docs := make([]GalaxyRole, 0)
	for _, val := range results.Hits {
		id := val.ID
		doc, _ := index.Document(id)
		fields := map[string]interface{}{}
		
		for _, field := range doc.Fields {
			var newval interface{}
			switch field := field.(type) {
			case *document.TextField:
				newval = string(field.Value())
			case *document.NumericField:
				n, err := field.Number()
				if err == nil {
					newval = n
				}
			case *document.DateTimeField:
				d, err := field.DateTime()
				if err == nil {
					newval = d.Format(time.RFC3339Nano)
				}
			}
			existing, existed := fields[field.Name()]
			if existed {
				switch existing := existing.(type) {
				case []interface{}:
					fields[field.Name()] = append(existing, newval)
				case interface{}:
					arr := make([]interface{}, 2)
					arr[0] = existing
					arr[1] = newval
					fields[field.Name()] = arr
				}
			} else {
				fields[field.Name()] = newval
			}
		}
		role := GalaxyRole{}
		mapstructure.Decode(fields, &role)	
		role.ID = id
		docs = append(docs, role)
	}
	return docs
}

func search(keyword string, idx bleve.Index) {
	q := bleve.NewTermQuery(keyword)
	req := bleve.NewSearchRequest(q)
	res, err := idx.Search(req)
	if err != nil {
		log.Fatal(err)
	}
	d := getDocsFromSearchResults(res, idx)
	for _, data := range d {
		fmt.Println(data)
	}
}


func main() {
	os.RemoveAll("/tmp/scorch")

	m := bleve.NewIndexMapping()
	idx, err := bleve.NewUsing("/tmp/scorch", m, scorch.Name, scorch.Name, nil)
	//idx, err := bleve.New("/tmp/scorch", m)
	if err != nil {
		log.Fatal(err)
	}

	batch := idx.NewBatch()

	err = batch.Index("a", createData("test", "docker"))
	if err != nil {
		log.Fatal(err)
	}

	err = batch.Index("b", createData("test", "jenkins"))
	if err != nil {
		log.Fatal(err)
	}

	err = idx.Batch(batch)
	if err != nil {
		log.Fatal(err)
	}

	search("jenkins", idx)
	search("docker", idx)
	idx.Delete("b")
	search("jenkins", idx)
	search("docker", idx)

	idx.Delete("a")
	idx.Delete("b")
	err = idx.Close()
	if err != nil {
		log.Fatal(err)
	}

}






/*
KV Datastore
Setup
CreateDocument
EditDocument
DeleteDocument
GetDocument
Search


