package main

import (
	env "github.com/Netflix/go-env"
	"github.com/joho/godotenv"
)

type configuration struct {
	SMTPServer string `env:"GOKINDLE_SMTP_SERVER"`
	SMTPPort   int    `env:"GOKINDLE_SMTP_PORT"`
	SMTPUser   string `env:"GOKINDLE_SMTP_USER"`
	SMTPPass   string `env:"GOKINDLE_SMTP_PASS"`

	FromEmail string `env:"GOKINDLE_FROM_EMAIL"`
	ToEmail   string `env:"GOKINDLE_TO_EMAIL"`
}

func readConfig() *configuration {
	config := &configuration{}

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	_, err := env.UnmarshalFromEnviron(config)
	if err != nil {
		panic(err)
	}
	return config
}
