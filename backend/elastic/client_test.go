package elastic

import (
	"errors"
	"testing"

	"github.com/elastic/go-elasticsearch/v8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock mock.Mock
	elasticsearch.Client
}

func (mc *MockClient) Insert(index string, doc map[string]interface{}) error {
	args := mc.mock.Called(index, doc)
	return args.Error(0)
}

func (mc *MockClient) Search(index string, filters []string) (map[string]interface{}, error) {
	args := mc.mock.Called(index, filters)

	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (mc *MockClient) Get(index string, docId string) (DocResponse, error) {
	args := mc.mock.Called(index, docId)

	return args.Get(0).(DocResponse), args.Error(1)
}

func (mc *MockClient) InsertBulk(index string, doc map[string]interface{}) error {
	args := mc.mock.Called(index, doc)
	return args.Error(0)
}

func (mc *MockClient) BulkSearch(queries []map[string]interface{}, index string) (interface{}, error) {
	args := mc.mock.Called(queries, index)

	return args.Get(0).(interface{}), args.Error(1)
}

type testGet struct {
	name    string
	wantErr bool
	input   struct {
		index string
		docId string
	}
	mockCall *mock.Call
	expected DocResponse
}

func TestMockGet(t *testing.T) {
	client := &MockClient{mock: mock.Mock{}}
	tests := []testGet{
		{
			name:    "fail not found",
			wantErr: true,
			input: struct {
				index string
				docId string
			}{index: "", docId: ""},
			mockCall: client.mock.On("Get", "", "").Return(DocResponse{}, errors.New("not found")),
			expected: DocResponse{},
		},
		{
			name:    "success",
			wantErr: false,
			input: struct {
				index string
				docId string
			}{index: "mocks-index", docId: "mocks-id"},
			mockCall: client.mock.On("Get", "mocks-index", "mocks-id").Return(DocResponse{Index: "mocks-index", Id: "mocks-id", Source: map[string]interface{}{"name": "mocks-name"}}, nil),
			expected: DocResponse{Index: "mocks-index", Id: "mocks-id", Source: map[string]interface{}{"name": "mocks-name"}},
		},
	}

	for _, test := range tests {
		res, err := client.Get(test.input.index, test.input.docId)
		assert.Equal(t, err != nil, test.wantErr, test.name)
		assert.Equal(t, test.expected, res, test.name)
	}
}

type testInsert struct {
	name    string
	wantErr bool
	input   struct {
		index string
		doc   map[string]interface{}
	}
	mockCall *mock.Call
}

func TestInsert(t *testing.T) {
	client := &MockClient{mock: mock.Mock{}}
	tests := []testInsert{
		{
			name:    "fail",
			wantErr: true,
			input: struct {
				index string
				doc   map[string]interface{}
			}{index: "",
				doc: map[string]interface{}{
					"doc": Doc{},
				}},
			mockCall: client.mock.On("Insert", "", map[string]interface{}{
				"doc": Doc{},
			}).Return(errors.New("failed to insert - index not found")),
		},
		{
			name:    "success",
			wantErr: false,
			input: struct {
				index string
				doc   map[string]interface{}
			}{index: "mocks-index", doc: map[string]interface{}{
				"doc": Doc{},
			}},
			mockCall: client.mock.On("Insert", "mocks-index", map[string]interface{}{
				"doc": Doc{},
			}).Return(nil),
		},
	}

	for _, test := range tests {
		err := client.Insert(test.input.index, test.input.doc)
		assert.Equal(t, err != nil, test.wantErr, test.name)
	}
}

type testSearch struct {
	name    string
	wantErr bool
	input   struct {
		index   string
		filters []string
	}
	mockCall *mock.Call
	expected map[string]interface{}
}

func TestSearch(t *testing.T) {
	client := MockClient{mock: mock.Mock{}}
	tests := []testSearch{
		{
			name:    "fail on search",
			wantErr: true,
			input: struct {
				index   string
				filters []string
			}{index: "", filters: []string{""}},
			expected: map[string]interface{}{},
			mockCall: client.mock.On("Search", "", []string{""}).Return(map[string]interface{}{}, errors.New("not found")),
		},
		{
			name:    "success",
			wantErr: false,
			input: struct {
				index   string
				filters []string
			}{index: "mocks-index", filters: []string{""}},
			expected: map[string]interface{}{"result": DocResponse{Index: "mocks-index", Id: "mocks-id", Source: map[string]interface{}{"name": "test-doc"}}},
			mockCall: client.mock.On("Search", "mocks-index", []string{""}).Return(map[string]interface{}{"result": DocResponse{Index: "mocks-index", Id: "mocks-id", Source: map[string]interface{}{"name": "test-doc"}}}, nil),
		},
	}

	for _, test := range tests {
		res, err := client.Search(test.input.index, test.input.filters)
		assert.Equal(t, err != nil, test.wantErr, test.name)
		assert.Equal(t, res, test.expected, test.name)
	}
}

type insertBulkTest struct {
	name    string
	wantErr bool
	Input   struct {
		index string
		docs  map[string]interface{}
	}
	mockFunc *mock.Call
}

func TestService_InsertBulk(t *testing.T) {
	client := MockClient{
		mock: mock.Mock{},
	}

	tests := []insertBulkTest{
		{
			name:    "fail test empty index",
			wantErr: true,
			Input: struct {
				index string
				docs  map[string]interface{}
			}{
				index: "",
				docs: map[string]interface{}{
					"mocks-doc": "",
				}},
			mockFunc: client.mock.On("InsertBulk", "", map[string]interface{}{"mocks-doc": ""}).Return(errors.New("failed")),
		},
		{
			name:    "success",
			wantErr: false,
			Input: struct {
				index string
				docs  map[string]interface{}
			}{
				index: "mocks-index",
				docs: map[string]interface{}{
					"mocks-doc": "",
				}},
			mockFunc: client.mock.On("InsertBulk", "mocks-index", map[string]interface{}{"mocks-doc": ""}).Return(nil),
		},
	}

	for _, test := range tests {
		err := client.InsertBulk(test.Input.index, test.Input.docs)
		assert.Equal(t, test.wantErr, err != nil, test.name)
	}
}
