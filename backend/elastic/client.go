package elastic

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/labstack/gommon/log"

	"bytes"
	"encoding/json"
)

func main() {
	// Configure Elasticsearch client
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"}, // Elasticsearch Docker URL
	})

	// Document to index
	doc := map[string]interface{}{
		"name":       "Alice",
		"age":        30,
		"created_at": "2024-11-26",
	}

	// Convert document to JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		log.Fatalf("Error encoding document: %s", err)
	}

	if err != nil {
		log.Errorf("failed to initiate client")
	}

	res, err := es.Index(
		"my-index",                   // Index name
		&buf,                         // Document body
		es.Index.WithDocumentID("1"), // Optional: Specify document ID
	)

	if err != nil {
		log.Fatalf("Error indexing document: %s", err)
	}

	defer res.Body.Close()
}
