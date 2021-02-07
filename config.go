package main

import (
	"encoding/base64"

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

	AuthUser string `env:"GOKINDLE_AUTH_USER"`
	AuthPass string `env:"GOKINDLE_AUTH_PASS"`
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

	// auth password is base64 encoded
	pass, err := base64.RawStdEncoding.DecodeString(config.AuthPass)
	if err != nil {
		panic(err)
	}
	config.AuthPass = string(pass)
	return config
}
