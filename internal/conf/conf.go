package conf

import "github.com/spf13/viper"

// Config struct will store the configuration values provided by user
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
	configName  = "config"
	configPaths = []string{
		".",
		"/etc/",
		"$HOME/",
	}
)

func LoadConfig() (config Config, err error) {
	for k, v := range defaults {
		viper.SetDefault(k, v)
	}
	for _, p := range configPaths {
		viper.AddConfigPath(p)
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("toml")

	viper.SetEnvPrefix("kc_ssh")  // Becomes "KC_SSH"
	viper.BindEnv("Realm")        // KC_SSH_REALM
	viper.BindEnv("Endpoint")     // KC_SSH_ENDPOINT
	viper.BindEnv("ClientID")     // KC_SSH_CLIENTID
	viper.BindEnv("ClientSecret") // KC_SSH_CLIENTSECRET
	viper.BindEnv("ClientScope")  // KC_SSH_CLIENTSCOPE

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
