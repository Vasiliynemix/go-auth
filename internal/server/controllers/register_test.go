package controllers

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegisterRequestValidation(t *testing.T) {
	testCases := []struct {
		Request *RegisterRequest
		Valid   bool
		Error   error
	}{
		{
			Request: &RegisterRequest{
				Login:           "",
				Password:        "",
				ConfirmPassword: "",
			},
			Valid: false,
			Error: LoginRequiredError,
		},
		{
			Request: &RegisterRequest{
				Login:           "",
				Password:        "password",
				ConfirmPassword: "",
			},
			Valid: false,
			Error: LoginRequiredError,
		},
		{
			Request: &RegisterRequest{
				Login:           "login",
				Password:        "",
				ConfirmPassword: "",
			},
			Valid: false,
			Error: PasswordRequiredError,
		},
		{
			Request: &RegisterRequest{
				Login:           "login",
				Password:        "password",
				ConfirmPassword: "",
			},
			Valid: false,
			Error: PasswordRequiredError,
		},
		{
			Request: &RegisterRequest{
				Login:           "login",
				Password:        "password",
				ConfirmPassword: "qwerty",
			},
			Valid: false,
			Error: PasswordNotEqualError,
		},
		{
			Request: &RegisterRequest{
				Login:           "login",
				Password:        "password",
				ConfirmPassword: "password",
			},
			Valid: true,
		},
	}

	for _, testCase := range testCases {
		valid, err := testCase.Request.Validate(RegisterRequestValidationConfig{
			LoginRequired:     true,
			PasswordRequired:  true,
			PasswordEqual:     true,
			PasswordMinLength: 8,
		})
		require.Equal(t, testCase.Valid, valid)
		require.Equal(t, testCase.Error, err)
	}
}

func TestRegisterControllerCreation(t *testing.T) {
	c := NewRegisterController(nil, nil, nil)
	require.Equal(t, c, &RegisterController{})
	require.Equal(t, "/register", c.GetHandlers()[0].GetPath())
	require.Equal(t, "POST", c.GetHandlers()[0].GetMethod())
}

func TestRegisterControllerGroup(t *testing.T) {
	c := NewRegisterController(nil, nil, nil)
	require.Equal(t, "/auth", c.GetGroup())
}
