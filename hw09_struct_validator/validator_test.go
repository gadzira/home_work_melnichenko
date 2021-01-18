package hw09_struct_validator //nolint:golint,stylecheck

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	BrokenResponse struct {
		Code int    `validate:""`
		Body string `json:"omitempty"`
	}

	EmptyStruct struct {
		EmptyField string `validate:"actually_empty_tag"`
	}
)

var user = User{
	ID:     "f2223b6a-398d-11eb-adc1-0242ac120002",
	Name:   "Bob Ross",
	Age:    37,
	Email:  "user@domen.com",
	Role:   "admin",
	Phones: []string{"79257878787", "79267878787", "79277878787"},
	meta:   []byte("vip_client"),
}

var fuser = User{
	ID:     "f2223b6a-398d-11eb-adc1-0242ac120002",
	Name:   "Bob Ross",
	Age:    57,
	Email:  "user@domen.com",
	Role:   "admin",
	Phones: []string{"79257878787", "79267878787", "79277878787"},
	meta:   []byte("vip_client"),
}

var app = App{
	Version: "55555",
}

var fapp = App{
	Version: "4444",
}

var token = Token{
	Header:    []byte("header"),
	Payload:   []byte("payload"),
	Signature: []byte("signature"),
}

var response = Response{
	Code: 404,
	Body: "Not Found",
}

var fresponse = Response{
	Code: 715,
	Body: "Good luck with that",
}

var justEmptyStruct = EmptyStruct{
	EmptyField: "",
}

func TestValidatePositive(t *testing.T) {
	tests := []struct {
		name string
		in   interface{}
	}{
		{
			name: "User",
			in:   user,
		},
		{
			name: "App",
			in:   app,
		},
		{
			name: "Token",
			in:   token,
		},
		{
			name: "Response",
			in:   response,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.in)
			require.NoError(t, err)
		})
	}
}

func TestValidateNegative(t *testing.T) {
	tests := []struct {
		name          string
		in            interface{}
		expectedField string
		expectedError error
	}{
		{
			name:          "fail user",
			in:            fuser,
			expectedField: "Age",
			expectedError: ErrMax,
		},
		{
			name:          "fail app",
			in:            fapp,
			expectedField: "Version",
			expectedError: ErrLen,
		},
		{
			name:          "fail response",
			in:            fresponse,
			expectedField: "Code",
			expectedError: ErrIncluded,
		},
		{
			name:          "Empty struct",
			in:            justEmptyStruct,
			expectedField: "-",
			expectedError: ErrDefault,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.in)

			var errs ValidationErrors
			if errors.As(err, &errs) {
				for i := range errs {
					require.Equal(t, tt.expectedField, errs[i].Field)
					require.True(t, errors.Is(errs[i].Err, tt.expectedError), "actual: %v", errs[i].Err)
				}

			}
		})
	}
}

func TestValidateNil(t *testing.T) {
	err := Validate(nil)
	require.Nil(t, err, "expect nil, actual: %v", err)
}
func TestValidatePointer(t *testing.T) {
	v := Response{200, "ok"}
	err := Validate(&v)
	require.NoError(t, err, "expect no error, actual: %v", err)
}

func TestValidateNotConditions(t *testing.T) {
	r := BrokenResponse{500, "Internal Server Error"}
	err := Validate(r)
	require.Nil(t, err, "expect nil, actual: %v", err)
}

func TestValidateInt(t *testing.T) {
	err := Validate(7)
	require.True(t, errors.Is(err, ErrNotStruct), "actual: %v", err)
}

func TestValidateString(t *testing.T) {
	err := Validate("A")
	require.True(t, errors.Is(err, ErrNotStruct), "actual: %v", err)
}
