package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/document"
	"github.com/blevesearch/bleve/index/scorch"
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

type UserVotes struct {
	username string
	vote     string
}

type GalaxyRole struct {
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

func getBleveDocsFromSearchResults(
	results *bleve.SearchResult,
	index bleve.Index,
) [][]byte {
	docs := make([][]byte, 0)

	for _, val := range results.Hits {
		id := val.ID
		doc, _ := index.Document(id)

		rv := struct {
			ID     string                 `json:"id"`
			Fields map[string]interface{} `json:"fields"`
		}{
			ID:     id,
			Fields: map[string]interface{}{},
		}
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
			existing, existed := rv.Fields[field.Name()]
			if existed {
				switch existing := existing.(type) {
				case []interface{}:
					rv.Fields[field.Name()] = append(existing, newval)
				case interface{}:
					arr := make([]interface{}, 2)
					arr[0] = existing
					arr[1] = newval
					rv.Fields[field.Name()] = arr
				}
			} else {
				rv.Fields[field.Name()] = newval
			}
		}
		j2, _ := json.Marshal(rv.Fields)
		docs = append(docs, j2)
	}

	return docs
}

func search(keyword string, idx bleve.Index) {
	q := bleve.NewTermQuery(keyword)
	req := bleve.NewSearchRequest(q)
	//req.Fields = []string{"*"}
	//req.Highlight = bleve.NewHighlightWithStyle(ansi.Name)
	//req.Explain = true
	res, err := idx.Search(req)
	if err != nil {
		log.Fatal(err)
	}

	d := getBleveDocsFromSearchResults(res, idx)
	for _, data := range d {

		s := string(data)
		data := GalaxyRole{}
		json.Unmarshal([]byte(s), &data)
		fmt.Printf("Operation: %s \n", data.Repo)
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
