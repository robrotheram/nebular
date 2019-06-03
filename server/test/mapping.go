package main

import (
	"bytes"
	"fmt"
	"log"
	"net/mail"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis"
	"github.com/blevesearch/bleve/analysis/datetime/flexible"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/bleve/registry"
	"github.com/blevesearch/bleve/search/query"
)

const email2 = `Subject: test email
To: Bob Smith <bob@example.com>
From: Jane Smith <jane@example.com>
Date: Tue, 29 Nov 2016 21:48:47 +1000

This is a test email
`
const email1 = `Subject: blah boo
To: Bob Smith <bob@example.com>
From: Jane Smith <jane@example.com>
Date: Wed, 2 Nov 2016 21:48:47 +1000

This is a test email
`

type bleveDoc struct {
	Type string
	Data mail.Header
}

func main() {
	var err error
	var index bleve.Index

	mapping := buildIndexMapping()

	index, err = bleve.NewMemOnly(mapping)
	if err != nil {
		log.Fatal(err)
	}

	var msg *mail.Message
	msg, err = mail.ReadMessage(bytes.NewReader([]byte(email1)))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("header %v\n", msg.Header)

	doc := bleveDoc{"header", msg.Header}

	if err := index.Index("testdoc1", doc); err != nil {
		log.Fatal(err)
	}

	msg, err = mail.ReadMessage(bytes.NewReader([]byte(email2)))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("header %v\n", msg.Header)

	doc = bleveDoc{"header", msg.Header}

	if err = index.Index("testdoc2", doc); err != nil {
		log.Fatal(err)
	}

	bSearchRequest := bleve.NewSearchRequest(query.NewDateRangeQuery(
		time.Date(2016, 11, 20, 0, 0, 0, 0, time.Local),
		time.Date(2016, 12, 30, 0, 0, 0, 0, time.Local),
	))
	searchResult, err := index.Search(bSearchRequest)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Hits:")
	for _, hit := range searchResult.Hits {
		fmt.Printf("ID %s, Locations %v, Fields %v\n", hit.ID, hit.Locations, hit.Fields)
	}

}

func buildIndexMapping() mapping.IndexMapping {
	mapping := bleve.NewIndexMapping()

	docMapping := bleve.NewDocumentMapping()

	headerMapping := bleve.NewDocumentMapping()
	dateFieldMapping := bleve.NewDateTimeFieldMapping()
	dateFieldMapping.DateFormat = dateTimeParserName
	headerMapping.AddFieldMappingsAt("Date", dateFieldMapping)
	docMapping.AddFieldMappingsAt("Date", dateFieldMapping)
	docMapping.AddSubDocumentMapping("Data", headerMapping)

	mapping.AddDocumentMapping("header", docMapping)
	mapping.TypeField = "Type"

	return mapping
}

const dateTimeParserName = "dateTimeParser"
const RFC1123ZnoPadDay = "Mon, _2 Jan 2006 15:04:05 -0700"

func init() {
	registry.RegisterDateTimeParser(dateTimeParserName, DateTimeParserConstructor)
}

var dateTimeParserLayouts = []string{
	RFC1123ZnoPadDay,
}

func DateTimeParserConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.DateTimeParser, error) {
	return flexible.New(dateTimeParserLayouts), nil
}
