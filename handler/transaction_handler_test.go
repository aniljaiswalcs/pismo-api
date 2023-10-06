package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/aniljaiswalcs/pismo/model"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) CreateTransaction(ctx context.Context, transaction model.Transaction) (*model.Transaction, error) {
	args := m.Called(ctx, transaction)
	return args.Get(0).(*model.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) SubtractTransaction(ctx context.Context, transaction model.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func TestCreateTransaction(t *testing.T) {
	mockRepo := new(MockTransactionRepository)

	payload := `{"account_id": 123456789, "operation_type_id": 1, "amount": -100.0}`
	req, _ := http.NewRequest("POST", "/accounts", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	expectedTransaction := &model.Transaction{AccountId: 123456789, OperationTypeId: 1, Amount: 100.0}
	mockRepo.On("CreateTransaction", mock.Anything, mock.AnythingOfType("model.Transaction")).Return(expectedTransaction, nil)

	handler := &TransactionHandler{repository: mockRepo}
	handler.CreateTransaction(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d but got %d", http.StatusCreated, w.Code)
	}

	expectedResponse := `{"transaction_id":0,"account_id":123456789,"operation_type_id":1,"amount":100}`
	actualResponse := w.Body.String()

	expectedResponseJson := map[string]string{}
	actualResponseJson := map[string]string{}

	json.Unmarshal([]byte(expectedResponse), &expectedResponseJson)
	json.Unmarshal([]byte(actualResponse), &actualResponseJson)

	if !reflect.DeepEqual(expectedResponseJson, actualResponseJson) {
		t.Errorf("Expected response body %s but got %s", expectedResponse, w.Body.String())
	}
}

func TestCreateTransactionFailsWhenInvalidRequest(t *testing.T) {
	var scenarios = []struct {
		payload            string
		expectedResponse   string
		expectedStatusCode int
	}{
		{
			`{"account_id": 123456789, "operation_type_id": 1, "amount": 100.0}`,
			`{"purchases and withdraw operations must have a negative amount. Payment operations must have a positive amount."}`,
			http.StatusBadRequest,
		},
		{
			`{"account_id": 123456789, "operation_type_id": 2, "amount": 100.0}`,
			`{"purchases and withdraw operations must have a negative amount. Payment operations must have a positive amount."}`,
			http.StatusBadRequest,
		},
		{
			`{"account_id": 123456789, "operation_type_id": 3, "amount": 100.0}`,
			`{"purchases and withdraw operations must have a negative amount. Payment operations must have a positive amount."}`,
			http.StatusBadRequest,
		},
		{
			`{"account_id": 123456789, "operation_type_id": 4, "amount": -100.0}`,
			`{"purchases and withdraw operations must have a negative amount. Payment operations must have a positive amount."}`,
			http.StatusBadRequest,
		},
		{
			`{"account_id": "invalid", "operation_type_id": 1, "amount": -100.0}`,
			`{"The account_id and operation_type_id must be valid positive integers. The amount must be a valid decimal."}`,
			http.StatusBadRequest,
		},
		{
			`{"account_id": -1, "operation_type_id": 1, "amount": -100.0}`,
			`{"The account_id and operation_type_id must be valid positive integers. The amount must be a valid decimal."}`,
			http.StatusBadRequest,
		},
		{
			`{"account_id": 0, "operation_type_id": 1, "amount": -100.0}`,
			`{"The account_id must be a valid positive integer."}`,
			http.StatusBadRequest,
		},
		{
			`{"account_id": null, "operation_type_id": 1, "amount": -100.0}`,
			`{"The account_id must be a valid positive integer."}`,
			http.StatusBadRequest,
		},
		{
			`{"account_id": 123456789, "operation_type_id": 0, "amount": -100.0}`,
			`{"The operation_type_id must be one of the following valid values: 1, 2, 3, 4"}`,
			http.StatusBadRequest,
		},
		{
			`{"account_id": 123456789, "operation_type_id": 5, "amount": -100.0}`,
			`{"The operation_type_id must be one of the following valid values: 1, 2, 3, 4"}`,
			http.StatusBadRequest,
		},
		{
			`{"account_id": 123456789, "operation_type_id": -1, "amount": -100.0}`,
			`{"The account_id and operation_type_id must be valid positive integers. The amount must be a valid decimal."}`,
			http.StatusBadRequest,
		},
		{
			`{"account_id": 123456789, "operation_type_id": null, "amount": -100.0}`,
			`{"The operation_type_id must be one of the following valid values: 1, 2, 3, 4"}`,
			http.StatusBadRequest,
		},
		{
			`{"account_id": 123456789, "operation_type_id": "invalid", "amount": -100.0}`,
			`{"The account_id and operation_type_id must be valid positive integers. The amount must be a valid decimal."}`,
			http.StatusBadRequest,
		},
	}

	for _, scenario := range scenarios {
		mockRepo := new(MockTransactionRepository)

		payload := scenario.payload
		req, _ := http.NewRequest("POST", "/accounts", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler := &TransactionHandler{repository: mockRepo}
		handler.CreateTransaction(w, req)

		if w.Code != scenario.expectedStatusCode {
			t.Errorf("Expected status code %d but got %d", http.StatusCreated, w.Code)
		}

		expectedResponse := scenario.expectedResponse
		actualResponse := w.Body.String()

		expectedResponseJson := map[string]string{}
		actualResponseJson := map[string]string{}

		json.Unmarshal([]byte(expectedResponse), &expectedResponseJson)
		json.Unmarshal([]byte(actualResponse), &actualResponseJson)

		if !reflect.DeepEqual(expectedResponseJson, actualResponseJson) {
			t.Errorf("Expected response body %s but got %s", expectedResponse, w.Body.String())
		}
	}
}

func TestCreateTransactionWhenTransactionCreatonFails(t *testing.T) {
	mockRepo := new(MockTransactionRepository)

	payload := `{"account_id": 123456789, "operation_type_id": 1, "amount": -100.0}`
	req, _ := http.NewRequest("POST", "/accounts", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockRepo.On("CreateTransaction", mock.Anything, mock.AnythingOfType("model.Transaction")).Return(&model.Transaction{}, errors.New("Error!"))

	handler := &TransactionHandler{repository: mockRepo}
	handler.CreateTransaction(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d but got %d", http.StatusCreated, w.Code)
	}

	expectedResponse := `{"no account found for the provided account ID"}`
	actualResponse := w.Body.String()

	expectedResponseJson := map[string]string{}
	actualResponseJson := map[string]string{}

	json.Unmarshal([]byte(expectedResponse), &expectedResponseJson)
	json.Unmarshal([]byte(actualResponse), &actualResponseJson)

	if !reflect.DeepEqual(expectedResponseJson, actualResponseJson) {
		t.Errorf("Expected response body %s but got %s", expectedResponse, w.Body.String())
	}
}

func TestNewTransactionHandler(t *testing.T) {
	repository := &MockTransactionRepository{}
	handler := NewTransactionHandler(repository)

	if handler.repository != repository {
		t.Errorf("The repository field for the handler wasn't assigned. Expect %s but got %s", repository, handler.repository)
	}
}
