package config

import (
	"github.com/spf13/viper"
)

type Environment struct {
	DbName     string `mapstructure:"DB_NAME"`
	DbHost     string `mapstructure:"DB_HOST"`
	DbPort     string `mapstructure:"DB_PORT"`
	DbUser     string `mapstructure:"DB_USER"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
	ApiPort    string `mapstructure:"API_PORT"`
	JwtSecret  string `mapstructure:"JWT_SECRET"`
}

func LoadEnvironment() (Environment, error) {
	var env Environment

	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()

	if err != nil {
		return env, err
	}

	err = viper.Unmarshal(&env)

	if err != nil {
		return env, err
	}

	return env, nil
}
