package configs

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	DSN string
}

type JWTConfig struct {
	Secret            string
	AccessExpiryHour  int `mapstructure:"access_expiry_hour"`
	RefreshExpiryHour int `mapstructure:"refresh_expiry_hour"`
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 1. Read config.yaml first (for defaults)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// 2. Read .env file (if it exists) to override defaults locally
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.MergeInConfig() // Ignore error if .env doesn't exist

	// 3. Environment Variables (Highest Priority)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
