package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	//"github.com/go-chi/chi"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aniljaiswalcs/pismo/model"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) CreateAccount(ctx context.Context, account model.Account) (*model.Account, error) {
	args := m.Called(ctx, account)
	return args.Get(0).(*model.Account), args.Error(1)
}

func (m *MockAccountRepository) FindAccount(ctx context.Context, accountId uint64) (*model.Account, error) {
	args := m.Called(ctx, accountId)
	return args.Get(0).(*model.Account), args.Error(1)
}

func TestGetAccount(t *testing.T) {
	mockRepo := new(MockAccountRepository)
	handler := &AccountHandler{repository: mockRepo}

	req, err := http.NewRequest("GET", "/v1/accounts/123", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	expectedAccountID := uint64(123)
	expectedDocumentNumber := uint64(44)
	expectedAccount := &model.Account{
		AccountId:      expectedAccountID,
		DocumentNumber: expectedDocumentNumber,
	}
	mockRepo.On("FindAccount", mock.Anything, expectedAccountID).Return(expectedAccount, nil)

	router := mux.NewRouter()
	accountMux := router.PathPrefix("/v1/accounts").Subrouter()
	accountMux.HandleFunc("/{accountId:[0-9]+}", handler.GetAccount).Methods("GET")

	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var responseAccount model.Account
	if err := json.Unmarshal(rr.Body.Bytes(), &responseAccount); err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}
	assert.Equal(t, expectedAccount, &responseAccount)
}

func TestGetAccountFailsWhenInvalidRequest(t *testing.T) {
	var scenarios = []struct {
		description        string
		accountId          uint64
		expectedError      error
		expectedStatusCode int
	}{
		{
			"Zero value test ",
			0,
			errors.New("Error!"),
			http.StatusBadRequest,
		},
		{
			"Row not found",
			99,
			sql.ErrNoRows,
			http.StatusNotFound,
		},
		{
			"Database errror",
			0,
			errors.New("Database error!"),
			http.StatusBadRequest,
		},
	}

	for _, scenario := range scenarios {
		mockRepo := new(MockAccountRepository)
		handler := &AccountHandler{repository: mockRepo}

		req, err := http.NewRequest("GET", fmt.Sprintf("/v1/accounts/%d", scenario.accountId), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		mockRepo.On("FindAccount", mock.Anything, scenario.accountId).Return(&model.Account{}, scenario.expectedError)

		router := mux.NewRouter().SkipClean(true)
		accountMux := router.PathPrefix("/v1/accounts").Subrouter()
		accountMux.HandleFunc("/{accountId:[0-9]+}", handler.GetAccount).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != scenario.expectedStatusCode {
			t.Errorf("Expected status code %d but got %d for the test case %s", scenario.expectedStatusCode, rr.Code, scenario.description)
		}
	}
}

func TestCreateAccount(t *testing.T) {
	mockRepo := new(MockAccountRepository)
	handler := NewAccountHandler(mockRepo)

	payload := &AccountPayload{
		DocumentNumber: 123456,
	}
	requestBody, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "/accounts", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	expectedAccount := &model.Account{
		AccountId:      payload.AccountId,
		DocumentNumber: payload.DocumentNumber,
	}
	mockRepo.On("CreateAccount", mock.Anything, *expectedAccount).Return(expectedAccount, nil)

	handler.CreateAccount(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)

	var responseAccount model.Account
	if err := json.Unmarshal(rr.Body.Bytes(), &responseAccount); err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}
	assert.Equal(t, expectedAccount.AccountId, responseAccount.AccountId)
	assert.Equal(t, expectedAccount.DocumentNumber, responseAccount.DocumentNumber)
}

func TestCreateAccountFailsWhenInvalidRequest(t *testing.T) {
	testCases := []struct {
		Name             string
		Payload          *AccountPayload
		ExpectedCode     int
		ExpectedResponse string
		ExpectedReturn   *model.Account
	}{
		{
			Name:             "Invalid JSON Payload",
			Payload:          &AccountPayload{},
			ExpectedCode:     http.StatusBadRequest,
			ExpectedResponse: "the document_number must be a valid positive integer",
			ExpectedReturn:   &model.Account{},
		},
		{
			Name: "Account number instead of Document number",
			Payload: &AccountPayload{
				AccountId: 222,
			},
			ExpectedCode:     http.StatusBadRequest,
			ExpectedResponse: "the document_number must be a valid positive integer",
			ExpectedReturn:   &model.Account{},
		},
	}
	mockRepo := new(MockAccountRepository)
	handler := NewAccountHandler(mockRepo)
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			requestBody, _ := json.Marshal(tc.Payload)

			req, err := http.NewRequest("POST", "v1/accounts", bytes.NewReader(requestBody))
			req.Header.Set("Content-Type", "application/json")
			if err != nil {
				t.Fatal(err)
			}
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()
			mockRepo.On("CreateAccount", mock.Anything, *tc.ExpectedReturn).Return(tc.ExpectedReturn, nil)

			handler.CreateAccount(recorder, req)

			assert.Equal(t, tc.ExpectedCode, recorder.Code)
			expectedResponseJson := map[string]string{}
			actualResponseJson := map[string]string{}

			json.Unmarshal([]byte(tc.ExpectedResponse), &expectedResponseJson)
			json.Unmarshal([]byte(recorder.Body.Bytes()), &actualResponseJson)

			if !reflect.DeepEqual(expectedResponseJson, actualResponseJson) {
				t.Errorf("Expected response body %s but got %s", tc.ExpectedResponse, recorder.Body.String())
			}
		})
	}
}
