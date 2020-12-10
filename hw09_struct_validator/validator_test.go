package hw09_struct_validator //nolint:golint,stylecheck

import (
	"fmt"
	"testing"
)

// Test the function on different structures and other types.
type TestStruct struct {
	ID   string `json:"id" validate:"len:36"`
	Name string
	Age  int `validate:"min:18|max:50"`
}

var tstStruct1 = TestStruct{
	ID:   "f2223b6a-398d-11eb-adc1-0242ac120002",
	Name: "Bob Ross",
	Age:  37,
}

var tstStruct2 = 1

func TestValidate(t *testing.T) {

	err := Validate(tstStruct1)
	if err != nil {
		fmt.Errorf("Something go wrong: %w", err)
	}
}
