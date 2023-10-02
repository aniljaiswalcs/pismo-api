package model

import "testing"

func TestValidateOperationType(t *testing.T) {
	var scenarios = []struct {
		operationTypeId  uint32
		expectedResponse bool
	}{
		{
			1,
			true,
		},
		{
			2,
			true,
		},
		{
			3,
			true,
		},
		{
			4,
			true,
		},
		{
			0,
			false,
		},
		{
			5,
			false,
		},
	}

	for _, scenario := range scenarios {
		response := ValidateOperationType(scenario.operationTypeId)

		if response != scenario.expectedResponse {
			t.Errorf("Expected response to be %t but got %t", scenario.expectedResponse, response)
		}
	}
}

func TestValidateOperationTypeAmount(t *testing.T) {
	var scenarios = []struct {
		operationTypeId  uint32
		amount           float32
		expectedResponse bool
	}{
		{
			1,
			-100.0,
			true,
		},
		{
			2,
			-100.0,
			true,
		},
		{
			3,
			-100.0,
			true,
		},
		{
			4,
			100.0,
			true,
		},
		{
			1,
			100.0,
			false,
		},
		{
			2,
			100.0,
			false,
		},
		{
			3,
			100.0,
			false,
		},
		{
			4,
			-100.0,
			false,
		},
	}

	for _, scenario := range scenarios {
		response := ValidateOperationTypeAmount(scenario.operationTypeId, scenario.amount)

		if response != scenario.expectedResponse {
			t.Errorf("Expected response to be %t but got %t", scenario.expectedResponse, response)
		}
	}
}
