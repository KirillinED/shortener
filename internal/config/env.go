package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func ParseEnv(c *config) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = env.Parse(c)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	v, ok := os.LookupEnv("APP_PORT")
	if ok {
		c.Address.Host = v
	}
}
