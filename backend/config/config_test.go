package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigurationVariables(t *testing.T) {
	testConfig, err := LoadTestConfigurationVariables()
	assert.NoError(t, err, "expected no error loading test configuration variables, got: %v.", err)

	assert.NotEmpty(t, testConfig.TestDatabaseHostname)
	assert.NotEmpty(t, testConfig.TestDatabasePort)
	assert.NotEmpty(t, testConfig.TestDatabaseName)
	assert.NotEmpty(t, testConfig.TestDatabaseUsername)
	assert.NotEmpty(t, testConfig.TestDatabasePassword)
	assert.NotEmpty(t, testConfig.TestDatabaseTimezone)
	assert.NotEmpty(t, testConfig.TestAccessTokenSecretKey)

	t.Logf("Successfully loaded and verified configuration variables for testing database.")
}
