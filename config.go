package tweethog

import 	(
	"github.com/kylelemons/go-gypsy/yaml"
)

type Config struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

func NewConfig() *Config {
	result := new(Config)

	return result
}

func (config *Config) LoadFromFile(fileName string) error {
	yamlConfig, configErr := yaml.ReadFile(fileName)

	if configErr != nil {
		return configErr
	}

	if consumerKey, err := yamlConfig.Get("consumer_key"); err != nil {
		return err
	} else {
		config.ConsumerKey = consumerKey
	}

	if consumerSecret, err := yamlConfig.Get("consumer_secret"); err != nil {
		return err
	} else {
		config.ConsumerSecret = consumerSecret
	}

	if accessToken, err := yamlConfig.Get("access_token"); err != nil {
		return err
	} else {
		config.AccessToken = accessToken
	}

	if accessSecret, err := yamlConfig.Get("access_secret"); err != nil {
		return err
	} else {
		config.AccessSecret = accessSecret
	}

	return nil
}