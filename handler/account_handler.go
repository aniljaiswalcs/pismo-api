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

type AccountHandler struct {
	repository repository.AccountRepository
}

func NewAccountHandler(repository repository.AccountRepository) *AccountHandler {
	return &AccountHandler{
		repository: repository,
	}
}

func (c *AccountHandler) CreateAccount(w http.ResponseWriter, req *http.Request) {

	newCtx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	payload := &AccountPayload{}
	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		lib.RenderJSON(w, http.StatusBadRequest, err.Error())
	}

	documentNumber := payload.DocumentNumber
	if documentNumber <= 0 {
		lib.RenderJSON(w, http.StatusBadRequest, lib.DocumentNumberError)
		return
	}

	account, err := c.repository.CreateAccount(newCtx, model.Account{
		DocumentNumber: documentNumber,
	})

	if err != nil {
		if err.Error() == lib.DatabaseTimeoutError {
			lib.RenderJSON(w, http.StatusInternalServerError, lib.TimeoutError)
			return
		} else if err.Error() == lib.ContextDeadline {
			lib.RenderJSON(w, http.StatusInternalServerError, lib.TimeoutError)
			return
		}
		lib.RenderJSON(w, http.StatusInternalServerError, lib.AccountCreationError)
		return
	}

	lib.RenderJSON(w, http.StatusCreated, account)

}

func (c *AccountHandler) GetAccount(w http.ResponseWriter, req *http.Request) {

	newCtx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	accountIdParam := mux.Vars(req)["accountId"]
	accountId, err := strconv.ParseUint(accountIdParam, 10, 64)
	if err != nil {
		lib.RenderJSON(w, http.StatusBadRequest, lib.ParsingAccountID)
		return
	}
	if accountId <= 0 {
		lib.RenderJSON(w, http.StatusBadRequest, lib.AccountIdValidation)
		return
	}

	account, err := c.repository.FindAccount(newCtx, accountId)

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

type AccountPayload struct {
	AccountId      uint64 `json:"account_id,omitempty"`
	DocumentNumber uint64 `json:"document_number"`
}
