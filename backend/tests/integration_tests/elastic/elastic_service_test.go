package elastic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"go.uber.org/zap/zaptest"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"github.com/ireuven89/hello-world/backend/elastic"
	"github.com/ireuven89/hello-world/backend/tests/config"
)

var configJson config.ConfigurationJson
var httpClient http.Client

const (
	indexName = "articles"
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
var docId string

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

<<<<<<< Updated upstream
	es, err := elastic.New(logger)
=======
	config, err := utils.LoadConfig("elastic", "ENV")
	if err != nil {
		t.Error(err)
	}
	es, err := elastic.New(logger, config.Databases["elastic"])
>>>>>>> Stashed changes

	if err != nil {
		panic(fmt.Sprintf("failed initialize service %v", err))
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

	Id, err := esService.Insert(ctx, indexName, jsonDoc)

	assert.Nil(t, err, "failed inserting")
	docId = Id
}

func TestElasticSearchByIndex(t *testing.T) {
	res, err := esService.Search(ctx, indexName)

	assert.Nil(t, err, "failed search")
	assert.NotEmpty(t, res, "failed search")

}

func TestDeleteDocByIndex(t *testing.T) {
	err := esService.Delete(ctx, indexName, docId)

	assert.Nil(t, err, "failed delete")
	t.Log("finished test delete doc")

}

func TestDeleteIndex(t *testing.T) {
	err := esService.DeleteIndex(ctx, indexName)

	assert.Nil(t, err, "failed delete test")
}

func TestMain(m *testing.M) {
	exitCode := m.Run() // Runs all tests
	os.Exit(exitCode)
}
