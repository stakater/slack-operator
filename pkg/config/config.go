package config

import (
	"os"
	"time"

	util "github.com/stakater/operator-utils/util"
	secretsUtil "github.com/stakater/operator-utils/util/secrets"
	"gopkg.in/yaml.v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

const (
	ErrorRequeueTime = 15 * time.Minute

	SlackDefaultSecretName     string = "slack-secret"
	SlackAPITokenSecretKey     string = "APIToken"
	SlackUserAPITokenSecretKey string = "UserAPIToken"
)

var (
	setupLog               = ctrl.Log.WithName("setup")
	SlackSecretName string = getConfigSecretName()
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

var log = zap.New()

func readConfig(filePath string) (*Config, error) {
	var config Config

	// Read YML
	source, err := os.ReadFile(filePath)
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

func getConfigSecretName() string {
	configSecretName, _ := os.LookupEnv("CONFIG_SECRET_NAME")
	if len(configSecretName) == 0 {
		configSecretName = SlackDefaultSecretName
		setupLog.Info("CONFIG_SECRET_NAME is unset, using default value: " + SlackDefaultSecretName)
	}
	return configSecretName
}

func ReadSlackTokenSecret(k8sReader client.Reader) (string, string) {
	operatorNamespace, _ := os.LookupEnv("OPERATOR_NAMESPACE")
	if len(operatorNamespace) == 0 {
		operatorNamespaceTemp, err := util.GetOperatorNamespace()
		if err != nil {
			setupLog.Error(err, "Unable to get operator namespace")
			os.Exit(1)
		}
		operatorNamespace = operatorNamespaceTemp
	}

	token, err := secretsUtil.LoadSecretData(k8sReader, SlackSecretName, operatorNamespace, SlackAPITokenSecretKey)
	if err != nil {
		setupLog.Error(err, "Could not read API token from key", "secretName", SlackSecretName, "secretKey", SlackAPITokenSecretKey)
		os.Exit(1)
	}

	userToken, err := secretsUtil.LoadSecretData(k8sReader, SlackSecretName, operatorNamespace, SlackUserAPITokenSecretKey)
	if err != nil {
		setupLog.Error(err, "Could not read API token from key", "secretName", SlackSecretName, "secretKey", SlackUserAPITokenSecretKey)
		os.Exit(1)
	}

	return token, userToken
}
