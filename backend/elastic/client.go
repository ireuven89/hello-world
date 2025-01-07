package elastic

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/ireuven89/hello-world/backend/environment"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

type DocResponse struct {
	Index  string `json:"_index"`
	Id     string `json:"_id"`
	Source source `json:"_source"`
}

type source struct {
	Name      string `json:"name"`
	Age       string `json:"age"`
	CreatedAt string `json:"created_at"`
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
		logger: zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(zapcore.EncoderConfig{}), nil, nil)),
	}, nil
}

func (s *Service) Insert(ctx context.Context, index string, doc map[string]interface{}) error {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(doc); err != nil {
		log.Fatalf("Error encoding document: %s", err)
	}
	res, err := s.client.Index(index, &body)

	if err != nil {
		return err
	}

	s.logger.Debug("Set operation: ", zap.Any("response is ", res))

	return nil
}

func (s *Service) InsertBulk(index string, docs map[string][]interface{}) error {
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
		return err
	}

	s.logger.Debug("bulk operation: ", zap.String("response: ", response.String()))

	return nil
}

func (s *Service) Get(index string, docId string) (DocResponse, error) {
	res, err := s.client.Get(index, docId)

	if err != nil {
		return DocResponse{}, err
	}

	if res.StatusCode != 200 {
		return DocResponse{}, errors.New("doc not found")
	}

	result, err := s.parseDoc(res.Body, DocResponse{})

	println(res.Body)

	return result.(DocResponse), nil
}

func (s *Service) Search(index string, filters ...string) ([]DocResponse, error) {
	response, err := s.client.Search(s.client.Search.WithIndex(index))

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New("failed to search")
	}

	result, err := s.parseDoc(response.Body, []DocResponse{})

	if err != nil {
		return nil, err
	}

	return result.([]DocResponse), nil
}

func (s *Service) parseDoc(reader io.ReadCloser, t any) (any, error) {

	body, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &t)

	return t, nil
}
