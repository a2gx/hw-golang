package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
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

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:    "123456789012345678901234567890123456",
				Age:   21,
				Email: "text@text.test",
				Role:  UserRole("admin"),
				Phones: []string{
					"12345678901",
					"10987654321",
				},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:    "1234567890",
				Age:   12,
				Email: "test@test@com",
				Role:  UserRole("legal"),
				Phones: []string{
					"1234567890",
					"0987654321",
				},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: fmt.Errorf("length must be exactly %d", 36)},
				{Field: "Age", Err: fmt.Errorf("must be greater than %d", 18)},
				{Field: "Email", Err: fmt.Errorf("must be a valid regular expression")},
				{Field: "Role", Err: fmt.Errorf("must be one of admin,stuff")},
				{Field: "Phones", Err: fmt.Errorf("length must be exactly %d", 11)},
			},
		},
		{
			in: App{
				Version: "1.2.3",
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "1.2.3.4",
			},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: fmt.Errorf("length must be exactly %d", 5)},
			},
		},
		{
			in: Token{
				Header:    []byte("Header"),
				Payload:   []byte("Payload"),
				Signature: []byte("Signature"),
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 200,
				Body: "Hello World",
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tt.expectedErr) {
				t.Errorf("in: %v", reflect.TypeOf(tt.in).String())
				t.Errorf("wait: %v", tt.expectedErr)
				t.Errorf("fact: %v", err)
			}
		})
	}
}
