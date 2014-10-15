package config

import "github.com/spf13/viper"

var (
	PG_CONFIG             = ""
	BUCKET                = ""
	EMBER_APP             = ""
	ASSETS                = ""
	ENV                   = ""
	AWS_ACCESS_KEY_ID     = ""
	AWS_SECRET_ACCESS_KEY = ""
	SESSION_AUTH_KEY      = ""
	SESSION_CRYPT_KEY     = ""
)

func init() {
	viper.SetDefault("ENV", "development")
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/flowfeeds")
	viper.AddConfigPath("$HOME/.flowfeeds")
	viper.ReadInConfig()

	PG_CONFIG = viper.GetString("PG_CONFIG")
	EMBER_APP = viper.GetString("EMBER_APP")
	BUCKET = viper.GetString("BUCKET")
	ASSETS = viper.GetString("ASSETS")
	ENV = viper.GetString("ENV")
	AWS_ACCESS_KEY_ID = viper.GetString("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY = viper.GetString("AWS_SECRET_ACCESS_KEY")

	SESSION_AUTH_KEY = viper.GetString("SESSION_AUTH_KEY")
	SESSION_CRYPT_KEY = viper.GetString("SESSION_CRYPT_KEY")
}
