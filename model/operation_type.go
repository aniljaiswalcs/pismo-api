package model

const CASH_PURCHASE = 1
const INSTALLMENT_PURCHASE = 2
const WITHDRAW = 3
const PAYMENT = 4

func ValidateOperationType(operationTypeId uint32) bool {
	for _, operationType := range getOperationTypes() {
		if operationType == operationTypeId {
			return true
		}
	}

	return false
}

func ValidateOperationTypeAmount(operationTypeId uint32, amount float32) bool {
	switch operationTypeId {
	case CASH_PURCHASE, INSTALLMENT_PURCHASE, WITHDRAW:
		if amount >= 0 {
			return false
		}
	case PAYMENT:
		if amount <= 0 {
			return false
		}
	}

	return true
}

func getOperationTypes() []uint32 {
	return []uint32{
		CASH_PURCHASE,
		INSTALLMENT_PURCHASE,
		WITHDRAW,
		PAYMENT,
	}
}
