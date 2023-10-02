package model

type Transaction struct {
	TransactionId   uint64  `json:"transaction_id"`
	AccountId       uint64  `json:"account_id"`
	OperationTypeId uint32  `json:"operation_type_id"`
	Amount          float32 `json:"amount"`
}
