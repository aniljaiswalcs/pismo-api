package model

type Account struct {
	AccountId      uint64 `json:"account_id,omitempty"`
	DocumentNumber uint64 `json:"document_number"`
}
