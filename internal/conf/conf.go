package conf

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Realm        string
	Endpoint     string
	ClientID     string
	ClientSecret string
	ClientScope  string
}

var (
	defaults = map[string]interface{}{
		"realm":    "master",
		"endpoint": "localhost",
		"scope":    "openid",
	}
)

func LoadConfig(configPath ...string) (config Config, err error) {
	for k, v := range defaults {
		viper.SetDefault(k, v)
	}

	configFilePath := os.Getenv("KC_SSH_CONFIG")

	viper.SetConfigFile(configFilePath)
	viper.SetConfigType("toml")
	viper.AddConfigPath("/opt/kc-ssh-pam")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config")
	viper.SetEnvPrefix("kc_ssh")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
