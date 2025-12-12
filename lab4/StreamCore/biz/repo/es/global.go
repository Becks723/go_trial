package es

import (
	"log"

	elastic "github.com/elastic/go-elasticsearch/v8"
)

var esRaw *elastic.TypedClient

func Init() {
	es, err := elastic.NewTypedClient(elastic.Config{
		Addresses: []string{"http://localhost:9200"},
		// Username:
		// Password:
	})
	if err != nil {
		log.Fatal(err)
	}
	esRaw = es
}
