package infra

import (
	"github.com/elastic/go-elasticsearch/v8"
)

func InitElastic() (*elasticsearch.TypedClient, error) {
	es, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
		// Username:
		// Password:
	})
	return es, err
}
