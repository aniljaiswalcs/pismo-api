CREATE TABLE IF NOT EXISTS "transactions" (
    "transaction_id" SERIAL PRIMARY KEY,
    "account_id" INT NOT NULL,
    "operation_type_id" INT NOT NULL,
    "amount" NUMERIC(12, 4) NOT NULL,
    CONSTRAINT fk_account
      FOREIGN KEY(account_id) 
	  REFERENCES accounts(account_id),
    CONSTRAINT fk_operation_type
      FOREIGN KEY(operation_type_id) 
	  REFERENCES operation_types(operation_type_id)
);
