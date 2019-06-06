package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/document"
	"github.com/blevesearch/bleve/mapping"
	"github.com/mitchellh/mapstructure"
)

func buildMapping() mapping.IndexMapping {
	enFieldMapping := bleve.NewTextFieldMapping()
	enFieldMapping.Analyzer = en.AnalyzerName

	eventMapping := bleve.NewDocumentMapping()
	eventMapping.AddFieldMappingsAt("Readme", enFieldMapping)
	eventMapping.AddFieldMappingsAt("Meta", enFieldMapping)

	kwFieldMapping := bleve.NewTextFieldMapping()

	eventMapping.AddFieldMappingsAt("Server", kwFieldMapping)
	eventMapping.AddFieldMappingsAt("Namespace", kwFieldMapping)
	eventMapping.AddFieldMappingsAt("Repo", kwFieldMapping)

	mapping := bleve.NewIndexMapping()
	mapping.DefaultMapping = eventMapping
	mapping.DefaultAnalyzer = en.AnalyzerName

	return mapping
}

//NebularDoc stores the structure of the data in the bleve datastore
type NebularDoc struct {
	ID        string `json:"ID"`
	Namespace string `json:"Namespace"`
	Meta      string `json:"Meta"`
	Rated     int    `json:"Rated"`
	Readme    string `json:"Readme"`
	Repo      string `json:"Repo"`
	Server    string `json:"Server"`
	MetaType  string `json:"MetaType"`
}

func (role GalaxyRole) ToDoc() NebularDoc {
	meta, _ := json.Marshal(role.Meta)
	fmt.Println(string(meta))
	return NebularDoc{
		ID:        role.ID,
		Namespace: role.Namespace,
		Meta:      string(meta),
		Rated:     role.Rated,
		Readme:    role.Readme,
		Repo:      role.Repo,
		Server:    role.Server,
		MetaType:  role.MetaType,
	}
}

func createDescription(description string, server string, namespace string, repo string) string {
	return fmt.Sprintf("%s | Use the following url: 'git+%s/%s/%s'", description, server, namespace, repo)
}

func (doc NebularDoc) ToRole() GalaxyRole {
	role := GalaxyRole{
		ID:        doc.ID,
		Namespace: doc.Namespace,
		Rated:     doc.Rated,
		Readme:    doc.Readme,
		Repo:      doc.Repo,
		Server:    doc.Server,
		MetaType:  doc.MetaType,
		Name:      doc.Repo,
	}
	switch doc.MetaType {
	case "COMPLEX":
		meta := GalaxyMetaComplex{}
		json.Unmarshal([]byte(doc.Meta), &meta)
		role.Meta = meta
		role.Username = meta.GalaxyInfo.Author
		role.Description = createDescription(meta.GalaxyInfo.Description, doc.Server, doc.Namespace, doc.Repo)
	case "SIMPLE":
		meta := GalaxyMetaSimple{}
		json.Unmarshal([]byte(doc.Meta), &meta)
		role.Meta = meta
		role.Username = meta.GalaxyInfo.Author
		role.Description = createDescription(meta.GalaxyInfo.Description, doc.Server, doc.Namespace, doc.Repo)
	}
	return role
}

type KeyValueDB struct {
	idx  bleve.Index
	path string
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func NewKVDB(path string) *KeyValueDB {
	kv := KeyValueDB{path: path}
	if ex, err := exists(path); err != nil || !ex {
		fmt.Println("CREATING NEW MAPPPING")
		idx, err := bleve.New(path, buildMapping())
		if err != nil {
			log.Fatal(err)
		}
		kv.idx = idx
	} else {
		index, _ := bleve.Open(path)
		kv.idx = index
	}
	return &kv
}

func (kv *KeyValueDB) Close() {
	err := kv.idx.Close()
	if err != nil {
		log.Fatal(err)
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (kv *KeyValueDB) Save(role GalaxyRole) error {
	id := fmt.Sprintf("%s-%s-%s", role.Namespace, role.Repo, RandStringRunes(5))
	batch := kv.idx.NewBatch()
	err := batch.Index(id, role.ToDoc())
	if err != nil {
		return err
	}
	err = kv.idx.Batch(batch)
	if err != nil {
		return err
	}
	return nil
}

func (kv *KeyValueDB) Edit() {}

func (kv *KeyValueDB) Delete(id string) {
	fmt.Println("Deleting document id: " + id)
	kv.idx.Delete(id)
}

func getRoleFromDoc(doc *document.Document) GalaxyRole {
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

	nd := NebularDoc{}
	mapstructure.Decode(fields, &nd)
	nd.ID = doc.ID
	return nd.ToRole()
}
func (kv *KeyValueDB) SearchTerm(keyword string) []GalaxyRole {
	q := bleve.NewMatchQuery(keyword)
	req := bleve.NewSearchRequest(q)
	res, _ := kv.idx.Search(req)

	docs := make([]GalaxyRole, 0)
	for _, val := range res.Hits {
		id := val.ID
		doc, _ := kv.idx.Document(id)
		docs = append(docs, getRoleFromDoc(doc))
	}
	return docs
}

func (kv *KeyValueDB) IndexSize() (int, error) {
	query := bleve.NewMatchAllQuery()
	sizeRequest := bleve.NewSearchRequest(query)
	sizeRequest.Size = 0
	results, err := kv.idx.Search(sizeRequest)
	if err != nil {
		return 0, err
	}
	return int(results.Total), nil
}

func (kv *KeyValueDB) SearchAll() []GalaxyRole {
	docs := make([]GalaxyRole, 0)

	q := bleve.NewMatchAllQuery()
	req := bleve.NewSearchRequest(q)
	req.Size, _ = kv.IndexSize()
	res, err := kv.idx.Search(req)
	if err != nil {
		fmt.Println(err)
		return docs
	}
	fmt.Println(res)
	for _, val := range res.Hits {
		id := val.ID
		doc, _ := kv.idx.Document(id)
		docs = append(docs, getRoleFromDoc(doc))
	}
	return docs
}

func (kv *KeyValueDB) Get(id string) GalaxyRole {
	doc, _ := kv.idx.Document(id)
	return getRoleFromDoc(doc)
}
