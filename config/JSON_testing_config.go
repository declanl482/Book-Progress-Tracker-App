package config

import (
	"encoding/json"
	"os"
)

type JSONTestingConfig struct {
	TestDatabaseHostname     string `json:"test_database_hostname"`
	TestDatabasePort         string `json:"test_database_port"`
	TestDatabaseName         string `json:"test_database_name"`
	TestDatabaseUsername     string `json:"test_database_username"`
	TestDatabasePassword     string `json:"test_database_password"`
	TestDatabaseTimezone     string `json:"test_database_timezone"`
	TestAccessTokenSecretKey string `json:"test_access_token_secret_key"`
}

// LoadJSONConfig loads the configuration values from a JSON file.
func LoadJSONTestingConfigurationVariables(filename string) (*JSONTestingConfig, error) {
	testingConfig := &JSONTestingConfig{}

	// Read the JSON file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the JSON data into the config struct
	err = json.NewDecoder(file).Decode(testingConfig)
	if err != nil {
		return nil, err
	}

	return testingConfig, nil
}
