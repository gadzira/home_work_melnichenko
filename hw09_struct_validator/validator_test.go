package hw09_struct_validator //nolint:golint,stylecheck

import (
	"encoding/json"
	"strings"
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
		expectedError string
	}{
		{
			name:          "fail user",
			in:            fuser,
			expectedError: "incoming value is more than a maximum limit",
		},
		{
			name:          "fail app",
			in:            fapp,
			expectedError: "length is invalid",
		},
		{
			name:          "fail response",
			in:            fresponse,
			expectedError: "not included in the set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.in)
			require.Equal(t, tt.expectedError, getErrMsg(err))
		})
	}
}

// I'm ashamed, but I don't regret for the code below
func getErrMsg(e error) string {
	s := e.Error()
	mes := strings.Split(s, ":")
	return strings.TrimSpace(mes[len(mes)-1])
}
