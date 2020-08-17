package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// Config struct for operator config yaml
type Config struct {
	Slack Slack `yaml:"slack"`
}

// Slack for config yaml structure
type Slack struct {
	APIToken APIToken `yaml:"APIToken"`
}

// APIToken for config yaml structure
type APIToken struct {
	SecretName string `yaml:"secretName"`
	Key        string `yaml:"key"`
}

var log = zap.Logger(true)

func readConfig(filePath string) (*Config, error) {
	var config Config

	// Read YML
	source, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshall
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// GetOperatorConfig returns the config object for the operator
func GetOperatorConfig() (*Config, error) {
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if len(configFilePath) == 0 {
		configFilePath = "config/operator/default-config.yaml"
	}

	log.Info("Reading config file", "configFilePath", configFilePath)
	config, err := readConfig(configFilePath)
	if err != nil {
		return nil, err
	}
	return config, nil
}
