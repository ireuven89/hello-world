package elastic

import (
	"reflect"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/ireuven89/hello-world/backend/environment"

	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

//go:generate mockgen -source=client.go -destination=mock/client.go

type elasticService interface {
	Insert(index string, doc map[string]interface{}) error
	InsertBulk(index string, doc map[string]interface{}) error
	Search(index string, filters ...string) (map[string]interface{}, error)
	Get(index string, docId string) (DocResponse, error)
}

type Service struct {
	client *elasticsearch.Client
	api    *esapi.API
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

func New() (*Service, error) {

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

	//logger
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := loggerConfig.Build()

	if err != nil {
		log.Error("failed to connect to es cluster:", err)
		return nil, err
	}

	// Document to index
	doc := map[string]interface{}{
		"name":       "Alice",
		"age":        "30",
		"created_at": time.Now(),
	}

	// Convert document to JSON
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(doc); err != nil {
		log.Fatalf("Error encoding document: %s", err)
	}

	if err != nil {
		log.Errorf("failed to initiate client")
	}

	api := esapi.New(es.Transport)

	exreq := esapi.ExistsRequest{Index: "my-index", DocumentID: "1"}
	reuqest := esapi.IndexRequest{
		Index: "my-index",
		Body:  &body,
	}

	response, err := exreq.Do(context.Background(), es)

	if err != nil {
		log.Errorf("error checking index ", err)
	}
	log.Info(response)

	result, err := reuqest.Do(context.Background(), es)

	if err != nil {
		return nil, err
	}

	log.Info(result)

	return &Service{
		client: es,
		api:    api,
		logger: logger,
	}, nil
}

func (s *Service) Insert(ctx context.Context, index string, doc interface{}) error {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(doc); err != nil {
		s.logger.Error(fmt.Sprintf("failed insert operation %v", ctx.Value("operation")))
		return err
	}
	res, err := s.client.Index(index, &body)

	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to insert to %v", ctx.Value("operation")))
		return err
	}

	s.logger.Debug(fmt.Sprintf("Insert operation: %v", res.Body))

	return nil
}

func (s *Service) InsertBulk(ctx context.Context, index string, docs map[string][]interface{}) error {
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

func (s *Service) Get(ctx context.Context, index string, docId string) (DocResponse, error) {
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

func (s *Service) Search(ctx context.Context, index string, filters ...string) (SearchResponse, error) {
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

func (s *Service) parse(reader io.ReadCloser, obj interface{}) error {

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
