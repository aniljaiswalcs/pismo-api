package lib

const (
	StatusInvalidRequest = "Invalid Request"
	StatusCodeBadRequest = "Bad Request"
	StatusServerError    = "Server Error"
	StatusForbidden      = "Request Forbidden"

	DocumentNumberError = "the document_number must be a valid positive integer"

	//Acoount
	AccountCreationError = "an error occurred when creating the account"
	ParsingAccountID     = "error in parsing accountId"
	AccountIdValidation  = "the account_id must be a valid positive integer"
	AccountIdNotFound    = "no account found for the provided account ID"

	//opertaion
	OperationTypeIdError = "the operation_type_id must be one of the following valid values: 1, 2, 3, 4"
	OperationTypeError   = "purchases and withdraw operations must have a negative amount. Payment operations must have a positive amount."

	//database
	DatabaseTimeoutError = "timeout: context deadline exceeded"
	DatabaseError        = "an error occurred when fetching the account from the database"
	TimeoutError         = "timeout during operation. Try Again"

	//context
	ContextDeadline = "context deadline exceeded"
)
