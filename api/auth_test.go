package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthFunctions(t *testing.T) {
	auth := NewAuth("mock-secret-key")

	userID := 1

	accessToken, err := auth.GenerateAccessToken(userID)
	assert.NoError(t, err, "expected no error generating access token, got %v.", err)

	validatedID, err := auth.ValidateAccessToken(accessToken)
	assert.NoError(t, err, "expected no error validating mock access token, got: %v.", err)
	assert.Equal(t, userID, validatedID)

	invalidAccessToken := "invalid-access-token"
	invalidUserID, err := auth.ValidateAccessToken(invalidAccessToken)

	assert.Error(t, err, "expected error validating invalid access token, got: %v.", err)
	assert.Equal(t, 0, invalidUserID)
}
