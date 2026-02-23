package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host              string
	Port              string
	User              string
	Password          string
	Name              string
	SSLMode           string
	BybitAPIKey       string
	BybitAPISecretKey string
	MEXCAPIKey        string
	MEXCAPISecretKey  string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		return nil, err
	}

	return &Config{
		Host:              os.Getenv("DB_HOST"),
		Port:              os.Getenv("DB_PORT"),
		User:              os.Getenv("DB_USER"),
		Password:          os.Getenv("DB_PASSWORD"),
		Name:              os.Getenv("DB_NAME"),
		SSLMode:           os.Getenv("DB_SSL_MODE"),
		BybitAPIKey:       os.Getenv("BYBIT_API_KEY"),
		BybitAPISecretKey: os.Getenv("BYBIT_SECRET_KEY"),
		MEXCAPIKey:        os.Getenv("MEXC_API_KEY"),
		MEXCAPISecretKey:  os.Getenv("MEXC_SECRET_KEY"),
	}, nil
}

func (c *Config) GetDSN() string {
	return "postgres://" + c.User + ":" + c.Password + "@" + c.Host + ":" + c.Port + "/" + c.Name + "?sslmode=" + c.SSLMode
}
