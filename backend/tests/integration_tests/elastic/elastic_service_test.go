package elastic

import (
	"encoding/json"
	"go.uber.org/zap/zaptest"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/ireuven89/hello-world/backend/elastic"
	"github.com/ireuven89/hello-world/backend/tests/config"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var configJson config.ConfigurationJson
var httpClient http.Client

const (
	indexName      = "articles"
	docName        = "test-doc"
	getEndpoint    = "/_cat/indices"
	InsertEndpoint = "%s/test-index/_doc"
	SearchEndpoint = "%s/test-index/_search"
)

type ElasticSearchResponse struct {
	Took int  `json:"took"`
	Hits Hits `json:"hits"`
}

type Hits struct {
	Score float32 `json:"max_score"`
	Hits  []Hit   `json:"hits"`
}

type Hit struct {
	Index string `json:"_index"`
	Type  string `json:"_type"`
	Id    string `json:"_id"`
}

var esService elastic.Service
var ctx = context.Background()

func init() {

	if err := os.Setenv("ELASTIC_HOST", "http://localhost:9200"); err != nil {
		panic("failed setting env for tests")
	}

	if err := os.Setenv("ELASTIC_USER", "elastic"); err != nil {
		panic("failed setting env for tests")
	}

	if err := os.Setenv("ELASTIC_PASSWORD", "vwS2rX-79_ReyjLB"); err != nil {
		panic("failed setting env for tests")
	}

	t := testing.T{}
	logger := zaptest.NewLogger(&t)

	es, err := elastic.New(logger)

	if err != nil {
		panic("failed initialize service")
	}

	esService = es
}

func TestInsert(t *testing.T) {
	file, err := os.Open("./files/doc.json")

	if err != nil {
		assert.Fail(t, "failed to open file for test name")
	}

	body, err := ioutil.ReadAll(file)
	if err != nil {
		assert.Fail(t, "failed to parse file")
	}

	// Parse JSON into a map
	var jsonDoc map[string]interface{}
	if err = json.Unmarshal(body, &jsonDoc); err != nil {
		assert.Fail(t, "failed to parse file")
	}

	_, err = esService.Insert(ctx, indexName, jsonDoc)

	assert.Nil(t, err, "failed inserting")
}

func TestElasticSearchByIndex(t *testing.T) {
	res, err := esService.Search(ctx, indexName)

	assert.Nil(t, err, "failed search")
	assert.NotEmpty(t, res, "failed search")

}

func TestDeleteDocByIndex(t *testing.T) {
	err := esService.Delete(ctx, indexName, docName)

	assert.Nil(t, err, "failed delete")

}

func TestDeleteIndex(t *testing.T) {
	err := esService.DeleteIndex(ctx, indexName)

	assert.Nil(t, err, "failed delete")

}
