package bider

import (
	"errors"
	"github.com/ireuven89/hello-world/backend/bider/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockDb struct {
	mock mock.Mock
}

func (mdb *MockDb) List(input model.BiddersInput) ([]Bidder, error) {
	args := mdb.mock.Called(input)

	return args.Get(0).([]Bidder), args.Error(1)
}

func (mdb *MockDb) Single(input model.BiddersInput) (Bidder, error) {
	args := mdb.mock.Called(input)

	return args.Get(0).(Bidder), args.Error(1)
}

func (mdb *MockDb) Upsert(input model.BiddersInput) (string, error) {
	args := mdb.mock.Called(input)

	return args.Get(0).(string), args.Error(1)
}

func (mdb *MockDb) Delete(uuid string) error {
	args := mdb.mock.Called(uuid)

	return args.Error(0)
}

type TestList struct {
	name     string
	wantErr  bool
	input    model.BiddersInput
	mockCall *mock.Call
	expected []Bidder
}

func TestRepository_List(t *testing.T) {
	repo := MockDb{mock: mock.Mock{}}
	successInput := model.BiddersInput{
		Name: "name",
		Uuid: "uuid",
		Item: "item",
	}
	successExpected := []Bidder{{ID: 0, Name: "name", Uuid: "uuid"}}
	failedInput := model.BiddersInput{
		Name: "",
		Uuid: "",
		Item: "",
	}
	tests := []TestList{{
		name:     "success",
		input:    successInput,
		mockCall: repo.mock.On("List", successInput).Return(successExpected, nil),
		expected: successExpected,
	},
		{
			name:     "fail on empty list",
			input:    failedInput,
			mockCall: repo.mock.On("List", failedInput).Return([]Bidder(nil), errors.New("not found")),
			wantErr:  true,
			expected: nil,
		},
	}

	for _, test := range tests {
		res, err := repo.List(test.input)
		assert.Equal(t, err != nil, test.wantErr)
		assert.Equal(t, res, test.expected)
	}
}

type TestSingle struct {
	name     string
	wantErr  bool
	input    model.BiddersInput
	mockCall *mock.Call
	expected Bidder
}

func TestRepository_Single(t *testing.T) {
	repo := MockDb{mock: mock.Mock{}}
	successInput := model.BiddersInput{
		Name: "name",
		Uuid: "uuid",
		Item: "item",
	}
	successExpected := Bidder{ID: 0, Name: "name", Uuid: "uuid"}
	tests := []TestSingle{{
		name:     "success",
		input:    successInput,
		mockCall: repo.mock.On("Single", successInput).Return(successExpected, nil),
		expected: successExpected,
	},
	}

	for _, test := range tests {
		res, err := repo.Single(test.input)
		assert.Equal(t, err != nil, test.wantErr)
		assert.Equal(t, res, test.expected)
	}
}

type TestDelete struct {
	name     string
	wantErr  bool
	input    string
	mockCall *mock.Call
}

func TestRepository_Delete(t *testing.T) {
	repo := MockDb{mock: mock.Mock{}}
	successInput := "uuid"
	tests := []TestDelete{
		{
			name:     "success",
			input:    successInput,
			mockCall: repo.mock.On("Delete", successInput).Return(nil),
			wantErr:  false,
		},
		{
			name:     "failed not found",
			input:    "",
			mockCall: repo.mock.On("Delete", "").Return(errors.New("failed to find bidder")),
			wantErr:  true,
		},
	}

	for _, test := range tests {
		err := repo.Delete(test.input)
		assert.Equal(t, err != nil, test.wantErr)

	}
}

type TestUpsert struct {
	name     string
	wantErr  bool
	input    model.BiddersInput
	mockCall *mock.Call
	expected string
}

func TestRepository_Upsert(t *testing.T) {
	repo := MockDb{mock: mock.Mock{}}
	successInput := model.BiddersInput{Name: "test_name", Uuid: "uuid", Item: "item"}
	tests := []TestUpsert{
		{
			name:     "success",
			input:    successInput,
			mockCall: repo.mock.On("Upsert", successInput).Return("mock-uuid", nil),
			wantErr:  false,
			expected: "mock-uuid",
		},
		{
			name:     "failed not found",
			input:    model.BiddersInput{},
			mockCall: repo.mock.On("Upsert", model.BiddersInput{}).Return("", errors.New("failed to find bidder")),
			wantErr:  true,
			expected: "",
		},
	}

	for _, test := range tests {
		res, err := repo.Upsert(test.input)
		assert.Equal(t, err != nil, test.wantErr)
		assert.Equal(t, res, test.expected)
	}
}
