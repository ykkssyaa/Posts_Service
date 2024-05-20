package config

import (
	"github.com/joho/godotenv"
)

func InitConfig() error {
	err := godotenv.Load()

	return err
}
