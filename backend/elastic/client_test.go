package elastic

import (
	"errors"
	"github.com/golang/mock/gomock"
	mockelastic "github.com/ireuven89/hello-world/backend/elastic/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockClient struct {
	mock.Mock
}

func (mc *MockClient) Insert(index string, doc map[string]interface{}) error {
	args := mc.Mock.Called(index, doc)
	return args.Error(0)
}

func (mc *MockClient) Search(index string, filters ...string) (map[string]interface{}, error) {
	args := mc.Mock.Called(index, filters)

	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (mc *MockClient) Get(index string, docId string) (DocResponse, error) {
	args := mc.Mock.Called(index, docId)

	return args.Get(0).(DocResponse), args.Error(1)
}

func (mc *MockClient) InsertBulk(index string, doc map[string]interface{}) error {
	args := mc.Mock.Called(index, doc)
	return args.Error(0)
}

func TestMockGetFail(t *testing.T) {
	client := &MockClient{Mock: mock.Mock{}}
	client.Mock.On("Get", "test", "id").Return(DocResponse{}, errors.New("failed"))

	res, err := client.Get("test", "id")

	assert.NotNil(t, err)
	assert.Emptyf(t, res, "failed get")
}

func TestMockGet(t *testing.T) {
	client := &MockClient{Mock: mock.Mock{}}

	expected := DocResponse{
		Index:  "test",
		Id:     "id",
		Source: source{Name: "demo", Age: "20", CreatedAt: "2020-01-12"}}
	client.Mock.On("Get", "test", "id").Return(expected, nil)

	res, err := client.Get("test", "id")
	assert.Nil(t, err)
	assert.Equal(t, res, expected)
}

func TestInsert(t *testing.T) {
	client := &MockClient{Mock: mock.Mock{}}
	docArg := map[string]interface{}{
		"index": "mock-index",
		"name":  "mock-name",
	}
	indexArg := "mock-index"
	client.Mock.On("Insert", indexArg, docArg).Return(nil)
	err := client.Insert(indexArg, docArg)

	assert.Nil(t, err)
}

func TestSearch(t *testing.T) {
	expectedDoc := DocResponse{Index: "mock", Id: "mock", Source: source{
		Name: "mock",
	}}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockelastic.NewMockelasticService(ctrl)

	mockClient.EXPECT().Search("mock-index", "mock-filter").Return([]DocResponse{
		{Index: "mock", Id: "mock", Source: source{Name: "mock"}},
	}, nil).Times(1)

	//our real client
	client, err := New()

	res, err := client.Search("mock-index", "mock-filter")

	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)
	assert.Equal(t, res[1], expectedDoc)
}

func TestService_InsertBulk(t *testing.T) {
	expectedDoc := DocResponse{Index: "mock", Id: "mock", Source: source{
		Name: "mock",
	}}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockelastic.NewMockelasticService(ctrl)

	mockClient.EXPECT().Search("mock-index", "mock-filter").Return([]DocResponse{
		{Index: "mock", Id: "mock", Source: source{Name: "mock"}},
	}, nil).Times(1)

	//our real client
	client, err := New()

	res, err := client.Search("mock-index", "mock-filter")

	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)
	assert.Equal(t, res[1], expectedDoc)
}
