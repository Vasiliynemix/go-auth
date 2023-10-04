package controllers

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAuthRequestValidation(t *testing.T) {
	testCases := []struct {
		Request *AuthRequest
		Valid   bool
		Error   error
	}{
		{
			Request: &AuthRequest{
				Login:    "",
				Password: "",
			},
			Valid: false,
			Error: LoginRequiredError,
		},
		{
			Request: &AuthRequest{
				Login:    "",
				Password: "password",
			},
			Valid: false,
			Error: LoginRequiredError,
		},
		{
			Request: &AuthRequest{
				Login:    "login",
				Password: "",
			},
			Valid: false,
			Error: PasswordRequiredError,
		},
		{
			Request: &AuthRequest{
				Login:    "login",
				Password: "password",
			},
			Valid: true,
		},
	}

	for _, testCase := range testCases {
		valid, err := testCase.Request.Validate()
		require.Equal(t, testCase.Valid, valid)
		require.Equal(t, testCase.Error, err)
	}
}

func TestAuthControllerCreation(t *testing.T) {
	c := NewAuthController(nil, nil)
	require.Equal(t, c, &AuthController{})
	require.Equal(t, "/", c.GetHandlers()[0].GetPath())
	require.Equal(t, "POST", c.GetHandlers()[0].GetMethod())
}

func TestAuthControllerGroup(t *testing.T) {
	c := NewAuthController(nil, nil)
	require.Equal(t, "/auth", c.GetGroup())
}
