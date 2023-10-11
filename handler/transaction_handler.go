package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/aniljaiswalcs/pismo/model"
	"github.com/aniljaiswalcs/pismo/pkg/lib"
	"github.com/aniljaiswalcs/pismo/repository"
	"github.com/gorilla/mux"
)

type TransactionHandler struct {
	repository repository.TransactionRepository
}

func NewTransactionHandler(repository repository.TransactionRepository) *TransactionHandler {
	return &TransactionHandler{
		repository: repository,
	}
}

func (c *TransactionHandler) CreateTransaction(w http.ResponseWriter, req *http.Request) {
	payload := &TransactionPayload{}

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		lib.RenderJSON(w, http.StatusBadRequest, err.Error())
	}

	payloadErrors := validatePayload(payload)

	if len(payloadErrors) > 0 {
		lib.RenderJSON(w, http.StatusBadRequest, payloadErrors)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	transaction, err := c.repository.CreateTransaction(ctx, model.Transaction{
		AccountId:       payload.AccountId,
		OperationTypeId: payload.OperationTypeId,
		Amount:          payload.Amount,
	})

	if err != nil {
		if err.Error() == lib.DatabaseTimeoutError {
			lib.RenderJSON(w, http.StatusInternalServerError, lib.TimeoutError)
			return
		} else if err.Error() == lib.ContextDeadline {
			lib.RenderJSON(w, http.StatusInternalServerError, lib.TimeoutError)
			return
		}
		lib.RenderJSON(w, http.StatusBadRequest, lib.AccountIdNotFound)
		return
	}

	lib.RenderJSON(w, http.StatusCreated, transaction)
}

func validatePayload(payload *TransactionPayload) []string {
	var errors []string

	if payload.AccountId <= 0 {
		errors = append(errors, lib.AccountIdValidation)
	}

	if !model.ValidateOperationType(payload.OperationTypeId) {
		errors = append(errors, lib.OperationTypeIdError)
	}

	if !model.ValidateOperationTypeAmount(payload.OperationTypeId, payload.Amount) {
		errors = append(errors, lib.OperationTypeError)
	}

	return errors
}

type TransactionPayload struct {
	AccountId       uint64  `json:"account_id"`
	OperationTypeId uint32  `json:"operation_type_id"`
	Amount          float32 `json:"amount"`
}

func (c *TransactionHandler) GetAccount(w http.ResponseWriter, req *http.Request) {

	newCtx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	accountIdParam := mux.Vars(req)["transactionid"]
	accountId, err := strconv.ParseUint(accountIdParam, 10, 64)
	if err != nil {
		lib.RenderJSON(w, http.StatusBadRequest, lib.ParsingAccountID)
		return
	}
	if accountId <= 0 {
		lib.RenderJSON(w, http.StatusBadRequest, lib.AccountIdValidation)
		return
	}

	account, err := c.repository.FindtransactionAccount(newCtx, accountId)

	if err != nil {
		if err == sql.ErrNoRows {
			lib.RenderJSON(w, http.StatusNotFound, lib.AccountIdNotFound)
			return
		} else if err.Error() == lib.DatabaseTimeoutError || err.Error() == lib.ContextDeadline {
			lib.RenderJSON(w, http.StatusInternalServerError, lib.TimeoutError)
			return
		}
		lib.RenderJSON(w, http.StatusInternalServerError, lib.DatabaseError)
		return
	}

	lib.RenderJSON(w, http.StatusOK, account)
}
