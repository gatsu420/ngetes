package config

import "github.com/spf13/viper"

type Config struct {
	PostgresHost     string
	PostgresPort     int
	PostgresDB       string
	PostgresUser     string
	PostgresPassword string

	RedisHost     string
	RedisPort     int
	RedisDB       int
	RedisPassword string

	TokenSecretKey string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("./.env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	viper.AutomaticEnv()

	return &Config{
		PostgresHost:     viper.GetString("POSTGRES_HOST"),
		PostgresPort:     viper.GetInt("POSTGRES_PORT"),
		PostgresDB:       viper.GetString("POSTGRES_DB"),
		PostgresUser:     viper.GetString("POSTGRES_USER"),
		PostgresPassword: viper.GetString("POSTGRES_PASSWORD"),

		RedisHost:     viper.GetString("REDIS_HOST"),
		RedisPort:     viper.GetInt("REDIS_PORT"),
		RedisDB:       viper.GetInt("REDIS_DB"),
		RedisPassword: viper.GetString("REDIS_PASSWORD"),

		TokenSecretKey: viper.GetString("TOKEN_SECRET_KEY"),
	}, nil
}
