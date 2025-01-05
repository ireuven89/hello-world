package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ireuven89/hello-world/backend/tests/config"
	"github.com/ireuven89/hello-world/backend/tests/utils"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

var configJson config.ConfigurationJson
var httpClient http.Client

const (
	indexName      = "test-index"
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

func init() {
	configJson = config.GetConfigJson()
	httpClient = utils.NewHttpClient()
}

func TestElasticSearchByIndex(t *testing.T) {
	var searchResponse ElasticSearchResponse
	url := fmt.Sprintf(SearchEndpoint, configJson.ElasticUrlLocal)

	response, err := httpClient.Get(url)

	assert.Nil(t, err)
	assert.NotEmpty(t, response.Body)

	resp, err := io.ReadAll(response.Body)

	assert.NotEmpty(t, resp)

	err = json.Unmarshal(resp, &searchResponse)
	assert.Nil(t, err)
	assert.True(t, len(searchResponse.Hits.Hits) > 0)

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
