package elastic

import (
	"reflect"

	"github.com/elastic/go-elasticsearch/v8"

	"github.com/elastic/go-elasticsearch/v8/esapi"

	"github.com/ireuven89/hello-world/backend/environment"

	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

//go:generate mockgen -source=client.go -destination=mock/client.go

type Service interface {
	Insert(ctx context.Context, index string, doc interface{}) (string, error)
	InsertBulk(ctx context.Context, index string, docs map[string][]interface{}) error
	Search(ctx context.Context, index string, filters ...string) (SearchResponse, error)
	Get(ctx context.Context, index string, docId string) (DocResponse, error)
	Delete(ctx context.Context, index string, docId string) error
	DeleteIndex(ctx context.Context, index string) error
}

type EsService struct {
	client *elasticsearch.Client
	logger *zap.Logger
}

type Doc struct {
	doc   bytes.Buffer
	index string
}

type SearchResponse struct {
	Took     int16   `json:"took"`
	TimedOut bool    `json:"timed_out"`
	MaxScore float32 `json:"max_score"`
	Hits     struct {
		Hits []DocResponse `json:"hits"`
	} `json:"hits"`
}

type DocResponse struct {
	Index  string                 `json:"_index"`
	Id     string                 `json:"_id"`
	Source map[string]interface{} `json:"_source"`
}

type User struct {
	Age       string `json:"age"`
	CreatedAt string `json:"created_at"`
	Name      string `json:"name"`
}

func New(logger *zap.Logger) (Service, error) {

	//env
	if err := environment.Load(); err != nil {
		return nil, err
	}
	host := environment.Variables.ElasticHost
	userName := environment.Variables.ElasticUsername
	password := environment.Variables.ElasticPassword

	// Configure Elasticsearch client
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{host}, // Elasticsearch Docker URL
		Username:  userName,
		Password:  password,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})

	//ping check
	res, err := es.Ping(es.Ping.WithContext(context.Background()))
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("failed with status code %v", res.StatusCode)
	}
	return &EsService{
		client: es,
		logger: logger,
	}, nil
}

func (s *EsService) Insert(ctx context.Context, index string, doc interface{}) (string, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(doc); err != nil {
		s.logger.Error(fmt.Sprintf("failed insert operation %v", ctx.Value("operation")))
		return "", err
	}
	res, err := s.client.Index(index, &body)

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to insert to %v", ctx.Value("operation")))
		return "", err
	}

	if res.IsError() {
		s.logger.Error(fmt.Sprintf("failed to insert to %v %v", ctx.Value("operation"), res.String()))
		return "", fmt.Errorf("failed initiating server with status code %v", res.StatusCode)
	}

	// Parse the response to retrieve the document ID
	var resBody map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		s.logger.Error(fmt.Sprintf("Error parsing response body: %v", err))
		return "", err
	}

	docID, ok := resBody["_id"].(string)

	if !ok {
		s.logger.Error(fmt.Sprintf("Error parsing response body: %v", err))
		return "", errors.New(fmt.Sprintf("Error parsing doc id from response: %v", err))
	}

	s.logger.Debug(fmt.Sprintf("Insert operation: %v", res.String()))

	return docID, nil
}

func (s *EsService) InsertBulk(ctx context.Context, index string, docs map[string][]interface{}) error {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(docs); err != nil {
		log.Error(err)
		return err
	}
	req := esapi.BulkRequest{
		Body:  &body,
		Index: index,
	}
	response, err := req.Do(context.Background(), s.client.Transport)

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to insert to %v", ctx.Value("operation")))
		return err
	}

	s.logger.Debug("bulk operation: ", zap.String("response: ", response.String()))

	return nil
}

func (s *EsService) Get(ctx context.Context, index string, docId string) (DocResponse, error) {
	var result DocResponse

	res, err := s.client.Get(index, docId)

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to insert to %v", ctx.Value("operation")))
		return DocResponse{}, err
	}

	if res.StatusCode != 200 {
		s.logger.Error(fmt.Sprintf("failed to insert to %v", ctx.Value("operation")))
		return DocResponse{}, errors.New("doc not found")
	}

	err = s.parse(res.Body, &result)

	println(res.Body)

	return result, nil
}

func (s *EsService) BulkSearch(ctx context.Context, queries []map[string]interface{}, index string) (interface{}, error) {
	var buf bytes.Buffer

	// Build the bulk search request body
	for _, query := range queries {
		// Metadata line for each query
		meta := map[string]string{"index": index}
		metaLine, err := json.Marshal(meta)
		if err != nil {
			return nil, fmt.Errorf("error marshalling metadata: %w", err)
		}
		buf.Write(metaLine)
		buf.WriteString("\n")

		// Query line
		queryLine, err := json.Marshal(query)
		if err != nil {
			return nil, fmt.Errorf("error marshalling query: %w", err)
		}
		buf.Write(queryLine)
		buf.WriteString("\n")
	}

	// Execute the bulk search request
	res, err := s.client.Msearch(&buf,
		s.client.Msearch.WithContext(ctx),
		s.client.Msearch.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("error executing bulk search: %w", err)
	}
	defer res.Body.Close()

	// Handle response errors
	if res.IsError() {
		return nil, fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}

	// Parse the response
	var result struct {
		Responses []struct {
			Hits struct {
				Hits []struct {
					Source User `json:"_source"`
				} `json:"hits"`
			} `json:"hits"`
		} `json:"responses"`
	}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding Elasticsearch response: %w", err)
	}

	// Aggregate results
	var users []User
	for _, response := range result.Responses {
		for _, hit := range response.Hits.Hits {
			users = append(users, hit.Source)
		}
	}

	return users, nil
}

func (s *EsService) Search(ctx context.Context, index string, filters ...string) (SearchResponse, error) {
	var result SearchResponse
	response, err := s.client.Search(s.client.Search.WithIndex(index))

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to insert to %v", ctx.Value("operation")))
		return SearchResponse{}, err
	}

	if response.StatusCode != 200 {
		s.logger.Error(fmt.Sprintf("failed to insert to %v", ctx.Value("operation")))
		return SearchResponse{}, errors.New("failed to search")
	}

	err = s.parse(response.Body, &result)

	if err != nil {
		return SearchResponse{}, err
	}

	return result, nil
}

func (s *EsService) Delete(ctx context.Context, index string, docId string) error {
	res, err := s.client.Delete(index, docId, s.client.Delete.WithContext(ctx))

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to delete doc:  %v %v", docId, err))
		return err
	}

	if res.IsError() {
		s.logger.Error(fmt.Sprintf("failed to delete doc:  %v %v", docId, err))
		return fmt.Errorf("failed to delete doc status code %v message is %v", res.StatusCode, res.String())
	}

	return nil
}

func (s *EsService) DeleteIndex(ctx context.Context, index string) error {
	res, err := s.client.Indices.Delete([]string{index}, s.client.Indices.Delete.WithContext(ctx))

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to delete doc:  %v %v", index, err))
		return err
	}

	if res.IsError() {
		s.logger.Error(fmt.Sprintf("failed to delete doc:  %v %v", index, err))
		return fmt.Errorf("failed to delete doc status code %v message is %v", res.StatusCode, res.String())
	}

	return nil
}

func (s *EsService) parse(reader io.ReadCloser, obj interface{}) error {

	body, err := io.ReadAll(reader)

	if err != nil {
		return err
	}

	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf("obj must be a non-nil pointer")
	}

	// Unmarshal JSON into the provided object
	return json.Unmarshal(body, obj)

}
