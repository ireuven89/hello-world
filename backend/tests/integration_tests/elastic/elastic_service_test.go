package elastic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/ireuven89/hello-world/backend/elastic"
	"github.com/ireuven89/hello-world/backend/tests/config"
	"github.com/ireuven89/hello-world/backend/tests/utils"
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

var esService *elastic.Service
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

	es, err := elastic.New()

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

	err = esService.Insert(ctx, indexName, jsonDoc)

	assert.Nil(t, err, "failed inserting")
}

func TestElasticSearchByIndex(t *testing.T) {
	res, err := esService.Search(ctx, indexName)

	assert.Nil(t, err, "failed search")
	assert.NotEmpty(t, res, "failed search")

}

func TestElasticGet(t *testing.T) {
	esUrl := configJson.ElasticUrlLocal

	client := utils.NewHttpClient()

	url := esUrl + getEndpoint + "/" + indexName
	resp, err := client.Get(url)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.NotEmptyf(t, resp.Body, "failed to get response")
	response, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	res := string(response)
	fmt.Print(res)
}

func TestElasticInsert(t *testing.T) {
	doc := map[string]interface {
	}{
		"Name":   "Alice",
		"Age":    30,
		"Origin": "North",
	}
	url := fmt.Sprintf(InsertEndpoint, configJson.ElasticUrlLocal)
	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(doc)

	assert.Nil(t, err)

	response, err := httpClient.Post(url, "application/json", &body)

	assert.Nil(t, err)
	assert.Equal(t, response.StatusCode, 201)

}
