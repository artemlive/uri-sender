package config

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
)

const (
	slackAuthTokenEnvName = "SLACK_AUTH_TOKEN"
)

type Credentials struct {
	SlackApiToken string `json:"slack_api_token"`
}

type Config struct {
	credentials Credentials
	Notifiers   []struct {
		Type       string   `json:"type"`
		Recipients []string `json:"recipients"`
		Message    string   `json:"message"`
	} `json:"notifiers"`
}

func NewConfig(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return &Config{}, err
	}

	config := Config{}
	err = json.Unmarshal(file, &config)
	return &config, err
}

func (c *Config) GetSlackApiToken() (string, error) {
	if len(c.credentials.SlackApiToken) > 0 {
		log.Debug().Msgf("Found the SLACK API token in the config file, let's using it")
		return c.credentials.SlackApiToken, nil
	}
	envVar := c.getEnv(slackAuthTokenEnvName, "")
	if len(envVar) > 0 {
		log.Debug().Msgf("Found the SLACK API token in the ENV variable, let's using it")
		return envVar, nil
	}
	return "", fmt.Errorf("couldn't find a slack API token either in the config and ENV")
}

// consider using viper
func (c *Config) getEnv(envVar, defaultEnvVar string) string {
	osEnvVar := os.Getenv(envVar)
	if len(osEnvVar) == 0 {
		return defaultEnvVar
	}
	return osEnvVar
}
