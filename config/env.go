package config

import (
	"os"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBProtocol string
	DB         string
	DBPort     string
	Secret     string
}

func InitConfig() *Config {
	return &Config{
		DBUser:     GetEnv("DB_USER"),
		DBPassword: GetEnv("DB_PASSWORD"),
		DBName:     GetEnv("DB_NAME"),
		DBProtocol: GetEnv("DB_PROTOCOL"),
		DB:         GetEnv("DB"),
		DBPort:     GetEnv("DB_PORT"),
		Secret:     GetEnv("SECRET"),
	}
}

func GetEnv(key string) string {
	value := os.Getenv(key)
	return value

}
